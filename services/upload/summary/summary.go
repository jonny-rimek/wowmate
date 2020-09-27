package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	_ "github.com/lib/pq"
)

var ConnStr string

type Event struct {
	Filename string `json:"filename"`
}

func handler(e Event) error {
	log.Println("hello world: " + e.Filename)

	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer db.Close()
	log.Println("openend connection")

	//TODO: check result if 0 rows imported throw an error
	q := fmt.Sprintf(`
		INSERT INTO summary(caster_name, damage, mythicplus_uuid)
		(SELECT
			caster_name, SUM(actual_amount), mythicplus_uuid 
		FROM
			combatlogs
		WHERE
			mythicplus_uuid = '%v'
			AND event_type = 'SPELL_DAMAGE'
			AND caster_id LIKE '%v'
		GROUP BY
			caster_name, mythicplus_uuid
		);
		`, strings.TrimSuffix(e.Filename, ".csv"), "Player-%")

	// log.Println(q)

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
	log.Println("summary successfull")
	return nil
}

func main() {
	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		log.Println("csv bucket env var is empty")
		return
	}

	proxyEndpoint := os.Getenv("DB_ENDPOINT")
	if secretArn == "" {
		log.Println("csv bucket env var is empty")
		return
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Println("failed to create new session")
		return
	}

	ConnStr, err = golib.DBCreds(secretArn, proxyEndpoint, sess)
	if err != nil {
		return
	}

	lambda.Start(handler)
}
