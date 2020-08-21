package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	_ "github.com/lib/pq"
)

func handle() {
	secretArn := os.Getenv("SECRET_ARN")
	_ = os.Getenv("CSV_BUCKET_NAME")

	log.Println("such invovation")

	sess, err := session.NewSession()
	if err != nil {
		log.Println("failed to create new session")
		return
	}

	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretArn),
		VersionStage: aws.String("lol"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		log.Println("failed to get the secret")
		return
	}

	log.Println(*result.SecretString)
	//TODO: get connecetion data from secretsmanager
	// connStr := "user= dbname= sslmode= host= password="
	/*
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Query(`
			SHOW TABLES;
		`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("finished")
	*/
}

func main() {
	lambda.Start(handle)

	//TODO: connect to db with hardcoded creds and run SELECT 1
}
