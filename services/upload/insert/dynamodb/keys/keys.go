package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/timestreamquery"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/jonny-rimek/wowmate/services/common/golib"
	"github.com/sirupsen/logrus"
)

type logData struct {
	Wcu float64
}

var svc *dynamodb.DynamoDB

func handler(ctx aws.Context, e events.SNSEvent) error {
	logData, err := handle(ctx, e)
	if err != nil {
		//goland:noinspection GoNilness
		golib.CanonicalLog(map[string]interface{}{
			"wcu":   logData.Wcu,
			"err":   err.Error(),
			"event": e,
		})
		return err
	}

	golib.CanonicalLog(map[string]interface{}{
		"wcu": logData.Wcu,
	})
	return err
}

func handle(ctx aws.Context, e events.SNSEvent) (logData, error) {
	var logData logData

	ddbTableName := os.Getenv("DYNAMODB_TABLE_NAME")
	if ddbTableName == "" {
		return logData, fmt.Errorf("dynamo db table name env var is empty")
	}

	queryResult, err := extractQueryResult(e)

	record, err := convertQueryResult(queryResult)
	if err != nil {
		return logData, err
	}

	response, err := golib.DynamoDBPutItem(ctx, svc, &ddbTableName, record)
	if err != nil {
		return logData, err
	}

	logData.Wcu = *response.ConsumedCapacity.CapacityUnits

	return logData, nil
}

func extractQueryResult(e events.SNSEvent) (*timestreamquery.QueryOutput, error) {
	message := e.Records[0].SNS.Message
	if message == "" {
		return nil, fmt.Errorf("message is empty")
	}

	var result *timestreamquery.QueryOutput

	err := json.Unmarshal([]byte(message), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal sns message which contains the query result: %v", err)
	}

	return result, err
}

// TODO: tests
func convertQueryResult(queryResult *timestreamquery.QueryOutput) (golib.DynamoDBKeys, error) {
	resp := golib.DynamoDBKeys{}

	var summaries []golib.PlayerDamage

	for i := 0; i < len(queryResult.Rows); i++ {
		dam, err := strconv.Atoi(*queryResult.Rows[i].Data[0].ScalarValue)
		if err != nil {
			return resp, err
		}

		d := golib.PlayerDamage{
			Damage:   dam,
			Name:     *queryResult.Rows[i].Data[1].ScalarValue,
			PlayerID: *queryResult.Rows[i].Data[2].ScalarValue,
			Class:    "unsupported",
			Specc:    "unsupported",
		}

		summaries = append(summaries, d)
	}
	combatlogUUID := *queryResult.Rows[0].Data[3].ScalarValue

	dungeonName := *queryResult.Rows[0].Data[4].ScalarValue

	dungeonID, err := strconv.Atoi(*queryResult.Rows[0].Data[5].ScalarValue)
	if err != nil {
		return resp, err
	}

	keyLevel, err := strconv.Atoi(*queryResult.Rows[0].Data[6].ScalarValue)
	if err != nil {
		return resp, err
	}

	durationInMilliseconds, err := golib.Atoi64(*queryResult.Rows[0].Data[7].ScalarValue)
	if err != nil {
		return resp, err
	}
	dur := float64(durationInMilliseconds)
	// converts duration to date 1970 + duration, of which I only display the minutes and seconds
	// time.Duration, doesn't allow mixed formatting like min:seconds
	t := time.Unix(0, durationInMilliseconds*1e6) // milliseconds > nanoseconds

	finished, err := strconv.Atoi(*queryResult.Rows[0].Data[8].ScalarValue)
	if err != nil {
		return resp, err
	}

	twoAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[9].ScalarValue)
	if err != nil {
		return resp, err
	}

	fourAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[10].ScalarValue)
	if err != nil {
		return resp, err
	}

	sevenAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[11].ScalarValue)
	if err != nil {
		return resp, err
	}

	tenAffixID, err := strconv.Atoi(*queryResult.Rows[0].Data[12].ScalarValue)
	if err != nil {
		return resp, err
	}

	patch := *queryResult.Rows[0].Data[13].ScalarValue

	durAsPercent, intime, err := timedAsPercent(dungeonID, dur)
	if err != nil {
		return resp, err
	}

	resp = golib.DynamoDBKeys{
		// hardcoding the patch like that might be too granular, maybe it makes more sense that e.g. 9.0.2 and 9.0.5 are both S1
		Pk: fmt.Sprintf("LOG#KEY#%s", patch),
		Sk: fmt.Sprintf("%02d#%3.6f#%v", keyLevel, durAsPercent, combatlogUUID),
		// sorting in dynamoDB is achieved via the sort key, in order to sort by key level and within the key level by
		// time I'm printing the value as string and sort the string.
		// As I'm sorting descending I can't just print the duration in milliseconds.
		// instead I print the duration as percent in relation to the intime duration
		Damage:        summaries,
		Gsi1pk:        fmt.Sprintf("LOG#KEY#%s#%v", patch, dungeonID),
		Gsi1sk:        fmt.Sprintf("%02d#%3.6f#%v", keyLevel, durAsPercent, combatlogUUID),
		Duration:      t.Format("04:05"), // formats to minutes:seconds
		Deaths:        0,                 // TODO:
		Affixes:       golib.AffixIDsToString(twoAffixID, fourAffixID, sevenAffixID, tenAffixID),
		Keylevel:      keyLevel,
		DungeonName:   dungeonName,
		DungeonID:     dungeonID,
		CombatlogUUID: combatlogUUID,
		Finished:      finished != 0, // if 0 false, else true
		Intime:        intime,
	}
	return resp, err
}

