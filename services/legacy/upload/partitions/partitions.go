package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/athena"
)

func addPartition() error {
	awscfg := &aws.Config{}
	awscfg.WithRegion("us-east-1")
	sess := session.Must(session.NewSession(awscfg))
	svc := athena.New(sess, aws.NewConfig().WithRegion("us-east-1"))

	athenaDB := os.Getenv("ATHENA_DB")
	athenaTable := os.Getenv("ATHENA_TABLE")
	athenaQueryResultBucket := os.Getenv("ATHENA_QUERY_RESULT_BUCKET")
	bucket := os.Getenv("SOURCE_BUCKET")

	today := time.Now().UTC()
	query := "ALTER TABLE " + athenaTable + " ADD IF NOT EXISTS"

	//NOTE: I'm always adding partitions that already exist, because I add 2 hours
	//		but run the cron hourly, this is just to be save, I don't see a downside
	//		so there is no need to optimize
	for hour := 0; hour < 3; hour++ {
		for minute := 0; minute < 60; minute++ {
			q := fmt.Sprintf(`
				PARTITION (partition_0 = '%d', partition_1 = '%d', partition_2 = '%d', partition_3 = '%d', partition_4 = '%d')
				LOCATION 's3://%s/%d/%d/%d/%d/%d' 
				`, today.Year(), today.Month(), today.Day(), today.Hour()+hour, minute, bucket, today.Year(), today.Month(), today.Day(), hour, minute)
			query += q
		}
	}

	err := queryAthena(query, "s3://"+athenaQueryResultBucket, athenaDB, svc)
	if err != nil {
		return err
	}
	log.Println("Started add partition query")
	return nil
}

//TODO: check result and alert on failure
func queryAthena(query string, outputLocation string, athenaDB string, svc *athena.Athena) error {
	var s athena.StartQueryExecutionInput
	s.SetQueryString(query)

	var q athena.QueryExecutionContext
	q.SetDatabase(athenaDB)
	s.SetQueryExecutionContext(&q)

	var r athena.ResultConfiguration
	r.SetOutputLocation(outputLocation)
	s.SetResultConfiguration(&r)

	_, err := svc.StartQueryExecution(&s)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(addPartition)
}
