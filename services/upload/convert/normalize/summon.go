package normalize

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

type pet struct {
	OwnerID   string
	Name      string
	OwnerName string
	SpellID   int
}

// 1/24 16:50:53.696  SPELL_SUMMON,Player-3674-0906D09A,"Bihla-TwistingNether",0x512,0x0,Creature-0-4234-2291-11942-19668-00000DA56B,"Shadowfiend",0xa28,0x0,34433,"Shadowfiend",0x20
func summon(params []string, uploadUUID *string, combatlogUUID *string, rec map[string]map[string][]*timestreamwrite.WriteRecordsInput, pets map[string]pet) error {
	// only track player pets
	casterType := trimQuotes(params[3])
	if casterType != "0x512" && casterType != "0x511" {
		return nil
	}

	if len(params) != 12 {
		return fmt.Errorf("SPELL_SUMMON should have 13 columns, it has %v: %v", len(params), params)
	}
	spellID, err := strconv.Atoi(params[9]) // 34433
	if err != nil {
		return fmt.Errorf("failed to convert pet summon event, field spell id. got: %v", params[10])
	}

	pets[params[5]] = pet{ // Creature-0-4234-2291-11942-19668-00000DA56B
		OwnerID:   params[1],
		OwnerName: trimQuotes(params[2]),
		Name:      trimQuotes(params[10]), // params[7] is the pet name and [11] is the spell name, in this case the same
		SpellID:   spellID,
	}

	return nil
}
