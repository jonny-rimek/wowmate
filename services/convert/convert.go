package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	uuid "github.com/gofrs/uuid"
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

//https://mholt.github.io/json-to-go/ best tool EVER

type Request struct {
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

//SQSEvent is all the data that gets passed into the lambda from the q
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

func handler(e SQSEvent) error {

	csvBucket := os.Getenv("CSV_BUCKET_NAME")
	if csvBucket == "" {
		return fmt.Errorf("csv bucket env var is empty")
	}

	sess, _ := session.NewSession()

	if len(e.Records) == 0 {
		return fmt.Errorf("SQS Event doesn't contain any messages")
	}
	log.Printf("amount of messages %v", len(e.Records))

	for j, msg := range e.Records {
		err := timeMessageInQueue(e, j)
		if err != nil {
			return nil
		}

		log.Printf("index j: %v", j+1)

		//get sqs message body which contains the s3 event
		body := msg.Body

		req := Request{}
		err = json.Unmarshal([]byte(body), &req)
		if err != nil {
			return err
		}

		if len(req.Records) > 1 {
			log.Printf("Failed: the S3 event contains more than 1 element, not sure how that would happen")
			return err
		}

		bucketName := req.Records[0].S3.Bucket.Name
		objectKey := req.Records[0].S3.Object.Key

		objectSize, err := sizeOfS3Object(sess, bucketName, objectKey)
		if err != nil {
			return err
		}
		log.Printf("Object is %v MB", objectSize)

		if objectSize < 300 {
			fileContent := &aws.WriteAtBuffer{}

			err := downloadS3(sess, bucketName, objectKey, fileContent)
			if err != nil {
				return err
			}

			s := bufio.NewScanner(bytes.NewReader(fileContent.Bytes()))
			uploadUUID := uuid.Must(uuid.NewV4()).String()

			err = Normalize(s, uploadUUID, sess, csvBucket)
			if err != nil {
				return err
			}
		} else if objectSize >= 300 && objectSize < 20000 {
			log.Println("file between 300MB and 20GB")

			file, err := os.Create("/mnt/efs/" + objectKey)
			if err != nil {
				return err
			}

			err = downloadS3(sess, bucketName, objectKey, file)
			if err != nil {
				return err
			}

		} else {
			return fmt.Errorf("can't process files larger than 20GB")
		}
	}
	return nil
}

func downloadS3(sess *session.Session, bucketName string, objectKey string, fileContent io.WriterAt) error {
	downloader := s3manager.NewDownloader(sess)

	_, err := downloader.Download(
		fileContent,
		&s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func convertToCSV(events *[]Event) (io.Reader, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	ss, err := EventsAsStringSlices(events)
	if err != nil {
		return nil, err
	}
	log.Println("converted to struct to string slice")

	//flushes the string slice as csv to buffer
	if err := w.WriteAll(ss); err != nil {
		return nil, err
	}
	log.Println("converted to csv")

	return io.Reader(&buf), nil
}

func sizeOfS3Object(sess *session.Session, bucketName string, objectKey string) (int, error) {
	svc := s3.New(sess)

	output, err := svc.HeadObject(
		&s3.HeadObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(objectKey),
		})
	if err != nil {
		return 0, fmt.Errorf("Unable to to send head request to item %q, %v", objectKey, err)
	}

	return int(*output.ContentLength / 1024 / 1024), nil
}

func uploadS3(r io.Reader, sess *session.Session, mythicplugUUID string, csvBucket string) error {
	if mythicplugUUID == "" {
		//sometimes there are more CHALLANGE_MODE_END events than there are start events
		//it shouldn't come to this, because we aren't adding anything unless we have a started event
		return nil
	}
	uploader := s3manager.NewUploader(sess)

	//TODO: check that mythicplusUUID is not ""

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(csvBucket),
		Key:    aws.String(fmt.Sprintf("%v.csv", mythicplugUUID)),
		Body:   r,
	})
	if err != nil {
		log.Println("Failed to upload to S3")
		return err
	}
	log.Println("Upload finished! location: " + result.Location)

	return nil
}

func timeMessageInQueue(e SQSEvent, i int) error {
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

func main() {
	lambda.Start(handler)
}
