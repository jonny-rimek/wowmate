package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

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

func handler() error {
	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		return fmt.Errorf("csv bucket env var is empty")
	}

	proxyEndpoint := os.Getenv("DB_ENDPOINT")
	if secretArn == "" {
		return fmt.Errorf("csv bucket env var is empty")
	}

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println(err.Error())
		log.Println("failed to create new session")
		return err
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
				fmt.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				return err

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				return err

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				return err

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				return err

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				return err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
			return err
		}
	}

	var creds = DatabasesCredentials{}
	err = json.Unmarshal([]byte(*result.SecretString), &creds)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	connStr := fmt.Sprintf(
		"user=%v dbname=%v sslmode=verify-full host=%v password=%v port=5432",
		creds.UserName,
		creds.DatabaseName,
		// creds.Host,
		proxyEndpoint,
		creds.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()
	log.Println("openend connection")

	ystd := time.Now().AddDate(0, 0, -1)
	td := time.Now()
	tmrw := time.Now().AddDate(0, 0, 1)

	const (
		layoutISO   = "2006-01-02"
		layoutTable = "2006_01_02"
	)

	q := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS combatlogs_%v PARTITION OF combatlogs 
			FOR VALUES FROM ('%v 00:00:00') TO ('%v 23:59:59');
		CREATE TABLE IF NOT EXISTS combatlogs_%v PARTITION OF combatlogs 
			FOR VALUES FROM ('%v 00:00:00') TO ('%v 23:59:59');
		CREATE TABLE IF NOT EXISTS combatlogs_%v PARTITION OF combatlogs 
			FOR VALUES FROM ('%v 00:00:00') TO ('%v 23:59:59');
		`,
		ystd.Format(layoutTable),
		ystd.Format(layoutISO),
		ystd.Format(layoutISO),
		td.Format(layoutTable),
		td.Format(layoutISO),
		td.Format(layoutISO),
		tmrw.Format(layoutTable),
		tmrw.Format(layoutISO),
		tmrw.Format(layoutISO),
	)

	log.Println(q)

	rows, err := db.Query(q)
	if err != nil {
		log.Println("query" + err.Error())
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var s string

		err = rows.Scan(&s)
		if err != nil {
			fmt.Println("rows scan" + err.Error())
			return err
		}
		log.Printf("query successfull: %v", s)
	}
	log.Println("partition successfully created")
	return nil
}

func main() {
	lambda.Start(handler)
}
