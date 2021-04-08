package normalize

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"log"
	"strconv"
	"time"
)

//1/24 17:23:07.443  CHALLENGE_MODE_END,2291,1,10,2136094
//1/24 16:47:37.880  CHALLENGE_MODE_END,2291,0,0,0
//if it's not a finish [2] is 0 and the other two aswell
func challengeModeEnd(params []string, uploadUUID string, combatlogUUID string) ([]*timestreamwrite.Record, error) {
	if len(params) != 5 {
		return nil, fmt.Errorf("combatlog version should have 5 columns, it has %v: %v", len(params), params)
	}

	finished, err := Atoi64(params[2]) //2136094
	if err != nil {
		log.Printf("failed to convert challange mode end event, field finished. got: %v", params[2])
		return nil, err
	}

	//in milli seconds
	duration, err := Atoi64(params[4]) //2136094
	if err != nil {
		log.Printf("failed to convert challange mode end event, field duration. got: %v", params[4])
		return nil, err
	}

	currentTimeInSeconds := time.Now().Unix()

	var e = []*timestreamwrite.Record{
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
			MeasureName:      aws.String("duration"),
			MeasureValue:     aws.String(strconv.FormatInt(duration, 10)),
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
			MeasureName:      aws.String("finished"),
			MeasureValue:     aws.String(strconv.FormatInt(finished, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		},
	}

	return e, nil
}
