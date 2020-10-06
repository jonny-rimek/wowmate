package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	_ "github.com/lib/pq"
)

var ConnStr string

var Sess *session.Session

//S3Event is the data that come from s3 and contains all the information about the event
type S3Event struct {
	Records []struct {
		EventVersion string    `json:"eventVersion"`
		EventSource  string    `json:"eventSource"`
		AwsRegion    string    `json:"awsRegion"`
		EventTime    time.Time `json:"eventTime"`
		EventName    string    `json:"eventName"`
		UserIdentity struct {
			PrincipalID string `json:"principalId"`
		} `json:"userIdentity"`
		RequestParameters struct {
			SourceIPAddress string `json:"sourceIPAddress"`
		} `json:"requestParameters"`
		ResponseElements struct {
			XAmzRequestID string `json:"x-amz-request-id"`
			XAmzID2       string `json:"x-amz-id-2"`
		} `json:"responseElements"`
		S3 struct {
			S3SchemaVersion string `json:"s3SchemaVersion"`
			ConfigurationID string `json:"configurationId"`
			Bucket          struct {
				Name          string `json:"name"`
				OwnerIdentity struct {
					PrincipalID string `json:"principalId"`
				} `json:"ownerIdentity"`
				Arn string `json:"arn"`
			} `json:"bucket"`
			Object struct {
				Key       string `json:"key"`
				Size      int    `json:"size"`
				ETag      string `json:"eTag"`
				Sequencer string `json:"sequencer"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

//TODO: use events.SQSEvent see summary lambda
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
		MessageAttributes struct {
		} `json:"messageAttributes"`
		Md5OfBody      string `json:"md5OfBody"`
		EventSource    string `json:"eventSource"`
		EventSourceARN string `json:"eventSourceARN"`
		AwsRegion      string `json:"awsRegion"`
	} `json:"Records"`
}

func handler(e SQSEvent) error {
	topicArn := os.Getenv("SUMMARY_TOPIC_NAME")
	if topicArn == "" {
		return fmt.Errorf("summary topic name env var is empty")
	}
	if ConnStr == "" {
		log.Println("ConnStr was empty")
		secretArn := os.Getenv("SECRET_ARN")
		dbEndpoint := os.Getenv("DB_ENDPOINT")
		ConnStr, _ = golib.DBCreds(secretArn, dbEndpoint, Sess)
	}
	db, err := sql.Open("postgres", ConnStr)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Println("opened connection")
	log.Printf("number of messages: %v", len(e.Records))

	for _, record := range e.Records {
		s3 := S3Event{}
		err = json.Unmarshal([]byte(record.Body), &s3)
		if err != nil {
			log.Println("Unmarshal sqs json failed")
			return err
		}

		if len(s3.Records) == 0 {
			return fmt.Errorf("failed: s3 event empty")
		}

		if len(s3.Records) > 1 {
			return fmt.Errorf("failed: the S3 event contains more than 1 element, not sure how that would happen")
		}
		q := fmt.Sprintf(`
				SELECT aws_s3.table_import_from_s3(
					'combatlogs',
					'column_uuid,upload_uuid,unsupported,combatlog_uuid,boss_fight_uuid,mythicplus_uuid,event_type,version,advanced_log_enabled,dungeon_name,dungeon_id,key_unkown_1,key_level,key_array,key_duration,encounter_id,encounter_name,encounter_unkown_1,encounter_unkown_2,killed,caster_id,caster_name,caster_type,source_flag,target_id,target_name,target_type,dest_flag,spell_id,spell_name,spell_type,extra_spell_id,extra_spell_name,extra_school,aura_type,another_player_id,d0,d1,d2,d3,d4,d5,d6,d7,d8,d9,d10,d11,d12,d13,damage_unknown_14,actual_amount,base_amount,overhealing,overkill,school,resisted,blocked,absorbed,critical,glancing,crushing,is_offhand',
					'(format csv, DELIMITER '','', HEADER true)',
					'(%v,%v,us-east-1)');
			`, s3.Records[0].S3.Bucket.Name, s3.Records[0].S3.Object.Key)

		rows, err := db.Query(q)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				//NOTE: this can happen because sometimes the import takes 10x the normal time and
				//		the import finished, after the lambda timed out
				//		which causes the lambda to be triggered again from SQS, instead of failing
				//		till it lands in the DLQ we will just delete it from the queue,
				//		because it got actually imported to the DB
				log.Printf("duplicate key for %v, error msg: %v ", s3.Records[0].S3.Object.Key, err.Error())
				return nil
			}
			return err
		}

		defer rows.Close()

		for rows.Next() {
			var s string //trigger

			err = rows.Scan(&s)
			if err != nil {
				log.Println(err.Error())
				return err
			}
			log.Printf("import query successfull: %v", s)

			svc := sns.New(Sess)
			result, err := svc.Publish(&sns.PublishInput{
				Message:  aws.String(s3.Records[0].S3.Object.Key),
				TopicArn: aws.String(topicArn),
			})
			if err != nil {
				return err
			}
			log.Printf("message posted to topic")
			fmt.Println(*result.MessageId)
		}
	}
	log.Println("summary successfully invoked")
	return nil
}

func main() {
	secretArn := os.Getenv("SECRET_ARN")
	if secretArn == "" {
		log.Println("secret arn env var is empty")
		return
	}

	dbEndpoint := os.Getenv("DB_ENDPOINT")
	if dbEndpoint == "" {
		log.Println("db endpoint env var is empty")
		return
	}

	//if I use := Sess is not set as a global var but as a local
	var err error
	Sess, err = session.NewSession()
	if err != nil {
		log.Println("failed to create new session")
		return
	}

	ConnStr, err = golib.DBCreds(secretArn, dbEndpoint, Sess)
	if err != nil {
		return
	}

	lambda.Start(handler)
}
