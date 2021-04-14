package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	if err != nil {
		log.Printf("Error creating session: %v", err.Error())
	}
	svc := dynamodb.New(sess)

	var startKey map[string]*dynamodb.AttributeValue
	input := dynamodb.ScanInput{
		ExclusiveStartKey: startKey,
		TableName:         aws.String("wm-dev-DynamoDBtableF8E87752-HSV525WR7KN3"),
	}

	for {
		input.ExclusiveStartKey = startKey

		scan, err := svc.Scan(&input)
		if err != nil {
			log.Printf("failed to scan the table %v", err)
			return
		}
		// scan.Items
		for _, item := range scan.Items {
			delete(item, "affixes")
			delete(item, "combatlog_uuid")
			delete(item, "deaths")
			delete(item, "dungeon_id")
			delete(item, "dungeon_name")
			delete(item, "duration")
			delete(item, "finished")
			delete(item, "keylevel")
			delete(item, "player_damage")
			delete(item, "gsi1pk")
			delete(item, "gsi1sk")
			log.Println(item)
			output, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
				Key:       item,
				TableName: aws.String("wm-dev-DynamoDBtableF8E87752-HSV525WR7KN3"),
			})
			if err != nil {
				log.Printf("deleting failed %v\n", err)
				return
			}
			log.Println(output)
		}

		startKey = scan.LastEvaluatedKey
		log.Println(startKey)

		if len(startKey) == 0 {
			continue
		}
	}

	log.Println("hi")
}
