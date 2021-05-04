package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/mitchellh/hashstructure/v2"
	"golang.org/x/net/http2"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wowmate/jonny-rimek/wowmate/services/upload/convert/normalize"
)

/*
CREATE TABLE IF NOT EXISTS combatlogs (
  column_uuid UUID PRIMARY KEY,
  upload_uuid UUID,
  unsupported boolean,
  combatlog_uuid UUID,
  boss_fight_uuid UUID,
  mythicplus_uuid UUID,
  --timestamp timestamp, NOTE: deactivated till I figure out how to import it
  event_type VARCHAR,
  version int,
  advanced_log_enabled int,
  dungeon_name VARCHAR,
  dungeon_id int,
  key_unknown_1 int,
  key_level int,
  key_array VARCHAR,
  key_duration bigint,
  encounter_id int,
  encounter_name VARCHAR,
  encounter_unknown_1 int,
  encounter_unknown_2 int,
  killed int,
  caster_id VARCHAR,
  caster_name VARCHAR,
  caster_type VARCHAR,
  source_flag VARCHAR,
  target_id VARCHAR,
  target_name VARCHAR,
  target_type VARCHAR,
  dest_flag VARCHAR,
  spell_id int,
  spell_name VARCHAR,
  spell_type VARCHAR,
  extra_spell_id int,
  extra_spell_name VARCHAR,
  extra_school VARCHAR,
  aura_type VARCHAR,
  another_player_id VARCHAR,
  d0 VARCHAR,
  d1 bigint,
  d2 bigint,
  d3 bigint,
  d4 bigint,
  d5 bigint,
  d6 bigint,
  d7 bigint,
  d8 bigint,
  d9 VARCHAR,
  d10 VARCHAR,
  d11 VARCHAR,
  d12 VARCHAR,
  d13 VARCHAR,
  damage_unknown_14 VARCHAR,
  actual_amount bigint,
  base_amount bigint,
  overhealing bigint,
  overkill VARCHAR,
  school VARCHAR,
  resisted VARCHAR,
  blocked VARCHAR,
  absorbed bigint,
  critical VARCHAR,
  glancing VARCHAR,
  crushing VARCHAR,
  is_offhand VARCHAR
  --TODO: created_at
  --TODO: updated_at
);
*/
// https://mholt.github.io/json-to-go/ best tool EVER

type logData struct {
	CombatlogUUIDs     []string
	UploadUUID         string
	BucketName         string
	ObjectKey          string
	FileType           string
	FileSize           int
	Records            int
	TimestreamAPICalls int
	// TODO: wcu and rcu
}

var s3Svc *s3.S3
var snsSvc *sns.SNS
var downloader *s3manager.Downloader
var writeSvc *timestreamwrite.TimestreamWrite
var dynamodbSvc *dynamodb.DynamoDB

