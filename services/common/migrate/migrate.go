package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jonny-rimek/wowmate/services/common/golib"
)
func handler() error {
	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		return fmt.Errorf("secret arn env var is empty")
	}

	sess, err := session.NewSession()
	if err != nil {
		return err
	}
	
	connStr, err := golib.DBCreds(secretArn, "", sess)
	if err != nil {
		return err
	}

	//checks all files in the sql directory and applies all the missing ones
	//the current state is saved inside the db
	m, err := migrate.New("file://sql", connStr)
	if err != nil {
		return err
	}

	//IMPROVE: as far as I can tell, it doesn't automatically close the db connection
	//			it's not that big a problem, because I don't run that many migrations
	//			but I should keep an eye on it
	//			it will be even less of a problem with the db proxy
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			//don't fail only because there was no change
			//I run the migration after every push
			log.Println("no change")
			return nil
		}
		return err
	}

	log.Println("migration successful")
	//TODO: print query excecuted
	return nil
}

func main() {
	lambda.Start(handler)
}
