package normalize

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"log"
	"strconv"
	"time"
)

//11/3 09:00:00.760  CHALLENGE_MODE_START,"Atal'Dazar",1763,244,10,[10,11,14,16]

//v16
//10/3 05:51:00.975  CHALLENGE_MODE_START,"Halls of Atonement",2287,378,2,[10]
//TODO: add new example
//NOTE: the array is definitely the affixes
func challengeModeStart(params []string, uploadUUID string, combatlogUUID string) (*timestreamwrite.Record, error) {
	if len(params) != 6 {
		return nil, fmt.Errorf("combatlog version should have 6 columns, it has %v: %v", len(params), params)
	}

	dungeonID, err := Atoi64(params[1]) //283810
	if err != nil {
		log.Printf("failed to challange mode start event, field dungeon id. got: %v", params[1])
		return nil, err
	}

	currentTimeInSeconds := time.Now().Unix()

	var e = &timestreamwrite.Record{
		Dimensions: []*timestreamwrite.Dimension{
			{
				Name:  aws.String("dungeon_name"),
				Value: aws.String(params[0]),
			},
		},
		MeasureName:      aws.String("dungeon_id"),
		MeasureValue:     aws.String(strconv.FormatInt(dungeonID, 10)),
		MeasureValueType: aws.String("BIGINT"),
		Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
		TimeUnit:         aws.String("SECONDS"),
	}
	//TODO: add affixes and key level as seperate timestream record and return slice of records
	return e, nil
}
