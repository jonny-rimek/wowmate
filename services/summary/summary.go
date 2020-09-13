package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

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

type Event struct {
	Filename string `json:"filename"`
}

func handler(e Event) error {
	log.Println("hello world: " + e.Filename)

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

	q := fmt.Sprintf(`
			INSERT INTO summary(caster_name, damage)
			(SELECT
				caster_name, SUM(actual_amount) AS damage
			FROM
				combatlogs
			WHERE
				upload_uuid = '%v'
				AND event_type = 'SPELL_DAMAGE'
				AND caster_id LIKE 'Player-%'
			GROUP BY
				caster_name
			);
			`, strings.TrimSuffix(e.Filename, ".csv"))

	rows, err := db.Query(q)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var s string

		err = rows.Scan(&s)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		log.Printf("query successfull: %v", s)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
