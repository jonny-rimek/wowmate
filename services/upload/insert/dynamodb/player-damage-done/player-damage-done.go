package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
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
		//goland:noinspection ALL
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
func convertQueryResult(queryResult *timestreamquery.QueryOutput) (golib.DynamoDBPlayerDamageDone, error) {
	resp := golib.DynamoDBPlayerDamageDone{}

	players := getPlayers(queryResult)

	m := make(map[string]golib.PlayerDamageDone)

	for _, row := range queryResult.Rows {
		for _, player := range players {
			if player == *row.Data[1].ScalarValue {

				dam, err := golib.Atoi64(*row.Data[0].ScalarValue)
				if err != nil {
					return resp, fmt.Errorf("failed to convert time damage cell to int for first entry in map: %v", err.Error())
				}
				spellID, err := strconv.Atoi(*row.Data[13].ScalarValue)
				if err != nil {
					return resp, fmt.Errorf("failed to convert time damage cell to int for first entry in map: %v", err.Error())
				}

				// check if map has entry with the name of the player
				if val, ok := m[player]; ok {
					// yes, map contains player
					// add damage value to existing struct in map
					d := golib.DamagePerSpell{
						SpellID:   spellID,
						SpellName: *row.Data[14].ScalarValue,
						Damage:    dam,
					}

					val.DamagePerSpell = append(val.DamagePerSpell, d)
					val.Damage += dam
					m[player] = val

				} else {
					// no, map does not contain player as key > initialize
					d := golib.PlayerDamageDone{
						Damage:   dam,
						Name:     *row.Data[1].ScalarValue,
						PlayerID: *row.Data[2].ScalarValue,
						Class:    "unsupported",
						Specc:    "unsupported",
						DamagePerSpell: []golib.DamagePerSpell{
							{
								SpellID:   spellID,
								SpellName: *row.Data[14].ScalarValue,
								Damage:    dam,
							},
						},
					}

					m[player] = d
				}
			}
		}
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
	// converts duration to date 1970 + duration, of which I only display the minutes and seconds
	// time.Duration, doesn't allow mixed formatting like min:seconds
	t := time.Unix(0, durationInMilliseconds*1e6) // milliseconds > nanoseconds

	_, intime, err := golib.TimedAsPercent(dungeonID, float64(durationInMilliseconds))
	if err != nil {
		return resp, err
	}

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

	date, err := golib.Atoi64(*queryResult.Rows[0].Data[15].ScalarValue)
	if err != nil {
		return resp, err
	}

	resp = golib.DynamoDBPlayerDamageDone{
		Pk:            fmt.Sprintf("LOG#KEY#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID),
		Sk:            fmt.Sprintf("LOG#KEY#%v#OVERALL_PLAYER_DAMAGE", combatlogUUID),
		Duration:      t.Format("04:05"), // formats to minutes:seconds
		Deaths:        0,
		Affixes:       golib.AffixIDsToString(twoAffixID, fourAffixID, sevenAffixID, tenAffixID),
		Keylevel:      keyLevel,
		DungeonName:   dungeonName,
		DungeonID:     dungeonID,
		CombatlogUUID: combatlogUUID,
		Finished:      finished != 0, // if 0 false, else 1
		Intime:        intime,
		Date:          date,
	}

	for _, el := range m {
		resp.Damage = append(resp.Damage, el)
	}
	// sort player damage desc
	sort.Slice(resp.Damage, func(i, j int) bool {
		return resp.Damage[i].Damage > resp.Damage[j].Damage // order descending
	})

	// sort spell damage desc
	for _, el := range resp.Damage {
		sort.Slice(el.DamagePerSpell, func(i, j int) bool {
			return el.Damage > el.Damage // order descending
		})
	}

	// TODO: merge arrays of DamagePerSpell by name and keep one id

	// prettyStruct, err := golib.PrettyStruct(resp)
	// if err != nil {
	// 	return golib.DynamoDBPlayerDamageDone{}, err
	// }
	// log.Println(prettyStruct)

	return resp, err
}

// getPlayers checks the query output and returns a list of all players that are in it
func getPlayers(result *timestreamquery.QueryOutput) []string {
	var players []string

	for _, el := range result.Rows {
		player := *el.Data[1].ScalarValue
		if contains(players, player) == false {
			players = append(players, player)
		}
	}

	return players
}

// contains dd slice
func contains(slice []string, el string) bool {
	for _, a := range slice {
		if a == el {
			return true
		}
	}
	return false
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
