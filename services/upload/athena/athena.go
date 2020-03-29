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
	UploadUUID string `json:"upload_uuid"`
	Year       int    `json:"year"`
	Month      int    `json:"month"`
	Day        int    `json:"day"`
	Hour       int    `json:"hour"`
	Minute     int    `json:"minute"`
}

//Response is the object the next step in Stepfunctions expects
type Response struct {
	BucketName string `json:"result_bucket"`
	Key        string `json:"file_name"`
	ID         string `json:"id"`
}

func handler(e SfnEvent) (Response, error) {
	region := os.Getenv("REGION")
	resultBucket := os.Getenv("RESULT_BUCKET")
	athenaDatabase := os.Getenv("ATHENA_DATABASE")
	resp := Response{BucketName: resultBucket}

	query := fmt.Sprintf(`
		SELECT cl.damage, ei.encounter_id, cl.boss_fight_uuid, cl.caster_id, cl.caster_name, cl.upload_uuid, cl.partition_0, cl.partition_1, cl.partition_2, cl.partition_3, cl.partition_4
		FROM (
			SELECT SUM(actual_amount) AS damage, caster_name, caster_id, boss_fight_uuid, upload_uuid, partition_0, partition_1, partition_2, partition_3, partition_4
			FROM  "wowmate"."combatlogs"
			WHERE caster_type LIKE '0x5%%' AND caster_name != 'nil' 
			GROUP BY caster_name, caster_id, boss_fight_uuid, upload_uuid, partition_0, partition_1, partition_2, partition_3,partition_4
			) AS cl
		JOIN (
			SELECT encounter_id, boss_fight_uuid, upload_uuid
			FROM "wowmate"."combatlogs"
			WHERE event_type = 'ENCOUNTER_START'
			GROUP BY encounter_id, boss_fight_uuid, upload_uuid) AS ei
			ON cl.upload_uuid = ei.upload_uuid
		WHERE 
			cl.upload_uuid = '%v'
			AND cl.partition_0 = '%v'
			AND cl.partition_1 = '%v'
			AND cl.partition_2 = '%v'
			AND cl.partition_3 = '%v'
			AND cl.partition_4 = '%v'
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
	rc.SetOutputLocation("s3://" + resultBucket + "/results/")
	sq.SetResultConfiguration(&rc)

	queryID, err := svc.StartQueryExecution(&sq)
	if err != nil {
		return resp, err
	}
	log.Println("DEBUG: athena query started")

	resp.ID = *queryID.QueryExecutionId
	resp.Key = "results/" + *queryID.QueryExecutionId + ".csv"
	log.Println("DEBUG: filename: " + resp.Key)

	return resp, nil
}

func main() {
	lambda.Start(handler)
}
