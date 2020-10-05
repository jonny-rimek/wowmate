package normalize

import (
	"fmt"
	"log"
)

//11/3 09:00:22.354  ENCOUNTER_START,2086,"Rezan",8,5,1763

//16
//10/3 05:59:07.379  ENCOUNTER_START,2401,"Halkias, the Sin-Stained Goliath",8,5,2287
func (e *Event) encounterStart(params []string) (err error) {
	if len(params) != 6 {
		return fmt.Errorf("combatlog version should have 6 columns, it has %v: %v", len(params), params)
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