// minSecToMilliseconds converts time in the "minute:seconds" format to milliseconds
func minSecToMilliseconds(input string) (int64, error) {
	input = fmt.Sprintf("1970 %s", input)
	t, err := time.Parse("2006 04:05", input)
	if err != nil {
		return 0, fmt.Errorf("failed to parse time input: %v", err)
	}
	milliseconds := t.UnixNano() / 1e6
	return milliseconds, nil
}

func timedAsPercent(dungeonID int, durationInMilliseconds float64) (durAsPercent float64, intime int, err error) {
	var intimeDuration, twoChestDuration, threeChestDuration float64

	switch dungeonID {
	case 2291: // De Other Side
		ms, err := minSecToMilliseconds("43:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("34:25")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("25:49")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2289: // Plaguefall
		ms, err := minSecToMilliseconds("38:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("30:24")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("22:48")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)
	case 2284: // Sanguine Depths
		ms, err := minSecToMilliseconds("41:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("32:48")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("24:36")
		if err != nil {
			return 0, 0, err
		}
	// TODO: parse time and convert to milli seconds do it in TDD
	/*
		https://www.wowhead.com/mythic-keystones-and-dungeons-guide
		Dungeon	Timer	+2	+3
		De Other Side	43:00	34:25	25:49
		Plaguefall	38:00	30:24	22:48
		Halls of Atonement	31:00	24:48	18:36
		Mists of Tirna Scithe	30:00	24:00	18:00
		Spires of Ascension	39:00	31:12	23:24
		Sanguine Depths	41:00	32:48	24:36
		Necrotic Wake	36:00	28:48	21:36
		Theater of Pain	37:00	29:36	22:12
	*/
	case 2287: // Halls of Atonement
		ms, err := minSecToMilliseconds("31:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("24:48")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("18:36")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2290: // Mists of Tirna Scithe
		ms, err := minSecToMilliseconds("30:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("24:00")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("18:00")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2285: // Spires of Ascension
		ms, err := minSecToMilliseconds("39:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("31:12")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("23:24")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2286: // Necrotic Wake
		ms, err := minSecToMilliseconds("36:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("28:48")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("21:36")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)

	case 2293:
		ms, err := minSecToMilliseconds("37:00")
		if err != nil {
			return 0, 0, err
		}
		intimeDuration = float64(ms)

		ms, err = minSecToMilliseconds("29:36")
		if err != nil {
			return 0, 0, err
		}
		twoChestDuration = float64(ms)

		ms, err = minSecToMilliseconds("22:12")
		if err != nil {
			return 0, 0, err
		}
		threeChestDuration = float64(ms)
	}
	intime = timed(durationInMilliseconds, intimeDuration, twoChestDuration, threeChestDuration)

	return durationAsPercent(intimeDuration, durationInMilliseconds), intime, err
}

func timed(durationInMilliseconds, intimeDuration, twoChestDuration, threeChestDuration float64) int {
	if durationInMilliseconds <= threeChestDuration {
		return 3 // three chest
	} else if durationInMilliseconds > threeChestDuration && durationInMilliseconds <= twoChestDuration {
		return 2 // two chest
	} else if durationInMilliseconds > twoChestDuration && durationInMilliseconds <= intimeDuration {
		return 1 // timed
	} else {
		return 0 // deplete
	}
}

func durationAsPercent(dungeonIntimeDuration, durationInMilliseconds float64) float64 {
	return (dungeonIntimeDuration / durationInMilliseconds) * 100
}

func main() {
	golib.InitLogging()

	sess, err := session.NewSession()
	if err != nil {
		logrus.Info(fmt.Sprintf("Error creating session: %v", err.Error()))
		return
	}

	svc = dynamodb.New(sess)
	if os.Getenv("LOCAL") == "false" {
		xray.AWS(svc.Client)
	}

	lambda.Start(handler)
}
