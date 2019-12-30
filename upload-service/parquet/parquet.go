package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/gofrs/uuid"
	"github.com/jonny-rimek/wowmate/lambda/lib/combatlog/normalize"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

//CSV ..
type CSV struct {
	Name string `parquet:"name=name, type=UTF8, encoding=PLAIN_DICTIONARY"`
	N1   int32  `parquet:"name=n1, type=INT32"`
}

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

	//FIXME: fatal error: runtime: out of memory 
	//REPORT RequestId: 489701bd-42c0-431c-bb6c-e5603b581720	Duration: 140806.08 ms	Billed Duration: 140900 ms	Memory Size: 3008 MB	Max Memory Used: 1730 MB	

	r := bytes.NewReader(file.Bytes())
	s := bufio.NewScanner(r)

	events, err := Import(s, uploadUUID) //TODO: handle errors
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Print("DEBUG: read combatlog to slice of Event structs")

	//WRITE TO PARQUET FILE
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

	fr, err := local.NewLocalFileReader("/tmp/flat.parquet")
	if err != nil {
		log.Fatalf("Can't open file")
	}
	//END

	//UPLOAD TO S3
	s3Svc := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(s3Svc)
	uploadFileName := fmt.Sprintf("test/test.parquet")

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(targetBucket),
		Key:    &uploadFileName,
		Body:   fr,
	})
	if err != nil {
		log.Println("Failed to upload to S3: " + err.Error())
	}

	log.Printf("DEBUG: Upload finished! location: %s", result.Location)
	//TODO: FIXME: delete file

	return nil
}

func main() {
	normalize.Normalize()
	//event.Event()
	lambda.Start(handler)
}
