package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Atoi32 converts a string directly to a int32, baseline golang parses string always into int64 and have to be converted
//to int32. You can however transform a string easily to int, which is somehow the same, but the parquet package expects int32
//specifically
func Atoi32(input string) (int32, error) {
	bigint, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, err
	}

	num := int32(bigint)
	return num, nil
}

//Atoi64 is just a small wrapper around ParseInt
func Atoi64(input string) (int64, error) {
	num, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}

	return num, nil
}

//NOTE: propably should check that it is surounded by quotes and fail otherwise
//		to make it fail early.
//		Because the columns must be surrounded by qoutes otherwise it is a wrong column
func trimQuotes(input string) string {
	output := strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

func convertToTimestampMilli(input string) (int64, error) {
	//TODO: this will break during new year, because go assumes UTC,
	//		but the combatlog has the time of the player afaik
	input = fmt.Sprintf("%v/%s", time.Now().Year(), input)
	stupidBlizzTimeformat := "2006/1/2 15:04:05.000"
	t, err := time.Parse(stupidBlizzTimeformat, input)
	if err != nil {
		return 0, err
	}

	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)), nil
}

//copying code from stackoverflow like a pro
//https://stackoverflow.com/questions/59297737/go-split-string-by-comma-but-ignore-comma-within-double-quotes
func splitAtCommas(s *string) []string {
	var res []string
	var beg int
	var inString bool

	for i := 0; i < len(*s); i++ {
		if (*s)[i] == ',' && !inString {
			res = append(res, (*s)[beg:i])
			beg = i + 1
		} else if (*s)[i] == '"' {
			if !inString {
				inString = true
			} else if i > 0 && (*s)[i-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, (*s)[beg:])
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
