package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
	"github.com/jonny-rimek/wowmate/services/golib"
)

//RequestParameters test ..
type RequestParameters struct {
	ID           string `json:"id"`
	ResultBucket string `json:"result_bucket"`
	FileName     string `json:"file_name"`
}

func handler(e RequestParameters) (RequestParameters, error) {

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := athena.New(sess, aws.NewConfig().WithRegion("eu-central-1"))

	var qri athena.GetQueryExecutionInput
	qri.SetQueryExecutionId(e.ID)
	var qrop *athena.GetQueryExecutionOutput

	qrop, err := svc.GetQueryExecution(&qri)
	if err != nil {
		log.Println("Failed to get status on the query: " + err.Error())
		return e, err
	}

	if *qrop.QueryExecution.Status.State == "RUNNING" {
		err = fmt.Errorf("Query is still running")
		return e, err

	} else if *qrop.QueryExecution.Status.State == "SUCCEEDED" {

		var ip athena.GetQueryResultsInput
		ip.SetQueryExecutionId(e.ID)

		op, err := svc.GetQueryResults(&ip)
		if err != nil {
			log.Println(err)
			return e, nil
		}
		log.Printf("%+v", op)
	} else {
		log.Println(*qrop.QueryExecution.Status.State)
		log.Println(*qrop.QueryExecution.Status.StateChangeReason)
	}

	return e, nil
}

func main() {
	golib.InitLogging()
	lambda.Start(handler)
}
