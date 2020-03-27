package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jonny-rimek/wowmate/services/golib"
	"github.com/sirupsen/logrus"
)

//Event is the data from StepFunctions
type Event struct {
	BucketName string `json:"result_bucket"`
	Key        string `json:"file_name"`
}

func handler(e Event) error {
	var errMsg string

	bytes, wcu, dup, rcu, err := handle(e)
	if err == nil {
		errMsg = ""
	} else {
		errMsg = err.Error()
	}

	golib.CanonicalLog(map[string]interface{}{
		"bucket":        e.BucketName,
		"key":           e.Key,
		"downloaded KB": bytes / 1024,
		"duplicate":     dup,
		"wcu":           wcu,
		"rcu":           rcu,
		"err":           errMsg,
	})
	return err
}

func handle(e Event) (int64, float64, bool, float64, error) {
	sess, _ := session.NewSession()
	var bytes int64
	var wcu float64
	var rcu float64
	var duplicate bool

	file, bytes, err := golib.DownloadFileFromS3(e.BucketName, e.Key, sess)
	if err != nil {
		return bytes, wcu, duplicate, rcu, err
	}

	records, err := parseCSV(file)
	if err != nil {
		return bytes, wcu, duplicate, rcu, err
	}

	rcu, err = newCombatlog(records, sess)
	if err != nil {
		if err.Error() == "duplicate combatlog" {
			duplicate = true
		}
		return bytes, wcu, duplicate, rcu, err
	}

	wcu, err = writeDynamoDB(records, sess)
	return bytes, wcu, duplicate, rcu, err
}

func newCombatlog(records []golib.DamageSummary, sess *session.Session) (float64, error) {
	svcdb := dynamodb.New(sess)

	//IMPROVE:
	//use GetItem, should consume less RCU as it is limited to 1 item, atm rcu is still 0.5
	//https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/DynamoDBReadItem.go
	//https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_GetItem.html
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v1": {
				S: aws.String(records[0].Hash),
			},
		},
		KeyConditionExpression: aws.String("pk = :v1"),
		TableName:              aws.String(os.Getenv("DDB_NAME")),
		ReturnConsumedCapacity: aws.String("TOTAL"),
		Limit:                  aws.Int64(1),
	}

	result, err := svcdb.Query(input)
	rcu := *result.ConsumedCapacity.CapacityUnits
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return rcu, err
			case dynamodb.ErrCodeResourceNotFoundException:
				return rcu, err
			case dynamodb.ErrCodeInternalServerError:
				return rcu, err
			default:
				return rcu, err
			}
		} else {
			return rcu, err
		}
	}

	if *result.Count > 0 {
		return rcu, fmt.Errorf("duplicate combatlog")
	}

	return rcu, nil
}

//it handles more than 25 entries, it's not tested tho, might have an
//edge case
//WISHLIST: test behaviour with =25 and >25 entries
func writeDynamoDB(records []golib.DamageSummary, sess *session.Session) (float64, error) {
	writeRequests, err := createDynamoDBWriteRequest(records)
	var writes []*dynamodb.WriteRequest

	var consumedWCU float64
	for _, value := range writeRequests {
		writes = append(writes, value)
		if len(writes) > 25 {
			logrus.Debug("batch size > 25")
			wcu, err := writeBatchDynamoDB(writes, sess)
			if err != nil {
				return consumedWCU, err
			}
			consumedWCU += wcu
			writes = nil
		}
	}
	wcu, err := writeBatchDynamoDB(writes, sess)
	if err != nil {
		return consumedWCU, err
	}
	consumedWCU += wcu

	return consumedWCU, nil
}

func createDynamoDBWriteRequest(records []golib.DamageSummary) ([]*dynamodb.WriteRequest, error) {
	writesRequets := []*dynamodb.WriteRequest{}

	for _, s := range records {
		av, err := dynamodbattribute.MarshalMap(s)
		if err != nil {
			return nil, fmt.Errorf("got error marshalling csv struct into dynamoDB element: %v", err)
		}

		wr := &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: av,
			},
		}
		writesRequets = append(writesRequets, wr)
	}
	return writesRequets, nil
}

func writeBatchDynamoDB(writeRequests []*dynamodb.WriteRequest, sess *session.Session) (float64, error) {
	svcdb := dynamodb.New(sess)
	ddbTableName := os.Getenv("DDB_NAME")

	input := &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]*dynamodb.WriteRequest{
			ddbTableName: writeRequests,
		},
		ReturnConsumedCapacity: aws.String("TOTAL"),
	}

	result, err := svcdb.BatchWriteItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeProvisionedThroughputExceededException, err)
			case dynamodb.ErrCodeResourceNotFoundException:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeResourceNotFoundException, err)
			case dynamodb.ErrCodeItemCollectionSizeLimitExceededException:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeItemCollectionSizeLimitExceededException, err)
			case dynamodb.ErrCodeRequestLimitExceeded:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeRequestLimitExceeded, err)
			case dynamodb.ErrCodeInternalServerError:
				return 0, fmt.Errorf("%v -- %v", dynamodb.ErrCodeInternalServerError, err)
			case dynamodb.ErrCodeTransactionCanceledException:
				return 0, err
			default:
				return 0, fmt.Errorf("default error: %v", err)
			}
		} else {
			return 0, fmt.Errorf("non aws error: %v", err)
		}
	}
	//NOTE: unprocessed items of result are never check, if it is not empty the lambda will
	//		fail and thus alert me, when the case arrises
	//		when does this occur, if I get an error I believe non in the batch got written to DDB
	if len(result.UnprocessedItems) > 0 {
		return 0, fmt.Errorf("handle unprocessed items")
	}
	return *result.ConsumedCapacity[0].CapacityUnits, nil
}

func parseCSV(file []byte) ([]golib.DamageSummary, error) {
	var records []golib.DamageSummary

	reader := bytes.NewReader(file)
	scanner := bufio.NewScanner(reader)

	var s strings.Builder

	scanner.Scan() //skips the first line, which is the header of the csv
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), ",")

		damage, err := strconv.ParseInt(trimQuotes(row[0]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert damage column to int64: %v", err)
		}

		encounterID, err := strconv.Atoi(trimQuotes(row[1]))
		if err != nil {
			return nil, fmt.Errorf("Failed to convert encounter id column to int: %v", err)
		}
		casterID := trimQuotes(row[3])
		casterName := trimQuotes(row[4])

		//strings.Builder is way faster than +=
		s.WriteString(fmt.Sprintf("|%v|%v", casterID, damage))

		r := golib.DamageSummary{
			EncounterID: fmt.Sprintf("%v_b_d_%v", encounterID, casterID),
			Damage:      damage,
			CasterID:    casterID,
			CasterName:  casterName,
		}
		records = append(records, r)
	}
	h := md5.New()
	io.WriteString(h, s.String())
	hash := fmt.Sprintf("%x%x", h.Sum(nil), records[0].EncounterID)

	for i := 0; i < len(records); i++ {
		//limiting the size of the pk to 12 chars
		//to reduce the storage and not have an unessesary long value
		//as they are exposed to the user in the URL
		records[i].Hash = hash[:12]
	}

	return records, nil
}

func trimQuotes(input string) string {
	output := strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

func main() {
	golib.InitLogging()
	lambda.Start(handler)
}
