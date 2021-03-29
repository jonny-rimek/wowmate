package normalize

//11/3 09:01:58.364  ENCOUNTER_END,2086,"Rezan",8,5,1

//v16
//10/3 06:00:02.433  ENCOUNTER_END,2401,"Halkias, the Sin-Stained Goliath",8,5,1
// func (e *Event) encounterEnd(params []string) (err error) {
// 	if len(params) != 6 {
// 		return fmt.Errorf("combatlog version should have 6 columns, it has %v: %v", len(params), params)
// 	}

// 	e.EncounterID, err = Atoi32(params[1])
// 	if err != nil {
// 		log.Println("failed to convert encounter end event")
// 		return err
// 	}

// 	e.EncounterName = trimQuotes(params[2])

// 	e.EncounterUnknown1, err = Atoi32(params[3])
// 	if err != nil {
// 		log.Println("failed to convert encounter end event")
// 		return err
// 	}

// 	e.EncounterUnknown2, err = Atoi32(params[4])
// 	if err != nil {
// 		log.Println("failed to convert encounter end event")
// 		return err
// 	}

// 	e.Killed, err = Atoi32(params[5])
// 	if err != nil {
// 		log.Println("failed to convert encounter end event")
// 		return err
// 	}
// 	return nil
// }
