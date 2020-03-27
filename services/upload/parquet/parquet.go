package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/gofrs/uuid"
	"github.com/jonny-rimek/wowmate/services/golib"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)
//IMPROVE: refactor to use up to date logging approach

//StepfunctionEvent provides config data to the lambda
type StepfunctionEvent struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

func handler(e StepfunctionEvent) error {
	uploadUUID := uuid.Must(uuid.NewV4()).String()
	targetBucket := os.Getenv("TARGET_BUCKET_NAME")

	log.Print("DEBUG: bucketname: " + e.BucketName)
	log.Print("DEBUG: filename: " + e.Key)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	file := &aws.WriteAtBuffer{}

	//30minutes of debugging later, S3 does in fact download into memory
	//and doesn't write to /tmp. or I could have just read the comment on
	//WriteAtBUFFER...
	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(e.BucketName),
			Key:    aws.String(e.Key),
		})
	if err != nil {
		log.Fatalf("Unable to download item %q, %v", e.Key, err)
	}
	log.Printf("DEBUG: Downloaded %v MB", numBytes/1024/1024)

	r := bytes.NewReader(file.Bytes())
	uncompressed, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	s := bufio.NewScanner(uncompressed)

	events, err := Import(s, uploadUUID)
	if err != nil {
		return err
	}

	log.Print("DEBUG: read combatlog to slice of Event structs")

	fw, err := local.NewLocalFileWriter("/tmp/flat.parquet")
	if err != nil {
		log.Fatalf("Can't create local file: %v", err)
	}

	log.Print("DEBUG: created local file")

	pw, err := writer.NewParquetWriter(fw, new(Event), 1) //4 is actually slower than 1 :o
	if err != nil {
		log.Fatalf("Can't create parquet writer: %v", err)
	}

	log.Print("DEBUG: created write")

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	for _, event := range events {
		if err = pw.Write(event); err != nil {
			log.Println("Write error", err)
		}

	}

	log.Print("DEBUG: wrote file to writer")

	if err = pw.WriteStop(); err != nil {
		log.Fatalf("WriteStop error: %v", err)
	}
	log.Println("DEBUG: Converting to parquet finished")

	fw.Close()
	uncompressed.Close()

	fr, err := local.NewLocalFileReader("/tmp/flat.parquet")
	if err != nil {
		log.Fatalf("Can't open file")
	}
	//END

	//file name with partitions per minute
	uploadFileName := fmt.Sprintf("%v/%v/%v/%v/%v/%v.parquet",
		time.Now().Year(),
		int(time.Now().Month()),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		strings.TrimPrefix(strings.TrimSuffix(e.Key, ".txt.gz"), "new/"),
	)

	err = golib.UploadFileToS3(fr, targetBucket, uploadFileName, sess)
	if err != nil {
		return err
	}

	newFilename := fmt.Sprintf("processed/%v", strings.TrimPrefix(e.Key, "new"))

	//If removed, remove write access to upload bucket
	svc := s3.New(sess)
	_, err = svc.CopyObject(&s3.CopyObjectInput{
		CopySource: aws.String(e.BucketName + "/" + e.Key),
		Bucket:     aws.String(e.BucketName),
		Key:        aws.String(newFilename),
	})
	if err != nil {
		log.Printf("unable to move file to processed dir. %v", err)
		return err
	}
	// Wait to see if the item got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(e.BucketName),
		Key:    aws.String(newFilename),
	})
	if err != nil {
		fmt.Printf("Error occurred while waiting for item %q to be copied to bucket processed folder", e.Key)
		return err
	}
	// Delete the item
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(e.BucketName), Key: aws.String(e.Key)})
	if err != nil {
		fmt.Printf("Unable to delete object %q from bucket %q", e.Key, e.BucketName)
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(e.BucketName),
		Key:    aws.String(e.Key),
	})
	if err != nil {
		fmt.Printf("Error occurred while waiting for object %q to be deleted from bucket %v", e.Key, e.BucketName)
		return err
	}

	err = os.Remove("/tmp/flat.parquet")
	if err != nil {
		log.Println("ERROR: failed to delete file")
		return err
	}
	log.Printf("DEBUG: file deleted")

	return nil
}

func main() {
	golib.InitLogging()
	lambda.Start(handler)
}
