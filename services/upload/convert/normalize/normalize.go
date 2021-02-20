package normalize

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/gofrs/uuid"
	"golang.org/x/net/http2"
)

//Normalize converts the combatlog to a slice of Event structs
func Normalize(scanner *bufio.Scanner, uploadUUID string, sess *session.Session) error {
	var combatEvents []*timestreamwrite.Record
	//IMPROVE: the UploadUUID logic should be part of the normalize package
	//UploadUUID //for the whole file
	combatlogUUID := "" //after every COMBAT_LOG_VERSION
	// BossFightUUID := ""
	// MythicplusUUID := ""

	//combatEvents = make([]Event, 0, 100000) //100.000 is an arbitrary value
	//initialising the slice with a capacitiy to reduce the amount reallocations
	//the difference in a small log was <1sec -> not worth

	for scanner.Scan() {
		//4/24 10:42:30.561  COMBAT_LOG_VERSION
		//every line starts with the date followed by the rest seperated with 2 spaces.
		//the rest is seperated with commas

		//IMPROVE:
		//write version of .Split that accepts a pointer to the string
		//this saved a lot of memory with splitAtComma
		//just look at the implementation of the strings.Split function
		//maybe there is a package that implements string functionality more efficient
		//the main problem is that the strings.Split calls a bunch of other functions
		//that all create a new version of the string and thus bloating the memory
		//UPDATE: it's a lot of work to rewrite everything to use []byte I'll resist the
		//premature optimization for now
		row := splitString(scanner.Text(), "  ")

		//NOTE: not written to DB atm https://github.com/jonny-rimek/wowmate/issues/129
		//gonna use it later again
		// timestamp, err := timestampMilli(row[0])
		// if err != nil {
		// 	return err
		// }

		params := splitAtCommas(&row[1])

		// e := &timestreamwrite.Record{}
		// e := Event{
		// 	// UploadUUID:     uploadUUID,
		// 	// CombatlogUUID:  CombatlogUUID,
		// 	// BossFightUUID:  BossFightUUID,
		// 	MythicplusUUID: MythicplusUUID,
		// 	// ColumnUUID:     uuid.Must(uuid.NewV4()).String(),
		// 	Timestamp: timestamp,
		// 	// EventType:      params[0],
		// }

		//TODO: never add anything if CombatlogUUID is empty, same logic as m+uuid

		// if MythicplusUUID == "" && params[0] != "CHALLENGE_MODE_START" {
		//I don't want to add events if they are outside of a combatlog
		// 	continue
		// }

		switch params[0] {
		case "COMBAT_LOG_VERSION":
			//TODO:
			//- [x] check version
			//- [x] check advanced logging
			//- [ ] generate report uuid
			//		- not sure what this is about

			combatlogUUID = uuid.Must(uuid.NewV4()).String()
			// e.CombatlogUUID = CombatlogUUID
			// err = e.combatLogVersion(params)
			// if err != nil {
			// 	return err
			// }
			//NOTE:
			//break is implicit in go, that means after the first match it exits
			//the switch statement

		case "ENCOUNTER_START":
			// e.BossFightUUID = uuid.Must(uuid.NewV4()).String()
			// err = e.encounterStart(params)
			// if err != nil {
			// 	return err
			// }

		case "ENCOUNTER_END":
			//I want the entry with encounter_end to have the id, just the records after should be nil again
			//it's been already set for this record (e) so it's okay to clear it befor calling encounterEnd
			// BossFightUUID = ""
			// err = e.encounterEnd(params)
			// if err != nil {
			// 	return err
			// }

		case "CHALLENGE_MODE_START":
			// MythicplusUUID = uuid.Must(uuid.NewV4()).String()
			// e.MythicplusUUID = MythicplusUUID
			// err = e.challengeModeStart(params)
			// if err != nil {
			// 	return err
			// }

		case "CHALLENGE_MODE_END":
			// err = e.challengeModeEnd(params)
			// if err != nil {
			// 	return err
			// }
			// combatEvents = append(combatEvents, e)

			// r, err := convertToCSV(&combatEvents)
			// if err != nil {
			// 	return err
			// }
			// err = uploadS3(&r, sess, e.MythicplusUUID, csvBucket)
			// if err != nil {
			// 	return err
			// }

			// combatEvents = nil
			// MythicplusUUID = ""

		// case "SPELL_HEAL", "SPELL_PERIODIC_HEAL":
		// 	err = e.importHeal(params)
		case "SPELL_DAMAGE":
			// log.Println("inside spell damage")
			e, err := spellDamage(params, uploadUUID, combatlogUUID)
			if err != nil {
				return err
			}
			combatEvents = append(combatEvents, e)

		default:
			// e.Unsupported = true
		}

		// if params[0] == "CHALLENGE_MODE_END" {
		// 	continue
		// }
	}

	err := uploadToTimestream(combatEvents)
	if err != nil {
		return err
	}

	return nil
}

func uploadToTimestream(e []*timestreamwrite.Record) error {
	log.Printf("%v dmg records: ", len(e))

	tr := &http.Transport{
		ResponseHeaderTimeout: 20 * time.Second,
		// Using DefaultTransport values for other parameters: https://golang.org/pkg/net/http/#RoundTripper
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			KeepAlive: 30 * time.Second,
			DualStack: true,
			Timeout:   30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// So client makes HTTP/2 requests
	err := http2.ConfigureTransport(tr)
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1"), MaxRetries: aws.Int(10), HTTPClient: &http.Client{Transport: tr}})
	if err != nil {
		return err
	}

	for i := 0; i < len(e); i += 100 {

		//get the upper bound of the record to write, in case it is the
		//last bit of records and i + 99 does not exist
		j := 0
		if i+99 > len(e) {
			j = len(e)
		} else {
			j = i + 99
		}

		//use common batching https://docs.aws.amazon.com/timestream/latest/developerguide/metering-and-pricing.writes.html#metering-and-pricing.writes.write-size-multiple-events
		//probably only applies to the uploadUuid tho
		writeSvc := timestreamwrite.New(sess)
		writeRecordsInput := &timestreamwrite.WriteRecordsInput{
			DatabaseName: aws.String("wowmate-analytics"),
			TableName:    aws.String("combatlogs"),
			Records:      e[i:j], //only upload a part of the records
		}

		_, err = writeSvc.WriteRecords(writeRecordsInput)
		if err != nil {
			return err
		}
	}
	log.Println("Write records is successful")
	return nil
}
