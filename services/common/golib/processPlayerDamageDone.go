package golib

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*
func PlayerDamageSimpleResponseToJson2(result *dynamodb2.QueryOutput, sorted, firstPage bool) (string, error) {
	var items []DynamoDBPlayerDamageSimple
	var lastPage bool //controls if the next btn is shown

	err := attributevalue2.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		return "", err
	}
	if result.LastEvaluatedKey != nil {
		items = items[:len(items)-1] //we don't care about the last item
		//getting 1 more item than we need and truncaten it off allows us to circumvent a dynamodb2 bug where
		//if 5 items are left to return and you get exactly return 5 items the last evaluated key is not null,
		//which leads to the next request being completely empty
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

	var firstSk types2.AttributeValue
	var lastSk types2.AttributeValue

	if sorted == false {
		firstSk = result.Items[0]["sk"]
		lastSk = result.Items[len(items)-1]["sk"]
		//-1 would be the last item, which we truncated
		//we need to 2nd last as last key, for the next page to not skip an item
	} else {
		//if we go back on the pagination the order is reveresed
		//which causes last and first key to be switched, if we
		//were to go back again, thing are messed up
		firstSk = result.Items[len(items)-1]["sk"]
		lastSk = result.Items[0]["sk"]
	}
	logrus.Debug(firstSk)
	logrus.Debug(lastSk)

	resp := JSONPlayerDamageSimpleResponse2{
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
*/

// PlayerDamageDoneToJson returns the log specific damage result, including damage
// per spell breakdown, both damage per player and per spell are sorted before saving
// to the db.
// We don't need an extra JSON struct, like for keys, because there is no
// pagination etc.
func PlayerDamageDoneToJson(result *dynamodb.GetItemOutput) (string, error) {
	var item DynamoDBPlayerDamageDone

	err := dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	return string(b), err
}
