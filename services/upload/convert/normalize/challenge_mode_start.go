package normalize

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

// challengeModeStart gets the dungeonID, keyLevel and affixes
// v17
// 1/24 16:47:38.068  CHALLENGE_MODE_START,"De Other Side",2291,377,10,[10,123,3,121]
func challengeModeStart(params []string, uploadUUID string, combatlogUUID string, rec map[string]map[string][]*timestreamwrite.WriteRecordsInput, timestamp *string) error {
	dungeonID, err := Atoi64(params[2]) // 2291
	if err != nil {
		log.Printf("failed to convert challange mode start event, field dungeon id. got: %v", params[2])
		return err
	}

	keyLevel, err := Atoi64(params[4]) // 10
	if err != nil {
		log.Printf("failed to convert challange mode start event, field key level. got: %v", params[4])
		return err
	}

	dungeonName := strings.Trim(params[1], "\"")
	currentTimeInSeconds := time.Now().Unix()

	var e = []*timestreamwrite.Record{
		{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("dungeon_name"),
					Value: aws.String(dungeonName),
				},
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
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
			},
			MeasureName:      aws.String("key_level"),
			MeasureValue:     aws.String(strconv.FormatInt(keyLevel, 10)),
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
			},
			MeasureName:      aws.String("date"),
			MeasureValue:     timestamp,
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		},
	}

	if len(params) >= 6 {
		twoAffixID, err := Atoi64(strings.Trim(strings.Trim(params[5], "["), "]")) // 10
		if err != nil {
			log.Printf("failed to convert challange mode start event, field level 2 affix. got: %v", params[5])
			return err
		}
		r := &timestreamwrite.Record{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
				},
			},
			MeasureName:      aws.String("two_affix_id"),
			MeasureValue:     aws.String(strconv.FormatInt(twoAffixID, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		}
		e = append(e, r)
	}
	if len(params) >= 7 {
		fourAffixID, err := Atoi64(strings.Trim(params[6], "]")) // 123
		if err != nil {
			log.Printf("failed to convert challange mode start event, field level 4 affix. got: %v", params[6])
			return err
		}
		r := &timestreamwrite.Record{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
				},
			},
			MeasureName:      aws.String("four_affix_id"),
			MeasureValue:     aws.String(strconv.FormatInt(fourAffixID, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		}
		e = append(e, r)
	}
	if len(params) >= 8 {
		sevenAffixID, err := Atoi64(strings.Trim(params[7], "]")) // 3
		if err != nil {
			log.Printf("failed to convert challange mode start event, field level 7 affix. got: %v", params[7])
			return err
		}
		r := &timestreamwrite.Record{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
				},
			},
			MeasureName:      aws.String("seven_affix_id"),
			MeasureValue:     aws.String(strconv.FormatInt(sevenAffixID, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		}
		e = append(e, r)
	}
	if len(params) == 9 {
		tenAffixID, err := Atoi64(strings.Trim(params[8], "]")) // 121
		if err != nil {
			log.Printf("failed to convert challange mode start event, field level 10 affix. got: %v", params[8])
			return err
		}
		r := &timestreamwrite.Record{
			Dimensions: []*timestreamwrite.Dimension{
				{
					Name:  aws.String("upload_uuid"),
					Value: aws.String(uploadUUID),
				},
			},
			MeasureName:      aws.String("ten_affix_id"),
			MeasureValue:     aws.String(strconv.FormatInt(tenAffixID, 10)),
			MeasureValueType: aws.String("BIGINT"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		}
		e = append(e, r)
	}

	rand.Seed(time.Now().UnixNano())
	key := strconv.Itoa(rand.Int())

	writeRecordsInputs := []*timestreamwrite.WriteRecordsInput{
		{
			// common attributes don't matter here, because this is a very rare
			// event
			CommonAttributes: &timestreamwrite.Record{
				Dimensions: []*timestreamwrite.Dimension{},
			},
			Records: e,
		},
	}
	rec[combatlogUUID][key] = writeRecordsInputs

	return nil
}
