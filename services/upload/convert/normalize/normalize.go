package normalize

import (
	"bufio"
	"fmt"

	"github.com/aws/aws-sdk-go/service/timestreamwrite"
	"github.com/gofrs/uuid"
)

// Normalize converts the combatlog to a slice of Event structs
// TODO: rename, we are not really normalizing the data anymore, because I'm not using a relational database anymore
func Normalize(scanner *bufio.Scanner, uploadUUID string) (map[string]map[string][]*timestreamwrite.WriteRecordsInput, map[string][]string, error) {
	var combatlogUUID string

	rec := make(map[string]map[string][]*timestreamwrite.WriteRecordsInput)
	dedup := make(map[string][]string)

	pets := make(map[string]pet)

	if uploadUUID == "" {
		return nil, nil, fmt.Errorf("can't provide an empty uploadUUID")
	}

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

		params := splitAtCommas(&row[1])

		// don' add events if they are outside of a combatlog
		if combatlogUUID == "" && params[0] != "CHALLENGE_MODE_START" {
			continue
		}

		switch params[0] {
		case "COMBAT_LOG_VERSION":
			err := combatLogVersion(params, uploadUUID, combatlogUUID, rec)
			if err != nil {
				return nil, nil, err
			}

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
			timestamp, err := parseTimestamp(&row[0])
			if err != nil {
				return nil, nil, err
			}

			combatlogUUID = uuid.Must(uuid.NewV4()).String()
			rec[combatlogUUID] = make(map[string][]*timestreamwrite.WriteRecordsInput)

			err = challengeModeStart(params, uploadUUID, combatlogUUID, rec, timestamp)
			if err != nil {
				return rec, nil, err
			}

		case "CHALLENGE_MODE_END":
			err := challengeModeEnd(params, uploadUUID, combatlogUUID, rec)
			if err != nil {
				return nil, nil, err
			}

			combatlogUUID = ""

		case "SPELL_DAMAGE", "SPELL_PERIODIC_DAMAGE", "RANGE_DAMAGE", "SWING_DAMAGE":
			err := damage(params, &uploadUUID, &combatlogUUID, rec, pets)
			if err != nil {
				return nil, nil, err
			}
			dedup[combatlogUUID] = append(dedup[combatlogUUID], params...)

		case "SPELL_SUMMON":
			err := summon(params, &uploadUUID, &combatlogUUID, rec, pets)
			if err != nil {
				return nil, nil, err
			}

		case "COMBATANT_INFO":
			dedup[combatlogUUID] = append(dedup[combatlogUUID], params...)

		default:
		}
	}

	return rec, dedup, nil
}
