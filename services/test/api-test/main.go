package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/synthetics"
)

var svc *synthetics.Synthetics

func main() {
	var sess *session.Session
	var err error

	sess, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	svc = synthetics.New(sess)

	err = startCanary()
	if err != nil {
		handleError(err)
	}
}

func startCanary() error {
	_, err := svc.StartCanary(&synthetics.StartCanaryInput{
		Name: aws.String("wmpreprodcanaryf2cb3d"),
	})
	if err != nil {
		return err
	}
	return nil
}

func handleError(err error) {
	log.Printf("%s", err)
	os.Exit(1)
}
