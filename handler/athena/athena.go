package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

//RequestParameters test ..
type RequestParameters struct {
	ID           string `json:"id"`
	ResultBucket string `json:"result_bucket"`
	Query        string `json:"query"`
	Region       string `json:"region"`
	Table        string `json:"table"`
	FileName     string `json:"file_name"`
}

func handler(es []RequestParameters) (RequestParameters, error) {
	e := es[0]

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(e.Region)},
	)
	svc := athena.New(sess)

	var sq athena.StartQueryExecutionInput
	sq.SetQueryString(e.Query)

	var q athena.QueryExecutionContext
	q.SetDatabase(e.Table)
	sq.SetQueryExecutionContext(&q)

	var rc athena.ResultConfiguration
	rc.SetOutputLocation("s3://" + e.ResultBucket + "/results/")
	sq.SetResultConfiguration(&rc)

	queryID, err := svc.StartQueryExecution(&sq)
	if err != nil {
		fmt.Errorf("Athena query failed: " + err.Error())
		return e, err
	}
	log.Println("DEBUG: athena query started")

	e.FileName = "results/" + *queryID.QueryExecutionId + ".csv"
	log.Println("DEBUG: filename: " + e.FileName)

	e.ID = *queryID.QueryExecutionId

	return e, nil
}

func main() {
	lambda.Start(handler)
}