// I need an array of combatlog uuids for my integration test, because I need
// a recently inserted combatlog uuid, because the timestream query only checks data
// of the last 15minutes, so if I hardcode a combatlogUUID it would eventually result
// in an empty query from timestream
//goland:noinspection GoNilness
func handler(ctx aws.Context, e golib.SQSEvent) ([]string, error) {
	logData, err := handle(ctx, e)
	if err != nil {
		// create custom error types https://blog.golang.org/error-handling-and-go
		// TODO: reactivate to release, ruins my load tests
		// err2 := golib.RenameFileS3(
		// 	ctx,
		// 	s3Svc,
		// 	logData.ObjectKey,
		// 	fmt.Sprintf("error/%v", logData.UploadUUID),
		// 	logData.BucketName,
		// )
		// if err2 != nil {
		// 	golib.CanonicalLog(map[string]interface{}{
		// 		"combatlog_uuid":  logData.CombatlogUUIDs,
		// 		"file_size_in_kb": logData.FileSize,
		// 		"file_type":       logData.FileType,
		// 		"upload_uuid":     logData.UploadUUID,
		// 		"object_key":      logData.ObjectKey,
		// 		"bucket_name":     logData.BucketName,
		// 		"err":             err.Error(),
		// 		"err2":            err2.Error(),
		// 		"records":         logData.Records, // printing record in case of error to better debug
		// "timestream_api_calls": logData.TimestreamAPICalls,
		// 		"event":           e,
		// 	})
		// 	return err
		// }
		golib.CanonicalLog(map[string]interface{}{
			"combatlog_uuid":       logData.CombatlogUUIDs,
			"file_size_in_kb":      logData.FileSize,
			"file_type":            logData.FileType,
			"upload_uuid":          logData.UploadUUID,
			"object_key":           logData.ObjectKey,
			"bucket_name":          logData.BucketName,
			"err":                  err.Error(),
			"records":              logData.Records, // printing record in case of error to better debug
			"timestream_api_calls": logData.TimestreamAPICalls,
			"event":                e,
		})
		return logData.CombatlogUUIDs, err
	}

	golib.CanonicalLog(map[string]interface{}{
		"combatlog_uuid":       logData.CombatlogUUIDs,
		"file_size_in_kb":      logData.FileSize,
		"file_type":            logData.FileType,
		"upload_uuid":          logData.UploadUUID,
		"object_key":           logData.ObjectKey,
		"bucket_name":          logData.BucketName,
		"records":              logData.Records,
		"timestream_api_calls": logData.TimestreamAPICalls,
	})
	return logData.CombatlogUUIDs, nil
}

