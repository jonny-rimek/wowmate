package golib

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

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

//Secret gets the secret duh
func Secret(secretArn string, sess *session.Session) (string, error) {
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
				return "", err

			case secretsmanager.ErrCodeInternalServiceError:
				// An error occurred on the server side.
				log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				return "", err

			case secretsmanager.ErrCodeInvalidParameterException:
				// You provided an invalid value for a parameter.
				log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				return "", err

			case secretsmanager.ErrCodeInvalidRequestException:
				// You provided a parameter value that is not valid for the current state of the resource.
				log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				return "", err

			case secretsmanager.ErrCodeResourceNotFoundException:
				// We can't find the resource that you asked for.
				log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				return "", err
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
			return "", err
		} 
	}

	return *result.SecretString, nil
}

func DBCreds(secretArn string, host string, sess *session.Session) (string, error) {
	secret, err := Secret(secretArn, sess)
	if err != nil {
		return "", err
	}

	//TODO: should move get secret outside of handler, because it dosn't need to run on every invocation
	var creds = DatabasesCredentials{}
	err = json.Unmarshal([]byte(secret), &creds)
	if err != nil {
		return "", err
	}

	//if no host is provided take the host from the secret
	if host == "" {
		host = creds.Host
	}

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=require",
		creds.UserName,
		url.PathEscape(creds.Password),
		host,
		5432,
		creds.DatabaseName,
	)

	return connStr, nil
}
