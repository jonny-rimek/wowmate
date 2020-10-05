package normalize

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

//TODO: check that it is surounded by quotes and fail otherwise
//		to make it fail early.
//		Because the columns must be surrounded by qoutes otherwise it is a wrong column
func trimQuotes(input string) string {
	output := strings.TrimSuffix(input, "\"")
	output = strings.TrimPrefix(output, "\"")
	return output
}

//TODO: test and fix the year problem
// this will break during new year, because go assumes UTC,
// but the combatlog has the time of the player afaik
func convertToTimestampMilli(input string) (int64, error) {
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
//atleast I added tests^^ and switched to string pointers to reduce memory
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

func convertToCSV(events *[]Event) (io.Reader, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	ss, err := EventsAsStringSlices(events)
	if err != nil {
		return nil, err
	}
	log.Println("converted to struct to string slice")

	//flushes the string slice as csv to buffer
	if err := w.WriteAll(ss); err != nil {
		return nil, err
	}
	log.Println("converted to csv")

	return io.Reader(&buf), nil
}

//TODO: test once the db table definition is stable
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

//IMPROVE:
//there shouldn't be any aws logic inside this package, but I want to upload directly after every
//m+, the idea was to free the memory of the uploaded combatlog, but I don't think this is working
//anyway
func uploadS3(r *io.Reader, sess *session.Session, mythicplugUUID string, csvBucket string) error {
	if mythicplugUUID == "" {
		//sometimes there are more CHALLANGE_MODE_END events than there are start events
		//it shouldn't come to this, because we aren't adding anything unless we have a started event
		return nil
	}
	uploader := s3manager.NewUploader(sess)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(csvBucket),
		Key:    aws.String(fmt.Sprintf("%v.csv", mythicplugUUID)),
		Body:   *r,
	})
	if err != nil {
		log.Println("Failed to upload to S3")
		return err
	}
	log.Println("Upload finished! location: " + result.Location)

	return nil
}

//IMPROVE: add all fields once table design is stable
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
