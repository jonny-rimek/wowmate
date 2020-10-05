package normalize

import (
	"fmt"
	"log"
)

//v16
//10/3 05:44:30.076  COMBAT_LOG_VERSION,16,ADVANCED_LOG_ENABLED,1,BUILD_VERSION,9.0.2,PROJECT_ID,1
func (e *Event) combatLogVersion(params []string) (err error) {
	if len(params) != 8 {
		return fmt.Errorf("combatlog version should have 8 columns, it has %v: %v", len(params), params)
	}

	e.Version, err = Atoi32(params[1])
	if err != nil {
		log.Println("failed to convert combatlog version event")
		return err
	}
	if e.Version != 16 {
		return fmt.Errorf("unsupported combatlog version: %v, only version 16 is supported", e.Version)
	}
	e.AdvancedLogEnabled, err = Atoi32(params[3])
	if err != nil {
		log.Println("failed to convert combatlog version event")
		return err
	}
	if e.AdvancedLogEnabled != 1 {
		return fmt.Errorf("advanced combatlogging must be enabled")
	}
	return nil
}