package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//Atoi32 converts a string directly to a int32, baseline golang parses string always into int64 and have to be converted
//to int32. You can however transform a string easily to int, which is somehow the same, but the parquet package expects int32
//specifically
func Atoi32(input string) (num int32, err error) {
	bigint, err := strconv.ParseInt(input, 10, 32)
	num = int32(bigint)
	return num, err
}

//Atoi64 is just a small wrapper around ParseInt
func Atoi64(input string) (num int64, err error) {
	num, err = strconv.ParseInt(input, 10, 64)
	return num, err
}

func trimQuotes(input string) (output string) {
	output = strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

func convertToTimestampMilli(input string) (int64, error) {
	input = fmt.Sprintf("%v/%s", time.Now().Year(), input)
	stupidBlizzTimeformat := "2006/1/2 15:04:05.000"
	t, err := time.Parse(stupidBlizzTimeformat, input)
	if err != nil {
		return 0, err
	}

	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)), nil

}

//11/3 09:00:00.760  CHALLENGE_MODE_START,"Atal'Dazar",1763,244,10,[10,11,14,16]

//v16
//10/3 05:51:00.975  CHALLENGE_MODE_START,"Halls of Atonement",2287,378,2,[10]
//NOTE: the array is definitely the affixes
func (e *Event) normalizeChallengeModeStart(params []string) (err error) {
	if len(params) != 6 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.DungeonName = trimQuotes(params[1])
	e.DungeonID, err = Atoi32(params[2])
	if err != nil {
		log.Println("failed to convert challange mode start event")
		return err
	}
	e.KeyUnkown1, err = Atoi32(params[3])
	if err != nil {
		log.Println("failed to convert challange mode start event")
		return err
	}
	e.KeyLevel, err = Atoi32(params[4])
	if err != nil {
		log.Println("failed to convert challange mode start event")
		return err
	}
	e.KeyArray = params[5]
	return nil
}

//11/3 09:34:07.310  CHALLENGE_MODE_END,1763,1,10,2123441

//v16
//10/3 05:51:00.879  CHALLENGE_MODE_END,2287,0,0,0
func (e *Event) normalizeChallengeModeEnd(params []string) (err error) {
	if len(params) != 5 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.DungeonID, err = Atoi32(params[1])
	if err != nil {
		log.Println("failed to convert challange mode end event")
		return err
	}
	e.KeyUnkown1, err = Atoi32(params[2])
	if err != nil {
		log.Println("failed to convert challange mode end event")
		return err
	}
	e.KeyLevel, err = Atoi32(params[3])
	if err != nil {
		log.Println("failed to convert challange mode end event")
		return err
	}
	e.KeyDuration, err = Atoi64(params[3])
	if err != nil {
		log.Println("failed to convert challange mode end event")
		return err
	}
	return nil
}

//11/3 09:00:22.354  ENCOUNTER_START,2086,"Rezan",8,5,1763

//16
//10/3 05:59:07.379  ENCOUNTER_START,2401,"Halkias, the Sin-Stained Goliath",8,5,2287
func (e *Event) normalizeEncounterStart(params []string) (err error) {
	if len(params) != 6 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.EncounterID, err = Atoi32(params[1])
	if err != nil {
		log.Println("failed to convert encounter start event")
		return err
	}

	e.EncounterName = trimQuotes(params[2])

	e.EncounterUnknown1, err = Atoi32(params[3])
	if err != nil {
		log.Println("failed to convert encounter start event")
		return err
	}

	e.EncounterUnknown2, err = Atoi32(params[4])
	if err != nil {
		log.Println("failed to convert encounter start event")
		return err
	}

	e.DungeonID, err = Atoi32(params[5])
	if err != nil {
		log.Println("failed to convert encounter start event")
		return err
	}
	return nil
}

//11/3 09:01:58.364  ENCOUNTER_END,2086,"Rezan",8,5,1

