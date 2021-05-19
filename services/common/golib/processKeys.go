package golib

import (
	"encoding/json"
	"sort"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sirupsen/logrus"
)

type JSONKeysResponse struct {
	Data      []JSONKeys `json:"data"`
	FirstSK   string     `json:"first_sk"`
	LastSK    string     `json:"last_sk"`
	FirstPage bool       `json:"first_page"`
	LastPage  bool       `json:"last_page"`
}

type JSONKeys struct {
	Damage        []PlayerDamage `json:"player_damage"`
	Duration      string         `json:"duration"`
	Deaths        int            `json:"deaths"`
	Affixes       string         `json:"affixes"`
	Keylevel      int            `json:"keylevel"`
	DungeonName   string         `json:"dungeon_name"`
	DungeonID     int            `json:"dungeon_id"`
	CombatlogHash string         `json:"combatlog_hash"`
	Intime        int            `json:"intime"`
	// don't need date atm, readd if needed
	// Date          int64          `json:"date"`
}

// KeysResponseToJson takes a dynamodb query output and converts it to be consumed by the frontend
// mostly it makes sure the pagination works correctly
// this is used for top keys and top keys per dungeon
func KeysResponseToJson(result *dynamodb.QueryOutput, sorted, firstPage bool) (string, error) {
	var items []DynamoDBKeys
	var lastPage bool // controls if the next btn is shown

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return "", err
	}
	if result.LastEvaluatedKey != nil {
		items = items[:len(items)-1] // we don't care about the last item
		// getting 1 more item than we need and truncate it off allows us to circumvent a dynamodb2 bug where
		// if 5 items are left to return and you get exactly return 5 items the last evaluated key is not null,
		// which leads to the next request being completely empty
		lastPage = false
	} else if result.LastEvaluatedKey == nil && sorted == true {
		firstPage = true
		lastPage = false
	} else if result.LastEvaluatedKey == nil && sorted == false {
		lastPage = true
	}

	if sorted == true {
		sort.Slice(items, func(i, j int) bool {
			return items[i].Sk > items[j].Sk // order descending
		})
	}

	var r []JSONKeys

	for _, el := range items {
		resp := JSONKeys{
			Damage:        el.Damage,
			Duration:      el.Duration,
			Deaths:        el.Deaths,
			Affixes:       el.Affixes,
			Keylevel:      el.Keylevel,
			DungeonName:   el.DungeonName,
			DungeonID:     el.DungeonID,
			CombatlogHash: el.CombatlogHash,
			Intime:        el.Intime,
		}
		r = append(r, resp)
	}
	logrus.Debug(r)

	var firstSk string
	var lastSk string

	if sorted == false {
		firstSk = *result.Items[0]["sk"].S
		lastSk = *result.Items[len(items)-1]["sk"].S
		// -1 would be the last item, which we truncated
		// we need to 2nd last as last key, for the next page to not skip an item
	} else {
		// if we go back on the pagination the order is reversed
		// which causes last and first key to be switched, if we
		// were to go back again, thing are messed up
		firstSk = *result.Items[len(items)-1]["sk"].S
		lastSk = *result.Items[0]["sk"].S
	}
	logrus.Debug(firstSk)
	logrus.Debug(lastSk)

	resp := JSONKeysResponse{
		Data:      r,
		FirstSK:   firstSk,
		LastSK:    lastSk,
		FirstPage: firstPage,
		LastPage:  lastPage,
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(b), err
}

// KeysToJson is unused right now, I think it was the predecessor of KeysResponseToJson
func KeysToJson(result *dynamodb.QueryOutput) (string, error) {
	var items []DynamoDBKeys

	err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return "", err
	}

	var r []JSONKeys

	for _, el := range items {
		resp := JSONKeys{
			Damage:      el.Damage,
			Duration:    el.Duration,
			Deaths:      el.Deaths,
			Affixes:     el.Affixes,
			Keylevel:    el.Keylevel,
			DungeonName: el.DungeonName,
			DungeonID:   el.DungeonID,
		}
		r = append(r, resp)
	}
	logrus.Debug(r)

	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), err
}
