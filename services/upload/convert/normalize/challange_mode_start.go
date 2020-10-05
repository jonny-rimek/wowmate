package normalize

import (
	"fmt"
	"log"
)

//11/3 09:00:00.760  CHALLENGE_MODE_START,"Atal'Dazar",1763,244,10,[10,11,14,16]

//v16
//10/3 05:51:00.975  CHALLENGE_MODE_START,"Halls of Atonement",2287,378,2,[10]
//NOTE: the array is definitely the affixes
func (e *Event) challengeModeStart(params []string) (err error) {
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
