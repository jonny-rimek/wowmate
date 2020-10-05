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
		return fmt.Errorf("combatlog version should have 5 columns, it has %v: %v", len(params), params)
	}

	e.DungeonID, err = Atoi32(params[1])
	if err != nil {
		log.Printf("failed to convert challenge_mode_end field dungeon id: %v", params[1])
		return err
	}

	//NOTE: this is if you timed the key and with how many chests.
	e.KeyUnkown1, err = Atoi32(params[2])
	if err != nil {
		log.Printf("failed to convert challenge_mode_end field key chests: %v", params[2])
		return err
	}

	e.KeyLevel, err = Atoi32(params[3])
	if err != nil {
		log.Printf("failed to convert challenge_mode_end field key level: %v", params[3])
		return err
	}

	e.KeyDuration, err = Atoi64(params[4])
	if err != nil {
		log.Printf("failed to convert challenge_mode_end field key duration: %v", params[4])
		return err
	}

	return nil
}
