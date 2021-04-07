package normalize

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"log"
	"strconv"
	"strings"
	"time"
)

//1/24 16:47:38.068  CHALLENGE_MODE_START,"De Other Side",2291,377,10,[10,123,3,121]
//NOTE: the array is definitely the affixes
func challengeModeStart(params []string, uploadUUID string, combatlogUUID string) ([]*timestreamwrite.Record, error) {
	/* the array inside the square bracket also gets split by comma
	if len(params) != 6 {
		return nil, fmt.Errorf("combatlog version should have 6 columns, it has %v: %v", len(params), params)
	}
	*/

	dungeonID, err := Atoi64(params[2]) //2291
	if err != nil {
		log.Printf("failed to convert challange mode start event, field dungeon id. got: %v", params[2])
		return nil, err
	}
	keyLevel, err := Atoi64(params[4]) //10
	if err != nil {
		log.Printf("failed to convert challange mode start event, field key level. got: %v", params[4])
		return nil, err
	}

	currentTimeInSeconds := time.Now().Unix()

	var e = []*timestreamwrite.Record{
		{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("dungeon_name"),
					Value: aws.String(strings.Trim(params[1], "\"")),
				},
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
				},
				{
					Name:  aws.String("combatlog_uuid"),
					Value: aws.String(combatlogUUID),
				},
			},
			MeasureName:      aws.String("dungeon_id"),
			MeasureValue:     aws.String(strconv.FormatInt(dungeonID, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		},
		{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
				},
				{
					Name:  aws.String("combatlog_uuid"),
					Value: aws.String(combatlogUUID),
				},
			},
			//measure name is always damage, read docs
			MeasureName:      aws.String("key_level"),
			MeasureValue:     aws.String(strconv.FormatInt(keyLevel, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		},
	}

	//TODO: add affixes as seperate timestream record and return slice of records
	return e, nil
}