func handle(ctx aws.Context, e golib.SQSEvent) (logData, error) {
	var logData logData

	topicArn := os.Getenv("TOPIC_ARN") // +deploy trigger
	if topicArn == "" {
		return logData, fmt.Errorf("arn topic env var is empty")
	}

	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if ddbTableName == "" {
		return logData, fmt.Errorf("dynamodb table name env var is empty")
	}

	if len(e.Records) != 1 {
		// I also don't plan to change this as invocation costs are irrelevant part of the bill atm
		// and it makes the code way easier
		return logData, fmt.Errorf("the code only supports a batch size of 1. the current size of the batch is %v", len(e.Records))
	}

	// we only ever have 1 msg in the batch from sqs, that's why we dont need to loop through them, or above will return an err
	msg := e.Records[0]
	req := golib.S3Event{}
	err := json.Unmarshal([]byte(msg.Body), &req)
	if err != nil {
		return logData, fmt.Errorf("failed to unmarshal s3 event from message body to json: body: %s error: %v", msg.Body, err)
	}

	if len(req.Records) != 1 {
		// had some cases where the len of records was 0, should keep an eye on that
		return logData, fmt.Errorf("s3 event contains %v elements instead of 1", len(req.Records))
	}

	bucketName := req.Records[0].S3.Bucket.Name
	objectKey := req.Records[0].S3.Object.Key

	var maxSizeInKB int
	maxSizeInKB, logData.FileType, err = fileType(objectKey)
	if err != nil {
		return logData, err
	}

	uploadUUID, err := uploadUUID(objectKey)
	if err != nil {
		return logData, fmt.Errorf("failed to extract the uploadUUID: %v", err.Error())
	}

	logData.BucketName = bucketName
	logData.ObjectKey = objectKey
	logData.UploadUUID = uploadUUID

	// zipped files are usually around 8% of the original size
	logData.FileSize, err = golib.SizeOfS3Object(ctx, s3Svc, bucketName, objectKey)
	if err != nil {
		return logData, fmt.Errorf("failed to get size of s3 object: %v", err.Error())
	}

	if logData.FileSize > maxSizeInKB {
		return logData, fmt.Errorf("file is to big(%v kb) for it's file type(%v)", logData.FileSize, logData.FileType)
	}

	fileContent := &aws.WriteAtBuffer{}

	err = golib.DownloadS3(ctx, downloader, bucketName, objectKey, fileContent)
	if err != nil {
		return logData, fmt.Errorf("failed to download from s3: %v", err.Error())
	}

	data, err := readFileTypes(logData.FileType, fileContent)
	if err != nil {
		return logData, fmt.Errorf("failed to read file content for file type: %v error:%v", logData.FileType, err.Error())
	}

	s := bufio.NewScanner(bytes.NewReader(data))

	nestedRecord, dedup, err := normalize.Normalize(s, uploadUUID)
	if err != nil {
		return logData, fmt.Errorf("normalizing failed: %v", err)
	}
	// deduplication doesn't have a noticeable impact on performance or memory usage!
	for combatlogUUID, record := range dedup {
		hash, err := hashstructure.Hash(record, hashstructure.FormatV2, nil)
		if err != nil {
			return logData, fmt.Errorf("failed to hash: %v", err.Error())
		}
		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		log.Println(combatlogUUID)
		log.Println(hash)
		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

		input := &dynamodb.GetItemInput{
			TableName: &ddbTableName,
			Key: map[string]*dynamodb.AttributeValue{
				"pk": {
					S: aws.String(fmt.Sprintf("DEDUP#%v", hash)),
				},
				"sk": {
					S: aws.String(fmt.Sprintf("DEDUP#%v", hash)),
				},
			},
			ReturnConsumedCapacity: aws.String("TOTAL"),
		}
		_, err = golib.DynamoDBGetItem(ctx, dynamodbSvc, input)
		if err != nil {
			return logData, err
		}

	}
	/*
		TODO:
			- move GetItem code from per log dmg api call to golib
			- ddb GetItem for existing hash pk and sk DEDUP#HASHVALUE + combatloguuid
			- write to ddb if not exist
			- save duplicate combatlogUUIDs
			- skip timestream write and sns publish for duplicate combatlogUUIDs
	*/

	// return logData, nil

	// can process single key log files with 26MB size and 1792 MB memory lambda in ~3sec once timestream is warm
	maxGoroutines := 15
	var ch = make(chan *timestreamwrite.WriteRecordsInput, 300) // This number 200 can be anything as long as it's larger than maxGoroutines
	var wg sync.WaitGroup

	var writeErr error

	wg.Add(maxGoroutines) // this start maxGoroutines number of goroutines that wait for something to do
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				writeErr = golib.WriteToTimestream(ctx, writeSvc, a)
			}
		}()
	}

	for _, record := range nestedRecord { // group by different keys
		for _, writeRecordsInputs := range record { // grouped by key to use common attribute
			logData.TimestreamAPICalls += len(writeRecordsInputs)
			for _, e := range writeRecordsInputs { // array of TimestreamWriteInputs
				ch <- e                           // add i to the queue
				logData.Records += len(e.Records) // not sure if this is problematic or I should use channels
			}
		}
	}

	// close(errorChannel)
	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish

	if writeErr != nil {
		return logData, writeErr
	}

	for combatlogUUID := range nestedRecord { // group by different keys
		logData.CombatlogUUIDs = append(logData.CombatlogUUIDs, combatlogUUID)
		err = golib.SNSPublishMsg(ctx, snsSvc, combatlogUUID, &topicArn)
		if err != nil {
			return logData, err
		}
	}
	return logData, nil
}

// probably also reasonably testable, but file handling is always weird
func readFileTypes(fileType string, fileContent *aws.WriteAtBuffer) ([]byte, error) {
	var data []byte

	if fileType == "gz" {
		buf := bytes.NewBuffer(fileContent.Bytes())
		r, err := gzip.NewReader(buf)
		if err != nil {
			return nil, err
		}

		var resB bytes.Buffer
		_, err = resB.ReadFrom(r)
		if err != nil {
			return nil, err
		}

		data = resB.Bytes()
	} else if fileType == "zip" {
		zipReader, err := zip.NewReader(bytes.NewReader(fileContent.Bytes()), int64(len(fileContent.Bytes())))
		if err != nil {
			return nil, err
		}

		for _, zipFile := range zipReader.File {
			// log.Printf("zip loop i = %v", i)
			// log.Println("Reading file:", zipFile.Name)
			unzippedFileBytes, err := readZipFile(zipFile)
			if err != nil {
				return nil, err
			}

			data = append(data, unzippedFileBytes...)
		}
	} else {
		// filetype txt
		data = fileContent.Bytes()
	}
	return data, nil
}

