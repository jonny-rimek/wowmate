package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handle() {
	log.Println("such invovation")
}

func main() {
	lambda.Start(handle)
}
