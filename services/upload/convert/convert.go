package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
//IMPROVE: use events.SQSEvent
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

		log.Printf("%v. message(s) in SQS event", j+1)

		//get s3 event from sqs event body
		req := Request{}
		err = json.Unmarshal([]byte(msg.Body), &req)
		if err != nil {
			return err
		}

		if len(req.Records) > 1 {
			log.Printf("Failed: the S3 event contains more than 1 element, not sure how that would happen")
			return err
		}

		bucketName := req.Records[0].S3.Bucket.Name
		objectKey := req.Records[0].S3.Object.Key

		var maxSize int
		var gz bool
		var z bool

		if strings.HasSuffix(objectKey, ".txt") {
			maxSize = 1000
		} else if strings.HasSuffix(objectKey, ".txt.gz") {
			maxSize = 100
			gz = true
		} else if strings.HasSuffix(objectKey, ".zip") {
			maxSize = 100
			z = true
		} else {
			return fmt.Errorf("file suffix is not supported")
		}

		objectSize, err := sizeOfS3Object(sess, bucketName, objectKey)
		if err != nil {
			return err
		}
		log.Printf("Object is %v MB", objectSize)

		if objectSize > maxSize {
			return fmt.Errorf("wow that's huge. um phrasing?")
		}

		fileContent := &aws.WriteAtBuffer{}

		err = downloadS3(sess, bucketName, objectKey, fileContent)
		if err != nil {
			return err
		}

		var data []byte

		if gz /* == true*/ {
			buf := bytes.NewBuffer(fileContent.Bytes())
			r, err := gzip.NewReader(buf)
			if err != nil {
				return err
			}

			var resB bytes.Buffer
			_, err = resB.ReadFrom(r)
			if err != nil {
				return err
			}

			data = resB.Bytes()
			log.Println("successfully ungziped")
		} else if z /* == true*/ {
			zipReader, err := zip.NewReader(bytes.NewReader(fileContent.Bytes()), int64(len(fileContent.Bytes())))
			if err != nil {
				return err
			}

			for i, zipFile := range zipReader.File {
				log.Printf("zip loop i = %v", i)
				fmt.Println("Reading file:", zipFile.Name)
				unzippedFileBytes, err := readZipFile(zipFile)
				if err != nil {
					return err
				}

				data = append(data, unzippedFileBytes...)
			}
			log.Println("successfully unziped")
		} else {
			data = fileContent.Bytes()
		}

		s := bufio.NewScanner(bytes.NewReader(data))
		uploadUUID := uploadUUID(objectKey)

		err = normalize.Normalize(s, uploadUUID, sess, csvBucket)
		if err != nil {
			return err
		}
	}
	return nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
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

//checks the file size without actually downloading it in MB
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

func uploadUUID(s string) string {
	s = strings.TrimSuffix(s, ".txt")
	s = strings.TrimSuffix(s, ".txt.gz")
	s = strings.TrimSuffix(s, ".zip")

	return strings.Split(s, "/")[3]
}

func main() {
	lambda.Start(handler)
}