//v16
//10/3 06:00:02.433  ENCOUNTER_END,2401,"Halkias, the Sin-Stained Goliath",8,5,1
func (e *Event) normalizeEncounterEnd(params []string) (err error) {
	if len(params) != 6 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.EncounterID, err = Atoi32(params[1])
	if err != nil {
		log.Println("failed to convert encounter end event")
		return err
	}

	e.EncounterName = trimQuotes(params[2])

	e.EncounterUnknown1, err = Atoi32(params[3])
	if err != nil {
		log.Println("failed to convert encounter end event")
		return err
	}

	e.EncounterUnknown2, err = Atoi32(params[4])
	if err != nil {
		log.Println("failed to convert encounter end event")
		return err
	}

	e.Killed, err = Atoi32(params[5])
	if err != nil {
		log.Println("failed to convert encounter end event")
		return err
	}
	return nil
}

//v16
//10/3 05:44:30.076  COMBAT_LOG_VERSION,16,ADVANCED_LOG_ENABLED,1,BUILD_VERSION,9.0.2,PROJECT_ID,1
func (e *Event) normalizeCombatlogVersion(params []string) (err error) {
	if len(params) != 8 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.Version, err = Atoi32(params[1])
	if err != nil {
		log.Println("failed to convert combatlog version event")
		return err
	}
	if e.Version != 16 {
		return fmt.Errorf("unsupported combatlog version: %v, only version 16 is supported", e.Version)
	}
	e.AdvancedLogEnabled, err = Atoi32(params[3])
	if err != nil {
		log.Println("failed to convert combatlog version event")
		return err
	}
	if e.AdvancedLogEnabled != 1 {
		return fmt.Errorf("advanced combatlogging must be enabled")
	}
	return nil
}

//11/3 09:00:29.792  SPELL_DAMAGE,Player-1302-09C8C064,"Hyrriuk-Archimonde",0x512,0x0,Vehicle-0-3892-1763-30316-122963-00005D638F,"Rezan",0x10a48,0x0,283810,"Reckless Flurry",0x1,Vehicle-0-3892-1763-30316-122963-00005D638F,0000000000000000,3600186,3811638,0,0,2700,1,0,0,0,-790.59,2265.96,935,0.8059,122,1287,1599,-1,1,0,0,0,nil,nil,nil

