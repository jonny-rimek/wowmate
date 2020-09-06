package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//DatabasesCredentials are the data to log into the db
type DatabasesCredentials struct {
	DatabaseName string `json:"dbname"`
	Password     string `json:"password"`
	UserName     string `json:"username"`
	Host         string `json:"host"`
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	serverError := events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
	}

	secretArn := os.Getenv("SECRET_ARN")

	sess, err := session.NewSession()
	if err != nil {
		log.Println(err.Error())
		log.Println("failed to create new session")
		return serverError, err
	}

	//TODO: should move get secret outside of handler, because it dosn't need to run on every invocation
	svc := secretsmanager.New(sess)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretArn),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeDecryptionFailure:
				// Secrets Manager can't decrypt the protected secret text using the provided KMS key.
				log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				return serverError, err

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				return serverError, err

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				return serverError, err

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				return serverError, err

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				return serverError, err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
			return serverError, err
		}
	}

	var creds = DatabasesCredentials{}
	err = json.Unmarshal([]byte(*result.SecretString), &creds)
	if err != nil {
		log.Println(err.Error())
		return serverError, err
	}

	connStr := fmt.Sprintf(
		"user=%v dbname=%v sslmode=verify-full host=%v password=%v port=5432",
		creds.UserName,
		creds.DatabaseName,
		// creds.Host,
		"dbproxy.proxy-cvevooy4lacx.us-east-1.rds.amazonaws.com",
		creds.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err.Error())
		return serverError, err
	}
	defer db.Close()
	log.Println("openend connection")

	// q := fmt.Sprintf(`SELECT COUNT(*) FROM combatlogs;`)

	rows, err := db.Query(`SELECT COUNT(*) FROM combatlogs;`)
	if err != nil {
		log.Println(err.Error())
		return serverError, err
		// return err
	}

	defer rows.Close()

	var s string
	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			log.Println(err.Error())
			return serverError, err
			// return err
		}
		log.Printf("import query successfull: %v", s)
	}
	resp := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       s,
	}
	return resp, nil
}

func main() {
	lambda.Start(handler)
}
