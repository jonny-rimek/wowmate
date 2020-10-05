package normalize

import (
	"fmt"
	"log"
)

//11/3 09:34:07.310  CHALLENGE_MODE_END,1763,1,10,2123441

//v16
// 10/3 05:51:00.879  CHALLENGE_MODE_END,2287,0,0,0 //this one was just after entering the dungeon
// 10/3 06:14:35.797  CHALLENGE_MODE_END,2287,1,2,1451286 // this was after finishing the key
func (e *Event) challengeModeEnd(params []string) (err error) {
	if len(params) != 5 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.DungeonID, err = Atoi32(params[1])
	if err != nil {
		log.Println("failed to convert challange mode end event")
		return err
	}
	//NOTE: this is if you timed the key and with how many chests.
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
