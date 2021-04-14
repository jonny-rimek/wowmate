package golib

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/sirupsen/logrus"
)

// TimestreamQuery runs a query against timestream and checks if the query is already finished.
// timestream returns a response after ~6 seconds even if it is not finished. It checks the NextToken
// and reruns the query if need be
func TimestreamQuery(ctx aws.Context, query *string, querySvc *timestreamquery.TimestreamQuery) (*timestreamquery.QueryOutput, error) {
	var result *timestreamquery.QueryOutput
	var err error

	if os.Getenv("LOCAL") == "true" {
		result, err = querySvc.Query(&timestreamquery.QueryInput{
			QueryString: query,
		})
	} else {
		result, err = querySvc.QueryWithContext(ctx, &timestreamquery.QueryInput{
			QueryString: query,
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed querying timestream: %v", err.Error())
	}
	logrus.Debug(result)
	logrus.Debug("after query")

	// TODO: refactor^^
	if result.NextToken != nil {
		result, err = querySvc.QueryWithContext(ctx, &timestreamquery.QueryInput{
			QueryString: query,
			NextToken:   result.NextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed querying timestream: %v", err.Error())
		}
	}
	if result.NextToken != nil {
		result, err = querySvc.QueryWithContext(ctx, &timestreamquery.QueryInput{
			QueryString: query,
			NextToken:   result.NextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed querying timestream: %v", err.Error())
		}
	}
	if result.NextToken != nil {
		result, err = querySvc.QueryWithContext(ctx, &timestreamquery.QueryInput{
			QueryString: query,
			NextToken:   result.NextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed querying timestream: %v", err.Error())
		}
	}

	return result, nil
}

/*
func TimestreamQuery2(query *string, client *timestreamquery2.Client) (*timestreamquery2.QueryOutput, error) {
	queryInput := timestreamquery2.QueryInput{
		QueryString: query,
	}

	//returns operation error Timestream Query: Query, https response error StatusCode: 404,
	//RequestID: ed6f01e1-bf16-426e-94b3-14a18aba2bbf, api error UnknownOperationException: UnknownError: OperationError
	result, err := client.Query(context.TODO(), &queryInput)
	if err != nil {
		return nil, err
	}
	if len(result.Rows) == 0 {
		return nil, fmt.Errorf("query returned empty result")
	}

	var result *timestreamquery2.QueryOutput
	logrus.Debug("after query")

	return result, nil
}
*/

// WriteToTimestream takes a slice of records to write and writes batches of 100 records to timestream
func WriteToTimestream(ctx aws.Context, writeSvc *timestreamwrite.TimestreamWrite, e *timestreamwrite.WriteRecordsInput) error {
	e.DatabaseName = aws.String("wowmate-analytics")
	e.TableName = aws.String("combatlogs")

	var err error
	if os.Getenv("LOCAL") == "true" {
		_, err = writeSvc.WriteRecords(e)
	} else {
		_, err = writeSvc.WriteRecordsWithContext(ctx, e)
	}
	if err != nil {
		// debug info
		// prettyStruct, err := PrettyStruct(e)
		// if err != nil {
		// 	return fmt.Errorf("failed to to get pretty struct: %v", err.Error())
		// }
		// log.Println(prettyStruct)
		log.Printf("failed to write to timestream, logging: %v", err)
		// I only handle one error from all goroutines, that's why I'm logging
		// each to see if errors occur in multiple goroutines

		return fmt.Errorf("failed to write to timestream, returning: %v", err)
	}
	return nil
}

/*
func UploadToTimestream2(cfg aws2.Config, e []types2.Record) error {
	log.Printf("%v dmg records", len(e))


	client := timestreamwrite2.NewFromConfig(cfg)


	for i := 0; i < len(e); i += 100 {

		//get the upper bound of the record to write, in case it is the
		//last bit of records and i + 99 does not exist
		j := 0
		if i+99 > len(e) {
			j = len(e)
		} else {
			j = i + 99
		}

		//use common batching https://docs.aws.amazon.com/timestream/latest/developerguide/metering-and-pricing.writes.html#metering-and-pricing.writes.write-size-multiple-events

		input := timestreamwrite2.WriteRecordsInput{
			DatabaseName: aws2.String("wowmate-analytics"), //TODO: dont hardcode
			Records:      e[i:j],
			TableName:    aws2.String("combatlogs"),
			//CommonAttributes: nil, TODO:
		}
		_, err := client.WriteRecords(context.TODO(), &input)
		if err != nil {
			return err
		}
	}
	log.Println("Write records to timestream was successful")
	return nil
}
*/

