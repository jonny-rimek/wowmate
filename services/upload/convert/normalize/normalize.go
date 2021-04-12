package normalize

import (
	"bufio"
	"fmt"

	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/gofrs/uuid"
)

// Normalize converts the combatlog to a slice of Event structs
func Normalize(scanner *bufio.Scanner, uploadUUID string) (map[string]map[string][]*timestreamwrite.WriteRecordsInput, error) {
	// var combatEvents []*timestreamwrite.Record
	var combatlogUUID string

	rec := make(map[string]map[string][]*timestreamwrite.WriteRecordsInput)

	if uploadUUID == "" {
		return nil, fmt.Errorf("can't provide an empty uploadUUID")
	}

	// combatEvents2 = make([]Event, 0, 100000) //100.000 is an arbitrary value
	// initialising the slice with a capacity to reduce the amount reallocation
	// the difference in a small log was <1sec -> not worth

	for scanner.Scan() {
		// 4/24 10:42:30.561  COMBAT_LOG_VERSION
		// every line starts with the date followed by the rest separated with 2 spaces.
		// the rest is separated with commas

		// IMPROVE:
		// write version of .Split that accepts a pointer to the string
		// this saved a lot of memory with splitAtComma
		// just look at the implementation of the strings.Split function
		// maybe there is a package that implements string functionality more efficient
		// the main problem is that the strings.Split calls a bunch of other functions
		// that all create a new version of the string and thus bloating the memory
		// UPDATE: it's a lot of work to rewrite everything to use []byte I'll resist the
		// premature optimization for now
		row := splitString(scanner.Text(), "  ")

		// NOTE: not written to DB atm https://github.com/jonny-rimek/wowmate/issues/129
		// gonna use it later again
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

		// don' add events if they are outside of a combatlog
		if combatlogUUID == "" && params[0] != "CHALLENGE_MODE_START" {
			continue
		}

		switch params[0] {
		case "COMBAT_LOG_VERSION":
			// TODO: get patch info from here
			// e.CombatlogUUID = CombatlogUUID
			// err = e.combatLogVersion(params)
			// if err != nil {
			// 	return err
			// }

			// NOTE:
			// break is implicit in go, that means after the first match it exits
			// the switch statement

		case "ENCOUNTER_START":
			// e.BossFightUUID = uuid.Must(uuid.NewV4()).String()
			// err = e.encounterStart(params)
			// if err != nil {
			// 	return err
			// }

		case "ENCOUNTER_END":
			// I want the entry with encounter_end to have the id, just the records after should be nil again
			// it's been already set for this record (e) so it's okay to clear it before calling encounterEnd
			// BossFightUUID = ""
			// err = e.encounterEnd(params)
			// if err != nil {
			// 	return err
			// }

		case "CHALLENGE_MODE_START":
			// COMBAT_LOG_VERSION is not a good place to create this, because after this event, a new
			// COMBAT_LOG_VERSION is always part of the log, but there are a couple of lines in between
			// which would mean those wouldn't have the combatlog_uuid info
			combatlogUUID = uuid.Must(uuid.NewV4()).String()
			rec[combatlogUUID] = make(map[string][]*timestreamwrite.WriteRecordsInput)

			err := challengeModeStart(params, uploadUUID, combatlogUUID, rec)
			if err != nil {
				return rec, err
			}

		case "CHALLENGE_MODE_END":
			err := challengeModeEnd(params, uploadUUID, combatlogUUID, rec)
			if err != nil {
				return rec, err
			}

			combatlogUUID = ""

		// case "SPELL_HEAL", "SPELL_PERIODIC_HEAL":
		// 	err = e.importHeal(params)
		case "SPELL_DAMAGE":
			err := spellDamage(params, &uploadUUID, &combatlogUUID, rec)
			if err != nil {
				return nil, err
			}

		default:
			// e.Unsupported = true
		}

		// if params[0] == "CHALLENGE_MODE_END" {
		// 	continue
		// }
	}

	return rec, nil
}