// v16
// 10/3 05:51:15.415  SPELL_DAMAGE,Player-4184-00130F03,"Unstaebl-Torghast",0x512,0x0,Creature-0-2085-2287-15092-165515-0005F81144,"Depraved Darkblade",0xa48,0x0,127802,"Touch of the Grave",0x20,Creature-0-2085-2287-15092-165515-0005F81144,0000000000000000,92482,96120,0,0,1071,0,3,100,100,0,-2206.68,5071.68,1663,2.1133,60,456,456,-1,32,0,0,0,nil,nil,nil
func (e *Event) normalizeDamage(params []string) (err error) {
	if len(params) != 39 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.CasterID = params[1]               //Player-1302-09C8C064 ✔
	e.CasterName = trimQuotes(params[2]) //"Hyrriuk-Archimonde" ✔
	e.CasterType = params[3]             //0x512
	e.SourceFlag = params[4]             //0x0
	e.TargetID = params[5]               //Vehicle-0-3892-1763-30316-122963-00005D638F
	e.TargetName = trimQuotes(params[6]) //"Rezan" ✔
	e.TargetType = params[7]             //0x10a48
	e.DestFlag = params[8]               //0x0
	e.SpellID, err = Atoi32(params[9])   //283810
	if err != nil {
		log.Println("failed to convert damage event")
		return err
	}
	e.SpellName = trimQuotes(params[10])           //"Reckless Flurry" ✔
	e.SpellType = params[11]                       //0x1
	e.AnotherPlayerID = params[12]                 //Vehicle-0-3892-1763-30316-122963-00005D638F
	e.D0 = params[13]                              //0000000000000000
	e.D1, _ = strconv.ParseInt(params[14], 10, 64) //3600186
	e.D2, _ = strconv.ParseInt(params[15], 10, 64) //3811638
	e.D3, _ = strconv.ParseInt(params[16], 10, 64) //0
	e.D4, _ = strconv.ParseInt(params[17], 10, 64) //0
	e.D5, _ = strconv.ParseInt(params[18], 10, 64) //2700
	e.D6, _ = strconv.ParseInt(params[19], 10, 64) //1
	e.D7, _ = strconv.ParseInt(params[20], 10, 64) //0
	e.D8, _ = strconv.ParseInt(params[21], 10, 64) //0
	e.D9 = params[22]                              //0
	e.D10 = params[23]                             //-790.59
	e.D11 = params[24]                             //2265.96
	e.D12 = params[25]                             //935 -- mb something like a map id?
	e.D13 = params[26]                             //0.8059
	e.DamageUnkown14 = params[27]                  //122
	e.ActualAmount, err = Atoi64(params[28])       //1287
	if err != nil {
		log.Println("failed to convert damage event")
		return err
	}
	e.BaseAmount, err = Atoi64(params[29]) //1599
	if err != nil {
		log.Println("failed to convert damage event")
		return err
	}
	e.Overkill = params[30]              // ✔ -1 no overkill, otherwise the dmg number it was overkilled with. TODO: convert to int64
	e.School = params[31]                //1 ✔
	e.Crushing = params[32]              //0 always 0 with ad10-disci TODO: double check with more data NOT CONFIRMED AS crushing
	e.Blocked = params[33]               //0 TODO: always a number and should be converted to int64, pretty sure it is not blocked bc it is not reflected by actual_amount vs base_amount like absorbed
	e.Absorbed, err = Atoi64(params[34]) //0 ✔
	if err != nil {
		log.Println("failed to convert damage event")
		return err
	}
	e.Critical = params[35]  //nil ✔ fairly certain this one is crit it plays into base and actual amount, nil or 1
	e.Glancing = params[36]  //nil always nil with ad10-disci TODO double check with more data NOT CONFIRMED AS glancing
	e.IsOffhand = params[37] //nil always nil with ad10-disci TODO double check with more data NOT CONFIRMED AS is_offhand

	return nil
}

//TODO:
// func (e *Event) normalizeHeal(params []string) (err error) {
// 	e.CasterID = params[1]               //Player-970-00307C5B
// 	e.CasterName = trimQuotes(params[2]) //"Brimidreki-Sylvanas"
// 	e.CasterType = params[3]             //0x512
// 	e.SourceFlag = params[4]             //0x0
// 	e.TargetID = params[5]               //Player-970-00307C5B
// 	e.TargetName = trimQuotes(params[6]) //"Brimidreki-Sylvanas"
// 	e.TargetType = params[7]             //0x512
// 	e.DestFlag = params[8]               //0x0
// 	e.SpellID, err = Atoi32(params[9])   //122281
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.SpellName = trimQuotes(params[10])           //"Healing Elixir"
// 	e.SpellType = params[11]                       //0x8
// 	e.AnotherPlayerID = params[12]                 //Player-970-00307C5B
// 	e.D0 = params[13]                              //0000000000000000
// 	e.D1, _ = strconv.ParseInt(params[14], 10, 64) //132358
// 	e.D2, _ = strconv.ParseInt(params[15], 10, 64) //135424
// 	e.D3, _ = strconv.ParseInt(params[16], 10, 64) //4706
// 	e.D4, _ = strconv.ParseInt(params[17], 10, 64) //1467
// 	e.D5, _ = strconv.ParseInt(params[18], 10, 64) //1455
// 	e.D6, _ = strconv.ParseInt(params[19], 10, 64) //3
// 	e.D7, _ = strconv.ParseInt(params[20], 10, 64) //79
// 	e.D8, _ = strconv.ParseInt(params[21], 10, 64) //100
// 	e.D9 = params[22]                              // 0
// 	e.D10 = params[23]                             // -934.51
// 	e.D11 = params[24]                             //2149.50
// 	e.D12 = params[25]                             //3.4243
// 	e.D13 = params[26]                             //307
// 	e.DamageUnkown14 = params[27]                  //.
// 	e.ActualAmount, err = Atoi64(params[28])       //20314
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.Overhealing, err = Atoi64(params[29]) //0
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.Absorbed, err = Atoi64(params[30]) //0
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.Critical = params[31] //nil

