package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jonny-rimek/wowmate/lambda/lib/combatlog/event"
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

//Event ..
type Event struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

func handler(e Event) error {
	targetBucket := os.Getenv("TARGET_BUCKET_NAME")

	log.Print("DEBUG: bucketname: " + e.BucketName)
	log.Print("DEBUG: filename: " + e.Key)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	downloader := s3manager.NewDownloader(sess)

	file := &aws.WriteAtBuffer{}

	numBytes, err := downloader.Download(
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(e.BucketName),
			Key:    aws.String(e.Key),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", e.Key, err)
	}

	log.Println("DEBUG: Downloaded", numBytes, "bytes")

	var records []CSV

	r := bytes.NewReader(file.Bytes())
	s := bufio.NewScanner(r)

	for s.Scan() {
		row := strings.Split(s.Text(), ",")

		bigint, err := strconv.ParseInt(row[1], 10, 32)
		d := int32(bigint)

		if err != nil {
			log.Fatalf("Failed to convert 2nd row to int32")
		}

		r := CSV{
			row[0],
			d,
		}

		records = append(records, r)
	}

	log.Print("DEBUG: read CSV into structs")

	fw, err := local.NewLocalFileWriter("/tmp/flat.parquet")
	if err != nil {
		log.Fatalf("Can't create local file: %v", err)
	}

	log.Print("DEBUG: created local file")

	pw, err := writer.NewParquetWriter(fw, new(CSV), 1) //4 is actually slower than 1 :o
	if err != nil {
		log.Fatalf("Can't create parquet writer: %v", err)
	}

	log.Print("DEBUG: created write")

	pw.RowGroupSize = 128 * 1024 * 1024 //128M
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	for _, r := range records {
		if err = pw.Write(r); err != nil {
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

	return nil
}

func main() {
	normalize.Normalize()
	event.Event()
	lambda.Start(handler)
}
