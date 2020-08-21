package main

import (
	"database/sql"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

func handle() {
	log.Println("such invovation")

	//TODO: get connecetion data from secretsmanager
	// connStr := "user= dbname= sslmode= host= password="
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
}

func main() {
	lambda.Start(handle)

	//TODO: connect to db with hardcoded creds and run SELECT 1
}
