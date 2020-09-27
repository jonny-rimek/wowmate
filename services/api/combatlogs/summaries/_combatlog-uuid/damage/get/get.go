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
	_ "github.com/lib/pq"
)

//DatabasesCredentials are the data to log into the db
type DatabasesCredentials struct {
	DatabaseName string `json:"dbname"`
	Password     string `json:"password"`
	UserName     string `json:"username"`
	Host         string `json:"host"`
}

type DamageResult struct {
	PlayerName string `json:"player_name"`
	Damage     int    `json:"damage"`
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	serverError := events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
	}
	log.Println(request.PathParameters["combatlog_uuid"])

	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		return serverError, fmt.Errorf("csv bucket env var is empty")
	}

	dbEndpoint := os.Getenv("DB_ENDPOINT")
	if dbEndpoint == "" {
		return serverError, fmt.Errorf("db endpoint env var is empty")
	}

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
		dbEndpoint,
		creds.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err.Error())
		return serverError, err
	}
	defer db.Close()
	log.Println("openend connection")

	rows, err := db.Query(`SELECT caster_name, damage FROM summary WHERE combatlog_uuid = $1;`, request.PathParameters["combatlog_uuid"])
	if err != nil {
		log.Println(err.Error())
		return serverError, err
	}

	defer rows.Close()

	var r []DamageResult
	var s string
	var i int
	for rows.Next() {
		err = rows.Scan(&s, &i)
		if err != nil {
			log.Println(err.Error())
			return serverError, err
		}
		res := DamageResult{
			PlayerName: s,
			Damage:     i,
		}
		r = append(r, res)
	}
	log.Printf("query successfull")

	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return serverError, err
	}

	resp := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(b),
	}
	return resp, nil
}

func main() {
	lambda.Start(handler)
}