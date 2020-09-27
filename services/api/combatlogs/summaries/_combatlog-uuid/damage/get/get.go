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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	_ "github.com/lib/pq"
)

var ConnStr string

type DamageResult struct {
	PlayerName string `json:"player_name"`
	Damage     int    `json:"damage"`
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	serverError := events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
	}

	log.Println(request.PathParameters["combatlog_uuid"])

	db, err := sql.Open("postgres", ConnStr)
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
	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		log.Println("failed csv bucket env var is empty")
		return
	}

	dbEndpoint := os.Getenv("DB_ENDPOINT")
	if dbEndpoint == "" {
		log.Println("failed db endpoint env var is empty")
		return
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Println("failed to create new session")
		return
	}

	ConnStr, err = golib.DBCreds(secretArn, "", sess)
	if err != nil {
		log.Println("failed to get the db creds")
		return
	}

	lambda.Start(handler)
}
