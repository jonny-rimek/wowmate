package normalize

import (
	"bufio"

	"github.com/aws/aws-sdk-go/aws/session"
	uuid "github.com/gofrs/uuid"
)

//Normalize converts the combatlog to a slice of Event structs
func Normalize(scanner *bufio.Scanner, uploadUUID string, sess *session.Session, csvBucket string) error {
	var combatEvents []Event
	//IMPROVE: the UploadUUID logic should be part of the normalize package
	//UploadUUID //for the whole file
	CombatlogUUID := "" //after every COMBAT_LOG_VERSION
	BossFightUUID := ""
	MythicplusUUID := ""

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
		//UPDATE: it's a lot of work to rewrite everything touse []byte I'll resist the 
		//premature optimization for now
		row := splitString(scanner.Text(), "  ")

		//NOTE: not written to DB atm https://github.com/jonny-rimek/wowmate/issues/129
		timestamp, err := timestampMilli(row[0])
		if err != nil {
			return err
		}

		params := splitAtCommas(&row[1])

		e := Event{
			UploadUUID:     uploadUUID,
			CombatlogUUID:  CombatlogUUID,
			BossFightUUID:  BossFightUUID,
			MythicplusUUID: MythicplusUUID,
			ColumnUUID:     uuid.Must(uuid.NewV4()).String(),
			Timestamp:      timestamp,
			EventType:      params[0],
		}

		//TODO: never add anything if CombatlogUUID is empty, same logic as m+uuid

		if MythicplusUUID == "" && params[0] != "CHALLENGE_MODE_START" {
			//I don't want to add events if they are outside of a combatlog
			continue
		}

		switch params[0] {
		case "COMBAT_LOG_VERSION":
			//TODO:
			//- [x] check version
			//- [x] check advanced logging
			//- [ ] generate report uuid
			//		- not sure what this is about

			CombatlogUUID = uuid.Must(uuid.NewV4()).String()
			e.CombatlogUUID = CombatlogUUID
			err = e.combatLogVersion(params)
			if err != nil {
				return err
			}
			//NOTE:
			//break is implicit in go, that means after the first match it exits
			//the switch statement

		case "ENCOUNTER_START":
			e.BossFightUUID = uuid.Must(uuid.NewV4()).String()
			err = e.encounterStart(params)
			if err != nil {
				return err
			}

		case "ENCOUNTER_END":
			//I want the entry with encounter_end to have the id, just the records after should be nil again
			//it's been already set for this record (e) so it's okay to clear it befor calling encounterEnd
			BossFightUUID = ""
			err = e.encounterEnd(params)
			if err != nil {
				return err
			}

		case "CHALLENGE_MODE_START":
			MythicplusUUID = uuid.Must(uuid.NewV4()).String()
			e.MythicplusUUID = MythicplusUUID
			err = e.challengeModeStart(params)
			if err != nil {
				return err
			}

		case "CHALLENGE_MODE_END":
			err = e.challengeModeEnd(params)
			if err != nil {
				return err
			}
			combatEvents = append(combatEvents, e)

			r, err := convertToCSV(&combatEvents)
			if err != nil {
				return err
			}
			err = uploadS3(&r, sess, e.MythicplusUUID, csvBucket)
			if err != nil {
				return err
			}

			combatEvents = nil
			MythicplusUUID = ""

		// case "SPELL_HEAL", "SPELL_PERIODIC_HEAL":
		// 	err = e.importHeal(params)
		case "SPELL_DAMAGE":
			err = e.spellDamage(params)
			if err != nil {
				return err
			}

		default:
			e.Unsupported = true
		}

		if params[0] == "CHALLENGE_MODE_END" {
			continue //I'm appending inside the case statement because im uploading the summary after
		}
		combatEvents = append(combatEvents, e)
	}
	return nil
}
