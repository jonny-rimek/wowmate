package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jonny-rimek/wowmate/services/common/golib"
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
	connStr, err := golib.DBCreds(secretArn, proxyEndpoint, sess)
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()
	log.Println("openend connection")

	yesterday := time.Now().AddDate(0, 0, -1)
	today := time.Now()
	tomorrow := time.Now().AddDate(0, 0, 1)
	monthAgo := time.Now().AddDate(0, 0, -30)

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

		DROP TABLE IF EXISTS combatlogs_%v;
		`,
		yesterday.Format(layoutTable),
		yesterday.Format(layoutISO),
		yesterday.Format(layoutISO),
		today.Format(layoutTable),
		today.Format(layoutISO),
		today.Format(layoutISO),
		tomorrow.Format(layoutTable),
		tomorrow.Format(layoutISO),
		tomorrow.Format(layoutISO),
		monthAgo.Format(layoutTable),
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
