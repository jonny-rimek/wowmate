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
//
// passing in the uuids as a pointer to the string reduced the mb usage by a couple MB
func spellDamage(params []string, uploadUUID *string, combatlogUUID *string, rec map[string]map[string][]*timestreamwrite.WriteRecordsInput) error {
	if len(params) != 39 {
		return fmt.Errorf("combatlog version should have 39 columns, it has %v: %v", len(params), params)
	}

	actualAmount, err := Atoi64(params[29]) // 1287
	if err != nil {
		log.Printf("failed to convert damage event, field actual amount. got: %v", params[27])
		return err
	}

	spellID, err := strconv.Atoi(params[9]) // 283810
	if err != nil {
		log.Printf("failed to convert damage event, field spell id. got: %v", params[9])
		return err
	}

	casterID := params[1]
	casterName := trimQuotes(params[2])
	casterType := trimQuotes(params[3])
	spellName := trimQuotes(params[10])
	spellType := params[11]
	currentTimeInSeconds := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())

	key := fmt.Sprintf("%s_%d", casterID, spellID)

	val, exists := rec[*combatlogUUID][key]
	if exists == true {
		for _, writeRecordsInput := range val {
			if len(writeRecordsInput.Records) < 100 {
				writeRecordsInput.Records = append(writeRecordsInput.Records, &timestreamwrite.Record{
					Dimensions: []*timestreamwrite.Dimension{
						{
							Name:  aws.String("rnd"),
							Value: aws.String(strconv.Itoa(rand.Int())),
							// replace with time from log
						},
					},
					MeasureValue: aws.String(strconv.FormatInt(actualAmount, 10)),
				})
			} else {
				/*
					writeInput := &timestreamwrite.WriteRecordsInput{
						DatabaseName: aws.String("wowmate-analytics"),
						TableName:    aws.String("combatlogs"),
						CommonAttributes: &timestreamwrite.Record{
							Dimensions: []*timestreamwrite.Dimension{
								{
									Name:  aws.String("caster_id"),
									Value: aws.String(casterID),
								},
								{
									Name:  aws.String("caster_name"),
									Value: aws.String(casterName),
								},
								{
									Name:  aws.String("caster_type"),
									Value: aws.String(casterType),
								},
								{
									Name:  aws.String("spell_id"),
									Value: aws.String(strconv.Itoa(spellID)),
								},
								{
									Name:  aws.String("spell_name"),
									Value: aws.String(spellName),
								},
								{
									Name:  aws.String("spell_type"),
									Value: aws.String(spellType),
								},
								{
									Name:  aws.String("upload_uuid"),
									Value: aws.String(*uploadUUID),
								},
								{
									Name:  aws.String("combatlog_uuid"),
									Value: aws.String(*combatlogUUID),
								},
							},
							MeasureName:      aws.String("damage"),
							MeasureValueType: aws.String("BIGINT"),
							TimeUnit:         aws.String("SECONDS"), // can specify seconds as rec for timestream instead of ms!
							Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
							// I don't care about this time, it just the time we create this entry, not the time of the combatlog event
							// I also don't care about the exact time this is written, so I always use the time the first record is created
							// and reuse it for the subsequent ones
						},
						Records: []*timestreamwrite.Record{
							{
								Dimensions: []*timestreamwrite.Dimension{
									{
										Name:  aws.String("rnd"),
										Value: aws.String(strconv.Itoa(rand.Int())),
										// replace with time from log
									},
								},
								MeasureValue: aws.String(strconv.FormatInt(actualAmount, 10)),
							},
						},
					}

					writeRecordsInputs := rec[*combatlogUUID][key]
					writeRecordsInputs = append(writeRecordsInputs, writeInput)

					rec[*combatlogUUID][key] = writeRecordsInputs*/
			}
		}
		return nil
	}

	writeRecordsInputs := []*timestreamwrite.WriteRecordsInput{
		{

			DatabaseName: aws.String("wowmate-analytics"),
			TableName:    aws.String("combatlogs"),
			CommonAttributes: &timestreamwrite.Record{
				Dimensions: []*timestreamwrite.Dimension{
					{
						Name:  aws.String("caster_id"),
						Value: aws.String(casterID),
					},
					{
						Name:  aws.String("caster_name"),
						Value: aws.String(casterName),
					},
					{
						Name:  aws.String("caster_type"),
						Value: aws.String(casterType),
					},
					{
						Name:  aws.String("spell_id"),
						Value: aws.String(strconv.Itoa(spellID)),
					},
					{
						Name:  aws.String("spell_name"),
						Value: aws.String(spellName),
					},
					{
						Name:  aws.String("spell_type"),
						Value: aws.String(spellType),
					},
					{
						Name:  aws.String("upload_uuid"),
						Value: aws.String(*uploadUUID),
					},
					{
						Name:  aws.String("combatlog_uuid"),
						Value: aws.String(*combatlogUUID),
					},
				},
				MeasureName:      aws.String("damage"),
				MeasureValueType: aws.String("BIGINT"),
				TimeUnit:         aws.String("SECONDS"), // can specify seconds as rec for timestream instead of ms!
				Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
				// I don't care about this time, it just the time we create this entry, not the time of the combatlog event
				// I also don't care about the exact time this is written, so I always use the time the first record is created
				// and reuse it for the subsequent ones
			},
			Records: []*timestreamwrite.Record{
				{
					Dimensions: []*timestreamwrite.Dimension{
						{
							Name:  aws.String("rnd"),
							Value: aws.String(strconv.Itoa(rand.Int())),
							// replace with time from log
						},
					},
					MeasureValue: aws.String(strconv.FormatInt(actualAmount, 10)),
				},
			},
		},
	}

	rec[*combatlogUUID][key] = writeRecordsInputs

	return nil
}
