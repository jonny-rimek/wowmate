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

//SfnEvent provides config data to the lambda
type SfnEvent struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

//Response is the object the next step in Stepfunctions expects
type Response struct {
	UploadUUID    string `json:"upload_uuid"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	Day           int    `json:"day"`
	Hour          int    `json:"hour"`
	Minute        int    `json:"minute"`
	ParquetBucket string `json:"parquet_bucket"`
	ParquetFile   string `json:"parquet_file"`
}

func handler(e SfnEvent) (Response, error) {
	uploadUUID := uuid.Must(uuid.NewV4()).String()
	resp := Response{UploadUUID: uploadUUID}

	targetBucket := os.Getenv("TARGET_BUCKET_NAME")

	log.Print("DEBUG: bucketname: " + e.BucketName)
	log.Print("DEBUG: filename: " + e.Key)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
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
		return resp, err
	}
	s := bufio.NewScanner(uncompressed)

	events, err := Import(s, uploadUUID)
	if err != nil {
		return resp, err
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

	resp.Year = time.Now().Year()
	resp.Month = int(time.Now().Month())
	resp.Day = time.Now().Day()
	resp.Hour = time.Now().Hour()
	resp.Minute = time.Now().Minute()
	resp.ParquetBucket = targetBucket
	resp.ParquetFile = fmt.Sprintf("year=%v/month=%v/day=%v/hour=%v/minute=%v/%v.parquet",
		resp.Year,
		resp.Month,
		resp.Day,
		resp.Hour,
		resp.Minute,
		strings.TrimPrefix(strings.TrimSuffix(e.Key, ".txt.gz"), "new/"),
	)

	err = golib.UploadFileToS3(fr, targetBucket, resp.ParquetFile, sess)
	if err != nil {
		return resp, err
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
		return resp, err
	}
	// Wait to see if the item got copied
	err = svc.WaitUntilObjectExists(&s3.HeadObjectInput{
		Bucket: aws.String(e.BucketName),
		Key:    aws.String(newFilename),
	})
	if err != nil {
		fmt.Printf("Error occurred while waiting for item %q to be copied to bucket processed folder", e.Key)
		return resp, err
	}
	// Delete the item
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(e.BucketName), Key: aws.String(e.Key)})
	if err != nil {
		fmt.Printf("Unable to delete object %q from bucket %q", e.Key, e.BucketName)
		return resp, err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(e.BucketName),
		Key:    aws.String(e.Key),
	})
	if err != nil {
		fmt.Printf("Error occurred while waiting for object %q to be deleted from bucket %v", e.Key, e.BucketName)
		return resp, err
	}

	err = os.Remove("/tmp/flat.parquet")
	if err != nil {
		log.Println("ERROR: failed to delete file")
		return resp, err
	}
	log.Printf("DEBUG: file deleted")

	return resp, nil
}

func main() {
	golib.InitLogging()
	lambda.Start(handler)
}
