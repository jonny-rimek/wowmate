package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jonny-rimek/wowmate/services/golib"
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

//Event struct is needed to represent the S3 event structure
type Event struct {
	Detail Detail `json:"detail"`
}

//Response is the object the next step in Stepfunctions expects
type Response struct {
	FileSize   int    `json:"file_size"`
	BucketName string `json:"bucketName"`
	Key        string `json:"key"`
}

//TODO: add event bridge as a step between s3 and sfn

func handler(e Event) (Response, error) {

	//TODO: add conanical log
	log.Print("DEBUG: bucketname: " + e.Detail.RequestParameters.BucketName)
	log.Print("DEBUG: filename: " + e.Detail.RequestParameters.Key)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	svc := s3.New(sess)

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

	return response, nil
}

func main() {
	golib.InitLogging()
	lambda.Start(handler)
}
