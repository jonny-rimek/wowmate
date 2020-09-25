package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

//SfnEvent provides config data to the lambda
type SfnEvent struct {
	UploadUUID    string `json:"upload_uuid"`
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	Day           int    `json:"day"`
	Hour          int    `json:"hour"`
	Minute        int    `json:"minute"`
	ParquetBucket string `json:"parquet_bucket"`
	ParquetFile   string `json:"parquet_file"`
}

//Response is the object the next step in Stepfunctions expects
type Response struct {
	BucketName    string `json:"result_bucket"`
	Key           string `json:"file_name"`
	ID            string `json:"id"`
	ParquetBucket string `json:"parquet_bucket"`
	ParquetFile   string `json:"parquet_file"`
}

func handler(e SfnEvent) (Response, error) {
	region := os.Getenv("REGION")
	resultBucket := os.Getenv("RESULT_BUCKET")
	athenaDatabase := os.Getenv("ATHENA_DATABASE")
	resp := Response{
		BucketName:    resultBucket,
		ParquetBucket: e.ParquetBucket,
		ParquetFile:   e.ParquetFile,
	}

	query := fmt.Sprintf(`
		SELECT cl.damage, ei.encounter_id, cl.boss_fight_uuid, cl.caster_id, cl.caster_name, cl.upload_uuid, cl.year, cl.month, cl.day, cl.hour, cl.minute
		FROM (
			SELECT SUM(actual_amount) AS damage, caster_name, caster_id, boss_fight_uuid, upload_uuid, year, month, day, hour, minute
			FROM  "wowmate"."combatlogs"
			WHERE caster_type LIKE '0x5%%' AND caster_name != 'nil' 
			GROUP BY caster_name, caster_id, boss_fight_uuid, upload_uuid, year, month, day, hour,minute
			) AS cl
		JOIN (
			SELECT encounter_id, boss_fight_uuid, upload_uuid
			FROM "wowmate"."combatlogs"
			WHERE event_type = 'ENCOUNTER_START'
			GROUP BY encounter_id, boss_fight_uuid, upload_uuid) AS ei
			ON cl.upload_uuid = ei.upload_uuid
		WHERE 
			cl.upload_uuid = '%v'
			AND cl.year = '%v'
			AND cl.month = '%v'
			AND cl.day = '%v'
			AND cl.hour = '%v'
			AND cl.minute = '%v'
		ORDER BY encounter_id, damage DESC
	`,
		//TODO: das query selbst gibt atm keine daten zurück
		e.UploadUUID,
		e.Year,
		e.Month,
		e.Day,
		e.Hour,
		e.Minute,
	)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	svc := athena.New(sess)

	var sq athena.StartQueryExecutionInput
	sq.SetQueryString(query)

	var q athena.QueryExecutionContext
	q.SetDatabase(athenaDatabase)
	sq.SetQueryExecutionContext(&q)

	var rc athena.ResultConfiguration
	rc.SetOutputLocation("s3://" + resultBucket + "/damage-summaries/")
	sq.SetResultConfiguration(&rc)

	queryID, err := svc.StartQueryExecution(&sq)
	if err != nil {
		return resp, err
	}
	log.Println("DEBUG: athena query started")

	resp.ID = *queryID.QueryExecutionId
	resp.Key = "damage-summaries/" + *queryID.QueryExecutionId + ".csv"
	log.Println("DEBUG: filename: " + resp.Key)

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
