package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

//SfnEvent provides config data to the lambda
type SfnEvent struct {
	UploadUUID    string `json:"upload_uuid"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	Day           int    `json:"day"`
	Hour          int    `json:"hour"`
	Minute        int    `json:"minute"`
	ParquetBucket string `json:"parquet_bucket"`
	ParquetFile   string `json:"parquet_file"`
}

//Response is the object the next step in Stepfunctions expects
type Response struct {
	BucketName    string `json:"result_bucket"`
	Key           string `json:"file_name"`
	ID            string `json:"id"`
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
	region := os.Getenv("REGION")
	resultBucket := os.Getenv("RESULT_BUCKET")
	athenaDatabase := os.Getenv("ATHENA_DATABASE")
	resp := Response{BucketName: resultBucket}

	//IMPROVE: don't hardcode table name
	query := "MSCK REPAIR TABLE combatlogs;"

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	svc := athena.New(sess)

	var sq athena.StartQueryExecutionInput
	sq.SetQueryString(query)

	var q athena.QueryExecutionContext
	q.SetDatabase(athenaDatabase)
	sq.SetQueryExecutionContext(&q)

	var rc athena.ResultConfiguration
	rc.SetOutputLocation("s3://" + resultBucket + "/repair/")
	sq.SetResultConfiguration(&rc)

	queryID, err := svc.StartQueryExecution(&sq)
	if err != nil {
		return resp, err
	}
	log.Println("DEBUG: log run: " + query)
	log.Println("DEBUG: athena query started")

	resp.ID = *queryID.QueryExecutionId
	resp.Key = "repair/" + *queryID.QueryExecutionId + ".csv"
	resp.UploadUUID = e.UploadUUID
	resp.Year = e.Year
	resp.Month = e.Month
	resp.Day = e.Day
	resp.Hour = e.Hour
	resp.Minute = e.Minute
	resp.ParquetBucket = e.ParquetBucket
	resp.ParquetFile = e.ParquetFile
	//TODO: add input to output
	log.Println("DEBUG: filename: " + resp.Key)

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
