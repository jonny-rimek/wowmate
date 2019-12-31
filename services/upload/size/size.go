package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//RequestParameters test ..
type RequestParameters struct {
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

//Detail test ..
type Detail struct {
	RequestParameters RequestParameters `json:"requestParameters"`
}

//Event ..
type Event struct {
	Detail Detail `json:"detail"`
}

//Response ..
type Response struct {
	FileSize   int    `json:"file_size"`
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

func handler(e Event) (Response, error) {

	log.Print("DEBUG: bucketname: " + e.Detail.RequestParameters.BucketName)
	log.Print("DEBUG: filename: " + e.Detail.RequestParameters.Key)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)

	// output, err := svc.GetObject(
	output, err := svc.HeadObject(
		&s3.HeadObjectInput{
			Bucket: aws.String(e.Detail.RequestParameters.BucketName),
			Key:    aws.String(e.Detail.RequestParameters.Key),
		})
	if err != nil {
		log.Fatalf("Unable to to send head request to item %q, %v", e.Detail.RequestParameters.Key, err)
	}

	var response Response
	MB := int(*output.ContentLength / 1024 / 1024)
	response.FileSize = MB
	response.BucketName = e.Detail.RequestParameters.BucketName
	response.Key = e.Detail.RequestParameters.Key

	//TODO: maybe we can get the filesize with an api call, before downloading it
	// 		extract that to seperate lambda
	//check filesize if > 400MB fail
	/*
		var response Response
		MB := int(numBytes / 1024 / 1024)
		response.FileSize = MB

		if response.FileSize > 400 {
			log.Println("File to large")
			return response, nil //if I fail the lambda, the fail state in sfn won't be executed I guess
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
			Bucket: aws.String("sfnstack-parquet0583a65d-1aklpe9quref8"),
			Key:    &uploadFileName,
			Body:   fr,
		})
		if err != nil {
			log.Println("Failed to upload to S3: " + err.Error())
		}

		log.Printf("DEBUG: Upload finished! location: %s", result.Location)
	*/

	return response, nil
}

func main() {
	lambda.Start(handler)
}
