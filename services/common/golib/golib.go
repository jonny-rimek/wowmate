package golib

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/sirupsen/logrus"
)

// ++++++++++++++++++++++++
// PlayerDamageDone Structs
// ++++++++++++++++++++++++

// DamagePerSpell is a part of PlayerDamageDone and contains the breakdown of damage per spell
type DamagePerSpell struct {
	SpellID   int    `json:"spell_id"`
	SpellName string `json:"spell_name"`
	Damage    int64  `json:"damage"`
}

// PlayerDamageDone contains player and damage per spell info for the log specific view
type PlayerDamageDone struct {
	Damage         int64            `json:"damage"`
	DamagePerSpell []DamagePerSpell `json:"damage_per_spell"`
	Name           string           `json:"player_name"`
	PlayerID       string           `json:"player_id"`
	Class          string           `json:"class"`
	Specc          string           `json:"specc"`
	/*
	   TODO:
	   	ItemLevel
	   	Covenant
	   	Traits
	   	Conduits
	   	Legendaries
	   	Trinkets
	   	Talents
	*/
}

// DynamoDBPlayerDamageDone is used to save player damage done to dynamodb, log specific view
type DynamoDBPlayerDamageDone struct {
	Pk            string             `json:"pk"`
	Sk            string             `json:"sk"`
	Damage        []PlayerDamageDone `json:"player_damage"`
	Duration      string             `json:"duration"`
	Deaths        int                `json:"deaths"`
	Affixes       string             `json:"affixes"`
	Keylevel      int                `json:"keylevel"`
	DungeonName   string             `json:"dungeon_name"`
	DungeonID     int                `json:"dungeon_id"`
	CombatlogUUID string             `json:"combatlog_uuid"`
	Finished      bool               `json:"finished"`
}

// ++++++++++++++++++++++++
// Keys Structs
// ++++++++++++++++++++++++

// PlayerDamage contains player and damage info for the top keys view etc.
type PlayerDamage struct {
	Damage   int    `json:"damage"` // TODO: convert to int64
	Name     string `json:"player_name"`
	PlayerID string `json:"player_id"`
	Class    string `json:"class"`
	Specc    string `json:"specc"`
}

type DynamoDBKeys struct {
	Pk            string         `json:"pk"`
	Sk            string         `json:"sk"`
	Damage        []PlayerDamage `json:"player_damage"`
	Gsi1pk        string         `json:"gsi1pk"`
	Gsi1sk        string         `json:"gsi1sk"`
	Duration      string         `json:"duration"`
	Deaths        int            `json:"deaths"`
	Affixes       string         `json:"affixes"`
	Keylevel      int            `json:"keylevel"`
	DungeonName   string         `json:"dungeon_name"`
	DungeonID     int            `json:"dungeon_id"`
	CombatlogUUID string         `json:"combatlog_uuid"`
	Finished      bool           `json:"finished"`
}

type JSONKeysResponse struct {
	Data      []JSONKeys `json:"data"`
	FirstSK   string     `json:"first_sk"`
	LastSK    string     `json:"last_sk"`
	FirstPage bool       `json:"first_page"`
	LastPage  bool       `json:"last_page"`
}

/*
type JSONPlayerDamageSimpleResponse2 struct {
	Data      []JSONKeys `json:"data"`
	FirstSK   types2.AttributeValue    `json:"first_sk"`
	LastSK    types2.AttributeValue    `json:"last_sk"`
	FirstPage bool                     `json:"first_page"`
	LastPage  bool                     `json:"last_page"`
}
*/

type JSONKeys struct {
	Damage        []PlayerDamage `json:"player_damage"`
	Duration      string         `json:"duration"`
	Deaths        int            `json:"deaths"`
	Affixes       string         `json:"affixes"`
	Keylevel      int            `json:"keylevel"`
	DungeonName   string         `json:"dungeon_name"`
	DungeonID     int            `json:"dungeon_id"`
	CombatlogUUID string         `json:"combatlog_uuid"`
}

// AGW404 returns a error agi gw v2 response
func AGW404() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Body:       http.StatusText(404),
	}
}

// AGW400 returns a error agi gw v2 response
func AGW400() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 400,
		Body:       http.StatusText(400),
	}
}