func fileType(objectKey string) (int, string, error) {
	var fileType string
	var maxSizeInKB int
	if strings.HasSuffix(objectKey, ".txt") {
		maxSizeInKB = 450 * 1024 // 1GB
		fileType = "txt"
	} else if strings.HasSuffix(objectKey, ".txt.gz") {
		maxSizeInKB = 40 * 1024 // 100MB
		fileType = "gz"
	} else if strings.HasSuffix(objectKey, ".zip") {
		maxSizeInKB = 40 * 1024 // 100MB
		fileType = "zip"
	} else {
		return 0, "", fmt.Errorf("file suffix is not supported, s3 prefix filtering is broken")
	}
	return maxSizeInKB, fileType, nil
}

func readZipFile(zf *zip.File) (data []byte, err error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			err = fmt.Errorf("failed to close zip file: %v", err)
		}
	}()
	return ioutil.ReadAll(f)
}

func uploadUUID(s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("input can't be empty")
	}
	s = strings.TrimSuffix(s, ".txt")
	s = strings.TrimSuffix(s, ".txt.gz")
	s = strings.TrimSuffix(s, ".zip")

	split := strings.Split(s, "/")
	correctLength := 6
	if len(split) != correctLength {
		return "", fmt.Errorf("input has the wrong length, got %d want %d", len(split), correctLength)
	}

	lastElement := split[correctLength-1]
	return lastElement, nil
}

func main() {
	golib.InitLogging()
	tr := &http.Transport{
		ResponseHeaderTimeout: 20 * time.Second,
		// Using DefaultTransport values for other parameters: https://golang.org/pkg/net/http/#RoundTripper
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			// DualStack: true, // deprecated
			Timeout: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// So client makes HTTP/2 requests
	err := http2.ConfigureTransport(tr)
	if err != nil {
		log.Printf("failed configuring http2 transport: %v", err.Error())
		return
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1"), MaxRetries: aws.Int(10), HTTPClient: &http.Client{Transport: tr}})
	// sess, err := session.NewSession()
	if err != nil {
		log.Printf("failed creating session: %v", err.Error())
		return
	}

	xray.SetLogger(xraylog.NullLogger)
	/*
		I get a lot of the following messages. The reason according to the support is that to upload background processes
		are used, which don't have access to the context. To not spam my log and reduce cost I'm disabling the message with the
		silent setting.

		2021-04-14T05:50:30 2021-04-14T05:50:30Z [ERROR] Suppressing AWS X-Ray context missing panic: failed to begin subsegment named 'Timestream Write': segment cannot be found.
		2021-04-14T05:50:30 2021-04-14T05:50:30Z [ERROR] Suppressing AWS X-Ray context missing panic: failed to begin subsegment named 'attempt': segment cannot be found.
		2021-04-14T05:50:30 2021-04-14T05:50:30Z [ERROR] Suppressing AWS X-Ray context missing panic: failed to begin subsegment named 'unmarshal': segment cannot be found.
	*/

	s3Svc = s3.New(sess)
	if os.Getenv("LOCAL") == "false" {
		xray.AWS(s3Svc.Client)
	}

	downloader = s3manager.NewDownloaderWithClient(s3Svc)

	snsSvc = sns.New(sess)
	if os.Getenv("LOCAL") == "false" {
		xray.AWS(snsSvc.Client)
	}

	writeSvc = timestreamwrite.New(sess)
	if os.Getenv("LOCAL") == "false" {
		xray.AWS(writeSvc.Client)
	}

	dynamodbSvc = dynamodb.New(sess)
	// the aws docs recommend to set custom http settings see here:
	// https://docs.aws.amazon.com/timestream/latest/developerguide/code-samples.write-client.html
	// I'm choosing to ignore them and go with default, I'll observer if it leads to any problems

	lambda.Start(handler)
}
