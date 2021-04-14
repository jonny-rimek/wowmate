package normalize

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

// v16
// 10/3 05:44:30.076  COMBAT_LOG_VERSION,16,ADVANCED_LOG_ENABLED,1,BUILD_VERSION,9.0.2,PROJECT_ID,1
func combatLogVersion(params []string, uploadUUID string, combatlogUUID string, rec map[string]map[string][]*timestreamwrite.WriteRecordsInput) error {
	if len(params) != 8 {
		return fmt.Errorf("combatlog version should have 7 columns, it has %v: %v", len(params), params)
	}

	patchVersion := params[5]

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
			MeasureName:      aws.String("patch_version"),
			MeasureValue:     aws.String(patchVersion),
			MeasureValueType: aws.String("VARCHAR"),
			Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
			TimeUnit:         aws.String("SECONDS"),
		},
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