// AGW500 returns a error agi gw v2 response
func AGW500() events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       http.StatusText(500),
	}
}

// AGW200 returns a agi gw v2 success response
func AGW200(body string, headers map[string]string) events.APIGatewayV2HTTPResponse {
	return events.APIGatewayV2HTTPResponse{
		Headers:    headers,
		StatusCode: 200,
		Body:       body,
	}
}

type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

// SQSEvent is all the data that gets passed into the lambda from the q
// IMPROVE: use events.SQSEvent, see summary lambda
type SQSEvent struct {
	Records []struct {
		MessageID     string `json:"messageId"`
		ReceiptHandle string `json:"receiptHandle"`
		Body          string `json:"body"`
		Attributes    struct {
			ApproximateReceiveCount          string `json:"ApproximateReceiveCount"`
			SentTimestamp                    string `json:"SentTimestamp"`
			SenderID                         string `json:"SenderId"`
			ApproximateFirstReceiveTimestamp string `json:"ApproximateFirstReceiveTimestamp"`
		} `json:"attributes"`
		EventSource    string `json:"eventSource"`
		EventSourceARN string `json:"eventSourceARN"`
		AwsRegion      string `json:"awsRegion"`
	} `json:"Records"`
}

// DynamoDBQuery is a helper to simplify querying a dynamo db table
func DynamoDBQuery(ctx aws.Context, svc *dynamodb.DynamoDB, input dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {

	result, err := svc.QueryWithContext(ctx, &input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/*
//DynamoDBQuery2 is a helper to simplify querying a dynamo db table
func DynamoDBQuery2(client *dynamodb2.Client, input dynamodb2.QueryInput) (*dynamodb2.QueryOutput, error) {

	result, err := client.Query(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
*/

func DynamoDBPutItem(ctx aws.Context, svc *dynamodb.DynamoDB, ddbTableName *string, record interface{}) (*dynamodb.PutItemOutput, error) {
	av, err := dynamodbattribute.MarshalMap(record)
	if err != nil {
		return nil, err
	}

	var response *dynamodb.PutItemOutput
	if os.Getenv("LOCAL") == "true" {
		response, err = svc.PutItem(&dynamodb.PutItemInput{
			Item:                   av,
			TableName:              ddbTableName,
			ReturnConsumedCapacity: aws.String("TOTAL"),
		})
	} else {
		response, err = svc.PutItemWithContext(ctx, &dynamodb.PutItemInput{
			Item:                   av,
			TableName:              ddbTableName,
			ReturnConsumedCapacity: aws.String("TOTAL"),
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to put item to dynamodb: %v", err)
	}

	return response, nil
}

/*
func DynamoDBPutItem2(client *dynamodb2.Client, ddbTableName *string, record interface{}) (*dynamodb2.PutItemOutput, error) {
	av, err := attributevalue2.MarshalMap(record)
	if err != nil {
		return nil, err
	}

	response, err := client.PutItem(context.TODO(), &dynamodb2.PutItemInput{
		Item:                   av,
		TableName:              ddbTableName,
		ReturnConsumedCapacity: types2.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
*/

func KeysResponseToJson(result *dynamodb.QueryOutput, sorted, firstPage bool) (string, error) {
	var items []DynamoDBKeys
	var lastPage bool // controls if the next btn is shown

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return "", err
	}
	if result.LastEvaluatedKey != nil {
		items = items[:len(items)-1] // we don't care about the last item
		// getting 1 more item than we need and truncaten it off allows us to circumvent a dynamodb2 bug where
		// if 5 items are left to return and you get exactly return 5 items the last evaluated key is not null,
		// which leads to the next request being completely empty
		lastPage = false
	} else if result.LastEvaluatedKey == nil && sorted == true {
		firstPage = true
		lastPage = false
	} else if result.LastEvaluatedKey == nil && sorted == false {
		lastPage = true
	}

	if sorted == true {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Sk > items[j].Sk // order descending
		})
	}

	var r []JSONKeys

	for _, el := range items {
		resp := JSONKeys{
			Damage:        el.Damage,
			Duration:      el.Duration,
			Deaths:        el.Deaths,
			Affixes:       el.Affixes,
			Keylevel:      el.Keylevel,
			DungeonName:   el.DungeonName,
			DungeonID:     el.DungeonID,
			CombatlogUUID: el.CombatlogUUID,
		}
		r = append(r, resp)
	}
	logrus.Debug(r)

	var firstSk string
	var lastSk string

	if sorted == false {
		firstSk = *result.Items[0]["sk"].S
		lastSk = *result.Items[len(items)-1]["sk"].S
		// -1 would be the last item, which we truncated
		// we need to 2nd last as last key, for the next page to not skip an item
	} else {
		// if we go back on the pagination the order is reversed
		// which causes last and first key to be switched, if we
		// were to go back again, thing are messed up
		firstSk = *result.Items[len(items)-1]["sk"].S
		lastSk = *result.Items[0]["sk"].S
	}
	logrus.Debug(firstSk)
	logrus.Debug(lastSk)

	resp := JSONKeysResponse{
		Data:      r,
		FirstSK:   firstSk,
		LastSK:    lastSk,
		FirstPage: firstPage,
		LastPage:  lastPage,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(b), err
}

/*
func PlayerDamageSimpleResponseToJson2(result *dynamodb2.QueryOutput, sorted, firstPage bool) (string, error) {
	var items []DynamoDBPlayerDamageSimple
	var lastPage bool //controls if the next btn is shown

	err := attributevalue2.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return "", err
	}
	if result.LastEvaluatedKey != nil {
		items = items[:len(items)-1] //we don't care about the last item
		//getting 1 more item than we need and truncaten it off allows us to circumvent a dynamodb2 bug where
		//if 5 items are left to return and you get exactly return 5 items the last evaluated key is not null,
		//which leads to the next request being completely empty
		lastPage = false
	} else if result.LastEvaluatedKey == nil && sorted == true {
		firstPage = true
		lastPage = false
	} else if result.LastEvaluatedKey == nil && sorted == false {
		lastPage = true
	}

	if sorted == true {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Sk > items[j].Sk // order descending
		})
	}

	var r []JSONKeys

	for _, el := range items {
		resp := JSONKeys{
			Damage:      el.Damage,
			Duration:    el.Duration,
			Deaths:      el.Deaths,
			Affixes:     el.Affixes,
			Keylevel:    el.Keylevel,
			DungeonName: el.DungeonName,
			DungeonID:   el.DungeonID,
		}
		r = append(r, resp)
	}
	logrus.Debug(r)

	var firstSk types2.AttributeValue
	var lastSk types2.AttributeValue

	if sorted == false {
		firstSk = result.Items[0]["sk"]
		lastSk = result.Items[len(items)-1]["sk"]
		//-1 would be the last item, which we truncated
		//we need to 2nd last as last key, for the next page to not skip an item
	} else {
		//if we go back on the pagination the order is reveresed
		//which causes last and first key to be switched, if we
		//were to go back again, thing are messed up
		firstSk = result.Items[len(items)-1]["sk"]
		lastSk = result.Items[0]["sk"]
	}
	logrus.Debug(firstSk)
	logrus.Debug(lastSk)

	resp := JSONPlayerDamageSimpleResponse2{
		Data:      r,
		FirstSK:   firstSk,
		LastSK:    lastSk,
		FirstPage: firstPage,
		LastPage:  lastPage,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(b), err
}
*/

// PlayerDamageDoneToJson returns the log specific damage result, including damage
// per spell breakdown, both damage per player and per spell are sorted before saving
// to the db.
// We don't need an extra JSON struct, like for keys, because there is no
// pagination etc.
func PlayerDamageDoneToJson(result *dynamodb.GetItemOutput) (string, error) {
	var item DynamoDBPlayerDamageDone

	err := dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func KeysToJson(result *dynamodb.QueryOutput) (string, error) {
	var items []DynamoDBKeys

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return "", err
	}

	var r []JSONKeys

	for _, el := range items {
		resp := JSONKeys{
			Damage:      el.Damage,
			Duration:    el.Duration,
			Deaths:      el.Deaths,
			Affixes:     el.Affixes,
			Keylevel:    el.Keylevel,
			DungeonName: el.DungeonName,
			DungeonID:   el.DungeonID,
		}
		r = append(r, resp)
	}
	logrus.Debug(r)

	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// InitLogging sets up the logging for every lambda and should be called before the handler
func InitLogging() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

// CanonicalLog writes a structured message to stdout if the log level is atleast INFO
func CanonicalLog(msg map[string]interface{}) {
	logrus.WithFields(msg).Info()
}

func SNSPublishMsg(ctx aws.Context, snsSvc *sns.SNS, input string, topicArn *string) error {
	if input == "" {
		return fmt.Errorf("input can't be empty")
	}
	logrus.Debug("sns input to publish: ", input)

	var err error

	if os.Getenv("LOCAL") == "true" {
		_, err = snsSvc.Publish(&sns.PublishInput{
			Message:  aws.String(input),
			TopicArn: topicArn,
		})
	} else {
		_, err = snsSvc.PublishWithContext(ctx, &sns.PublishInput{
			Message:  aws.String(input),
			TopicArn: topicArn,
		})
	}
	if err != nil {
		return fmt.Errorf("failed publishing a message to sns: %v", err)
	}

	logrus.Debug("message successfully sent to topic")
	return nil
}

/*
func SNSPublishMsg2(client *sns2.Client, input string, topicArn *string) error {
	if input == "" {
		return fmt.Errorf("combatlog_uuid can't be empty")
	}
	logrus.Debug("sns2 input to publish: ", input)

	_, err := client.Publish(context.TODO(), &sns2.PublishInput{
		Message:  aws2.String(input),
		TopicArn: topicArn,
	})
	if err != nil {
		return err
	}

	logrus.Debug("message successfully sent to topic")
	return nil
}
*/

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
	if len(result.Rows) == 0 {
		return result, fmt.Errorf("query returned empty result")
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

func UploadToTimestream(ctx aws.Context, writeSvc *timestreamwrite.TimestreamWrite, e []*timestreamwrite.Record) error {

	for i := 0; i < len(e); i += 100 {

		// get the upper bound of the record to write, in case it is the
		// last bit of records and i + 99 does not exist
		j := 0
		if i+99 > len(e) {
			j = len(e)
		} else {
			j = i + 99
		}

		// use common batching https://docs.aws.amazon.com/timestream/latest/developerguide/metering-and-pricing.writes.html#metering-and-pricing.writes.write-size-multiple-events
		writeRecordsInput := &timestreamwrite.WriteRecordsInput{
			// TODO: add to and read from env
			DatabaseName: aws.String("wowmate-analytics"),
			TableName:    aws.String("combatlogs"),
			Records:      e[i:j], // only upload a part of the records
		}

		var err error

		if os.Getenv("LOCAL") == "true" {
			_, err = writeSvc.WriteRecords(writeRecordsInput)
		} else {
			_, err = writeSvc.WriteRecordsWithContext(ctx, writeRecordsInput)
		}
		if err != nil {
			return fmt.Errorf("failed to write to timestream: %v", err)
		}
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

func PrettyStruct(input interface{}) (string, error) {
	prettyJSON, err := json.MarshalIndent(input, "", "    ")
	if err != nil {
		return "", err
	}
	return string(prettyJSON), nil
}

func DownloadS3(ctx aws.Context, downloader *s3manager.Downloader, bucketName string, objectKey string, fileContent io.WriterAt) error {
	var err error

	if os.Getenv("LOCAL") == "true" {
		_, err = downloader.Download(
			fileContent,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
			},
		)
	} else {
		_, err = downloader.DownloadWithContext(
			ctx,
			fileContent,
			&s3.GetObjectInput{
				Bucket: aws.String(bucketName),
				Key:    aws.String(objectKey),
			},
		)
	}
	if err != nil {
		return err
	}

	return nil
}

/*
func downloadS32(cfg aws2.Config, bucketName string, objectKey string, fileContent io.WriterAt) error {
	client := s32.NewFromConfig(cfg)
	downloader := manager2.NewDownloader(client)
	downloader.Download(context.TODO(), fileContent, &s32.GetObjectInput{
		Bucket: aws2.String(bucketName),
		Key:    aws2.String(objectKey),
	})
	return nil
}
*/

func RenameFileS3(ctx aws.Context, s3Svc *s3.S3, oldName, newName, bucketName string) error {
	source := fmt.Sprintf("%s/%s", bucketName, oldName)
	escapedSource := url.PathEscape(source)

	_, err := s3Svc.CopyObjectWithContext(ctx, &s3.CopyObjectInput{
		Key:        &newName,
		Bucket:     &bucketName,
		CopySource: &escapedSource,
	})
	if err != nil {
		return err
	}

	_, err = s3Svc.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Key:    &oldName,
		Bucket: &bucketName,
	})
	if err != nil {
		return err
	}

	return nil
}

// SizeOfS3Object checks the file size without actually downloading it in KB
func SizeOfS3Object(ctx aws.Context, s3Svc *s3.S3, bucketName string, objectKey string) (int, error) {
	var output *s3.HeadObjectOutput
	var err error
	// if it's local do without context
	if os.Getenv("LOCAL") == "true" {
		output, err = s3Svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	} else {
		output, err = s3Svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	}
	if err != nil {
		return 0, fmt.Errorf("unable to to send head request to item %q, %v", objectKey, err)
	}

	return int(*output.ContentLength / 1024), nil
}

// checks the file size without actually downloading it in MB
/*
func sizeOfS3Object2(cfg aws2.Config, bucketName string, objectKey string) (int, error) {
	client := s32.NewFromConfig(cfg)

	input := s32.HeadObjectInput{
		Bucket: aws2.String(bucketName),
		Key:    aws2.String(objectKey),
	}
	output, err := client.HeadObject(context.TODO(), &input)
	if err != nil {
		return 0, fmt.Errorf("Unable to to send head request to item %q, %v", objectKey, err)
	}

	return int(output.ContentLength / 1024 / 1024), nil
}
*/

// TimeMessageInQueue
// IMPROVE: only pass in pointer to SQSEvent
// we can see an overall trend of this in the SQS metrics, I created this to double check how fast messages are polled
// with fargate, but since I use lambda the lambda service takes care of this and messages are usually only 10ms old
// before they are pushed into lambda
func TimeMessageInQueue(e SQSEvent, i int) error {
	j, err := strconv.ParseInt(e.Records[i].Attributes.ApproximateFirstReceiveTimestamp, 10, 64)
	if err != nil {
		log.Printf("Failed to parse int: %v", err)
		return err
	}
	tm1 := time.Unix(0, j*int64(1000000))

	ii, err := strconv.ParseInt(e.Records[i].Attributes.SentTimestamp, 10, 64)
	if err != nil {
		return err
	}
	tm2 := time.Unix(0, ii*int64(1000000))

	log.Printf("seconds the message was unprocessed in the queue: %v", tm1.Sub(tm2).Seconds())

	return nil
}

// Atoi64 is just a small wrapper around ParseInt
func Atoi64(input string) (int64, error) {
	num, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

var levelTwoAffixes = map[int]string{
	9:  "Tyrannical",
	10: "Fortified",
}

var levelFourAffixes = map[int]string{
	123: "Spiteful",
	7:   "Bolstering",
	11:  "Bursting",
	8:   "Sanguine",
	6:   "Raging",
	122: "Inspiring",
}

var levelSevenAffixes = map[int]string{
	13:  "Explosive",
	4:   "Necrotic",
	3:   "Volcanic",
	124: "Storming",
	14:  "Quaking",
	12:  "Grievous",
}

var levelTenAffixes = map[int]string{
	121: "Prideful",
}

// AffixIDsToString takes an array of affix ids and converts them to a readable list of array names
// separated by commas
// TODO:
// 	- rethink if I actually need this, displaying a list of affixes as string takes a lot of space, should just return an array of ids and display icons
func AffixIDsToString(levelTwoID, levelFourID, levelSevenID, levelTenID int) string {
	affixes := levelTwoAffixes[levelTwoID]

	if levelFourID == 0 {
		return affixes
	}

	affixes += ", " + levelFourAffixes[levelFourID]

	if levelSevenID == 0 {
		return affixes
	}

	affixes += ", " + levelSevenAffixes[levelSevenID]

	if levelTenID == 0 {
		return affixes
	}

	affixes += ", " + levelTenAffixes[levelTenID]

	return affixes
}
