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
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"

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
  key_unkown_1 int,
  key_level int,
  key_array VARCHAR,
  key_duration bigint,
  encounter_id int,
  encounter_name VARCHAR,
  encounter_unkown_1 int,
  encounter_unkown_2 int,
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
	CombatlogUUID []string
	UploadUUID    string
	BucketName    string
	ObjectKey     string
	FileType      string
	FileSize      int
	Records       int
}

var s3Svc *s3.S3
var snsSvc *sns.SNS
var downloader *s3manager.Downloader
var writeSvc *timestreamwrite.TimestreamWrite

//goland:noinspection GoNilness
func handler(ctx aws.Context, e golib.SQSEvent) error {
	logData, err := handle(ctx, e)
	if err != nil {
		err2 := golib.RenameFileS3(
			ctx,
			s3Svc,
			logData.ObjectKey,
			fmt.Sprintf("error/%v", logData.UploadUUID),
			logData.BucketName,
		)
		if err2 != nil {
			golib.CanonicalLog(map[string]interface{}{
				"combatlog_uuid":  logData.CombatlogUUID,
				"file_size_in_kb": logData.FileSize,
				"file_type":       logData.FileType,
				"upload_uuid":     logData.UploadUUID,
				"object_key":      logData.ObjectKey,
				"bucket_name":     logData.BucketName,
				"err":             err.Error(),
				"err2":            err2.Error(),
				"records":         logData.Records, // printing record in case of error to better debug
				"event":           e,
			})
			return err
		}
		golib.CanonicalLog(map[string]interface{}{
			"combatlog_uuid":  logData.CombatlogUUID,
			"file_size_in_kb": logData.FileSize,
			"file_type":       logData.FileType,
			"upload_uuid":     logData.UploadUUID,
			"object_key":      logData.ObjectKey,
			"bucket_name":     logData.BucketName,
			"err":             err.Error(),
			"records":         logData.Records, // printing record in case of error to better debug
			"event":           e,
		})
		return err
	}

	golib.CanonicalLog(map[string]interface{}{
		"combatlog_uuid":  logData.CombatlogUUID,
		"file_size_in_kb": logData.FileSize,
		"file_type":       logData.FileType,
		"upload_uuid":     logData.UploadUUID,
		"object_key":      logData.ObjectKey,
		"bucket_name":     logData.BucketName,
		"records":         logData.Records,
	})
	return err
}

func handle(ctx aws.Context, e golib.SQSEvent) (logData, error) {
	var logData logData

	topicArn := os.Getenv("TOPIC_ARN")
	if topicArn == "" {
		return logData, fmt.Errorf("arn topic env var is empty")
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
		return logData, err
	}

	if len(req.Records) != 1 {
		// had some cases where the len of records was 0, should keep an eye on that
		return logData, fmt.Errorf("s3 event contains %v elements instead of 1", len(req.Records))
	}

	bucketName := req.Records[0].S3.Bucket.Name
	objectKey := req.Records[0].S3.Object.Key
	uploadUUID := uploadUUID(objectKey)

	logData.BucketName = bucketName
	logData.ObjectKey = objectKey
	logData.UploadUUID = uploadUUID

	var maxSizeInKB int

	maxSizeInKB, logData.FileType, err = fileType(objectKey)
	if err != nil {
		return logData, err
	}

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

	nestedRecord, err := normalize.Normalize(s, uploadUUID)
	if err != nil {
		return logData, err
	}

	// if os.Getenv("LOCAL") == "true" {
	// uploading to timestream takes too long locally, option to not run it
	//	return logData, nil
	// }
	// logData.Records = len(combatEvents)

	for _, record := range nestedRecord {
		for combatlogUUID, writeRecordsInputs := range record {
			for _, e := range writeRecordsInputs {
				err = golib.WriteToTimestream(ctx, writeSvc, e)
				if err != nil {
					return logData, err
				}
			}

			logData.CombatlogUUID = append(logData.CombatlogUUID, combatlogUUID)

			err = golib.SNSPublishMsg(ctx, snsSvc, combatlogUUID, &topicArn)
			if err != nil {
				return logData, err
			}
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

		for i, zipFile := range zipReader.File {
			log.Printf("zip loop i = %v", i)
			fmt.Println("Reading file:", zipFile.Name)
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

// TODO: ez test
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

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func uploadUUID(s string) string {
	// TODO: check if string is empty and return error
	//	update test too
	s = strings.TrimSuffix(s, ".txt")
	s = strings.TrimSuffix(s, ".txt.gz")
	s = strings.TrimSuffix(s, ".zip")

	return strings.Split(s, "/")[5]
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		log.Printf("Error creating session: %v", err.Error())
	}

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
	// the aws docs recommend to set custom http settings see here:
	// https://docs.aws.amazon.com/timestream/latest/developerguide/code-samples.write-client.html
	// I'm choosing to ignore them and go with default, I'll observer if it leads to any problems

	lambda.Start(handler)
}
