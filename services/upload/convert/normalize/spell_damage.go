package normalize

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"log"
	"math/rand"
	"strconv"
	"time"
)

//11/3 09:00:29.792  SPELL_DAMAGE,Player-1302-09C8C064,"Hyrriuk-Archimonde",0x512,0x0,Vehicle-0-3892-1763-30316-122963-00005D638F,"Rezan",0x10a48,0x0,283810,"Reckless Flurry",0x1,Vehicle-0-3892-1763-30316-122963-00005D638F,0000000000000000,3600186,3811638,0,0,2700,1,0,0,0,-790.59,2265.96,935,0.8059,122,1287,1599,-1,1,0,0,0,nil,nil,nil

// v16
// 10/3 05:51:15.415  SPELL_DAMAGE,Player-4184-00130F03,"Unstaebl-Torghast",0x512,0x0,Creature-0-2085-2287-15092-165515-0005F81144,"Depraved Darkblade",0xa48,0x0,127802,"Touch of the Grave",0x20,Creature-0-2085-2287-15092-165515-0005F81144,0000000000000000,92482,96120,0,0,1071,0,3,100,100,0,-2206.68,5071.68,1663,2.1133,60,456,456,-1,32,0,0,0,nil,nil,nil
func spellDamage(params []string, uploadUUID string, combatlogUUID string) (*timestreamwrite.Record, error) {
	if len(params) != 39 {
		return nil, fmt.Errorf("combatlog version should have 39 columns, it has %v: %v", len(params), params)
	}

	actualAmount, err := Atoi64(params[29]) //1287
	if err != nil {
		log.Printf("failed to convert damage event, field actual amount. got: %v", params[27])
		return nil, err
	}

	spellID, err := strconv.Atoi(params[9]) //283810
	if err != nil {
		log.Printf("failed to convert damage event, field spell id. got: %v", params[9])
		return nil, err
	}
	// currentTimeInMilliseconds := time.Now().UnixNano() / 1000000
	//TODO: timestream expects ms not seconds, double check this
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
				Value: aws.String(strconv.Itoa(rand.Intn(999-1) + 1)),
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