// 	return nil
// }

//copying code from stackoverflow like a pro
//https://stackoverflow.com/questions/59297737/go-split-string-by-comma-but-ignore-comma-within-double-quotes
func splitAtCommas(s string) []string {
	res := []string{}
	var beg int
	var inString bool

	for i := 0; i < len(s); i++ {
		if s[i] == ',' && !inString {
			res = append(res, s[beg:i])
			beg = i + 1
		} else if s[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && s[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, s[beg:])
}

func EventsAsStringSlices(events *[]Event) ([][]string, error) {
	var ss [][]string

	for _, e := range *events {
		s := []string{
			e.ColumnUUID,
			e.UploadUUID,
			strconv.FormatBool(e.Unsupported),
			e.CombatlogUUID,
			e.BossFightUUID,
			e.MythicplusUUID,
			// strconv.FormatInt(e.Timestamp, 10),
			e.EventType,
			strconv.FormatInt(int64(e.Version), 10),
			strconv.FormatInt(int64(e.AdvancedLogEnabled), 10),
			e.DungeonName,
			strconv.FormatInt(int64(e.DungeonID), 10),
			strconv.FormatInt(int64(e.KeyUnkown1), 10),
			strconv.FormatInt(int64(e.KeyLevel), 10),
			e.KeyArray,
			strconv.FormatInt(e.KeyDuration, 10),
			strconv.FormatInt(int64(e.EncounterID), 10),
			e.EncounterName,
			strconv.FormatInt(int64(e.EncounterUnknown1), 10),
			strconv.FormatInt(int64(e.EncounterUnknown2), 10),
			strconv.FormatInt(int64(e.Killed), 10),
			e.CasterID,
			e.CasterName,
			e.CasterType,
			e.SourceFlag,
			e.TargetID,
			e.TargetName,
			e.TargetType,
			e.DestFlag,
			strconv.FormatInt(int64(e.SpellID), 10),
			e.SpellName,
			e.SpellType,
			strconv.FormatInt(int64(e.ExtraSpellID), 10),
			e.ExtraSpellName,
			e.ExtraSchool,
			e.AuraType,
			e.AnotherPlayerID,
			e.D0,
			strconv.FormatInt(e.D1, 10),
			strconv.FormatInt(e.D2, 10),
			strconv.FormatInt(e.D3, 10),
			strconv.FormatInt(e.D4, 10),
			strconv.FormatInt(e.D5, 10),
			strconv.FormatInt(e.D6, 10),
			strconv.FormatInt(e.D7, 10),
			strconv.FormatInt(e.D8, 10),
			e.D9,
			e.D10,
			e.D11,
			e.D12,
			e.D13,
			e.DamageUnkown14,
			strconv.FormatInt(e.ActualAmount, 10),
			strconv.FormatInt(e.BaseAmount, 10),
			strconv.FormatInt(e.Overhealing, 10),
			e.Overkill,
			e.School,
			e.Resisted,
			e.Blocked,
			strconv.FormatInt(e.Absorbed, 10),
			e.Critical,
			e.Glancing,
			e.Crushing,
			e.IsOffhand,
		}
		ss = append(ss, s)
	}
	return ss, nil
}

func (s *Event) String() string {
	return fmt.Sprintf(`[
  UUID            -> %s
  TimeStamp       -> %v
  EventType       -> %s
  CasterID        -> %s
  CasterName      -> %s
  CasterType      -> %s
  SourceFlag      -> %s
  TargetID        -> %s
  TargetName      -> %s
  TargetType      -> %s
  DestFlag        -> %s
  SpellID         -> %v
  SpellName       -> %s
  SpellType       -> %s
]
`, s.UploadUUID, s.Timestamp, s.EventType, s.CasterID, s.CasterName, s.CasterType, s.SourceFlag, s.TargetID, s.TargetName, s.TargetType, s.DestFlag, s.SpellID, s.SpellName, s.SpellType)
	//  AnotherPlayerID -> %s
}
