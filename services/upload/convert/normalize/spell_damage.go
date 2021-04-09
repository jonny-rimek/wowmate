package normalize

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

// v17
// 1/24 16:47:51.662  SPELL_DAMAGE,Player-581-04A01EDA,"Ayléén-Blackrock",0x511,0x0,Creature-0-4234-2291-11942-168992-00000DA4A7,
// "Risen Cultist",0x10a48,0x0,50842,"Blood Boil", 0x20,Creature-0-4234-2291-11942-168992-00000DA4A7,0000000000000000,177038,177822,
// 0,0,1071,0,0,2434,2434,0,2943.32,-2219.79,1680,3.6342,60,784,783,-1,32,0,0,0,nil,nil,nil
// TODO: try out []*string for params or []byte
// 	- pass in timestreamwrite array as pointer
// 	- uuids as pointer
func spellDamage(params []string, uploadUUID string, combatlogUUID string) (*timestreamwrite.Record, error) {
	if len(params) != 39 {
		return nil, fmt.Errorf("combatlog version should have 39 columns, it has %v: %v", len(params), params)
	}

	actualAmount, err := Atoi64(params[29]) // 1287
	if err != nil {
		log.Printf("failed to convert damage event, field actual amount. got: %v", params[27])
		return nil, err
	}

	spellID, err := strconv.Atoi(params[9]) // 283810
	if err != nil {
		log.Printf("failed to convert damage event, field spell id. got: %v", params[9])
		return nil, err
	}
	// can specify seconds as input for timestream instead of ms!
	currentTimeInSeconds := time.Now().Unix()

	rand.Seed(time.Now().UnixNano())

	e := &timestreamwrite.Record{
		Dimensions: []*timestreamwrite.Dimension{
			{
				Name:  aws.String("caster_id"),
				Value: aws.String(params[1]),
			},
			{
				Name:  aws.String("caster_name"),
				Value: aws.String(trimQuotes(params[2])),
			},
			{
				Name:  aws.String("caster_type"),
				Value: aws.String(trimQuotes(params[3])),
			},
			{
				Name:  aws.String("spell_id"),
				Value: aws.String(strconv.Itoa(spellID)),
			},
			{
				Name:  aws.String("spell_name"),
				Value: aws.String(trimQuotes(params[10])),
			},
			{
				Name:  aws.String("spell_type"),
				Value: aws.String(params[11]),
			},
			{
				Name:  aws.String("upload_uuid"),
				Value: aws.String(uploadUUID),
			},
			{
				Name:  aws.String("combatlog_uuid"),
				Value: aws.String(combatlogUUID),
			},
			{
				Name:  aws.String("rnd"),
				Value: aws.String(strconv.Itoa(rand.Int())),
				// only randomizing between 1 and 999 didnt work
				// maybe I can remove this when I add the real time as a value
			},
		},
		MeasureName:      aws.String("damage"),
		MeasureValue:     aws.String(strconv.FormatInt(actualAmount, 10)),
		MeasureValueType: aws.String("BIGINT"),
		Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
		TimeUnit:         aws.String("SECONDS"),
	}

	return e, nil
}
