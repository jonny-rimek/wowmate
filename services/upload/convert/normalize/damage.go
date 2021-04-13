package normalize

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

// v17
// 1/24 16:47:51.662  SPELL_DAMAGE,Player-581-04A01EDA,"Ayléén-Blackrock",0x511,0x0,Creature-0-4234-2291-11942-168992-00000DA4A7,
// "Risen Cultist",0x10a48,0x0,50842,"Blood Boil", 0x20,Creature-0-4234-2291-11942-168992-00000DA4A7,0000000000000000,177038,177822,
// 0,0,1071,0,0,2434,2434,0,2943.32,-2219.79,1680,3.6342,60,784,783,-1,32,0,0,0,nil,nil,nil
//
//
// 1/24 16:47:51.662  SPELL_DAMAGE,Player-581-04A01EDA,"Ayléén-Blackrock",0x511,0x0,Creature-0-4234-2291-11942-168992-00000DA4A7,"Risen Cultist",0x10a48,0x0,50842,"Blood Boil", 0x20,Creature-0-4234-2291-11942-168992-00000DA4A7,0000000000000000,177038,177822,0,0,1071,0,0,2434,2434,0,2943.32,-2219.79,1680,3.6342,60,784,783,-1,32,0,0,0,nil,nil,nil
// 1/24 16:48:17.916  SPELL_PERIODIC_DAMAGE,Player-3674-0906D09A,"Bihla-TwistingNether",0x512,0x0,Creature-0-4234-2291-11942-168992-00020DA4A7,"Risen Cultist",0xa48,0x0,204213,"Purge the Wicked",0x4,Creature-0-4234-2291-11942-168992-00020DA4A7,0000000000000000,748,177822,0,0,1071,0,0,2434,2434,0,2903.19,-2224.57,1680,4.6220,60,483,241,-1,4,0,0,0,1,nil,nil
// 1/24 16:48:22.569  RANGE_DAMAGE,Player-1403-09B9285B,"Luminal-Draenor",0x512,0x0,Creature-0-4234-2291-11942-174773-00008DA4CF,"Spiteful Shade",0xa48,0x0,75,"Auto Shot",0x1,Creature-0-4234-2291-11942-174773-00008DA4CF,Creature-0-4234-2291-11942-168949-00010DA4A7,40914,88911,0,0,1071,0,1,0,0,0,2904.33,-2216.53,1680,1.2998,60,1295,924,-1,1,0,0,0,1,nil,nil
// 1/24 16:48:22.684  SWING_DAMAGE,Player-581-04A01EDA,"Ayléén-Blackrock",0x511,0x0,Creature-0-4234-2291-11942-169905-00010DA4A7,"Risen Warlord",0x10a48,0x0,Player-581-04A01EDA, 0000000000000000,45866,56060,1665,242,1918,0,6,543,1250,0,2901.84,-2241.36,1680,1.1696,204,791,1129,-1,1,0,0,0,nil,nil,nil
// pet event, note caster_type = 0x2112
// 1/24 16:50:53.696  SPELL_SUMMON,Player-3674-0906D09A,"Bihla-TwistingNether",0x512,0x0,Creature-0-4234-2291-11942-19668-00000DA56B,"Shadowfiend",0xa28,0x0,34433,"Shadowfiend",0x20
// SWING_DAMAGE_LANDED,Creature-0-4234-2291-11942-19668-00000DA56B,"Shadowfiend",0x2112,0x0,Creature-0-4234-2291-11942-168934-00000DA4A7,"Enraged Spirit",0x10a48,0x0,Creature-0-4234-2291-11942-168934-00000DA4A7,0000000000000000,84989,320080,0,0,1071,0,1,0,0,0,2605.79,-2121.73,1680,6.2740,60,788,788,-1,32,0,0,0,nil,nil,nil
// 1/24 16:55:26.219  SWING_DAMAGE,Pet-0-4234-2291-11942-58964-02021134B7,"Kerkek",0x1112,0x0,Creature-0-4234-2291-11942-170572-00008DA4A7,"Atal'ai Hoodoo Hexxer",0x10a48,0x1,Pet-0-4234-2291-11942-58964-02021134B7,Player-3674-0877EE37,19795,19795,849,1697,951,220,3,200,200,0,2743.39,-1819.85,1679,0.4638,201,210,299,-1,1,0,0,0,nil,nil,nil
// 1/24 16:55:26.505  SPELL_DAMAGE,Pet-0-4234-2291-11942-58964-02021134B7,"Kerkek",0x1112,0x0,Creature-0-4234-2291-11942-170572-00008DA4A7,"Atal'ai Hoodoo Hexxer",0x10a48,0x1,54049,"Shadow Bite",0x20,Creature-0-4234-2291-11942-170572-00008DA4A7,0000000000000000,165445,284515,0,0,1071,0,0,2434,2434,0,2745.74,-1818.67,1679,5.5495,60,534,534,-1,32,0,0,0,nil,nil,nil
//
// next summon, note new caster id. didn't find an event that the pet is gone (for shadowfiend)
// 1/24 16:57:32.383  SPELL_SUMMON,Player-3674-0906D09A,"Bihla-TwistingNether",0x512,0x0,Creature-0-4234-2291-11942-19668-00000DA6FA,"Shadowfiend",0xa28,0x0,34433,"Shadowfiend",0x20
// passing in the uuids as a pointer to the string reduced the mb usage by a couple MB
func damage(params []string, uploadUUID *string, combatlogUUID *string, rec map[string]map[string][]*timestreamwrite.WriteRecordsInput, pets map[string]pet) error {

	var actualAmount, currentTimeInSeconds int64
	var spellID int
	var err error
	var casterID, casterName, casterType, spellName string

	// TODO: refactor and test
	pet, exists := pets[params[1]]
	if exists == true {
		// Pet damage:
		// instead of getting the exact spell id and having a break down by spell per pet etc, I group everything
		// under the spell id that was used to summon the pet, no matter if the actual damage was an auto attack by the pet
		// or a cast
		spellID = pet.SpellID
		spellName = pet.Name

		casterName = pet.OwnerName
		casterType = "0x512"
		casterID = pet.OwnerID

		if params[0] != "SWING_DAMAGE" {
			actualAmount, err = Atoi64(params[29]) // 783
			if err != nil {
				log.Printf("failed to convert damage event, field actual amount. got: %v", params[29])
				return err
			}
		} else {
			actualAmount, err = Atoi64(params[26])
			if err != nil {
				log.Printf("failed to convert damage event, field actual amount. got: %v", params[26])
				return err
			}
		}
	} else {
		if params[0] != "SWING_DAMAGE" {
			// SWING_DAMAGE is 3 elements shorter because it doesn't have fields for spellID and spellName
			// everything else is the same
			if len(params) != 39 {
				return fmt.Errorf("*_DAMAGE should have 39 columns, it has %v: %v", len(params), params)
			}

			actualAmount, err = Atoi64(params[29]) // 783
			if err != nil {
				log.Printf("failed to convert damage event, field actual amount. got: %v", params[29])
				return err
			}

			spellID, err = strconv.Atoi(params[9]) // 50842
			if err != nil {
				log.Printf("failed to convert damage event, field spell id. got: %v", params[9])
				return err
			}
			spellName = trimQuotes(params[10])
		} else {
			if len(params) != 36 {
				return fmt.Errorf("SWING_DAMAGE should have 36 columns, it has %v: %v", len(params), params)
			}

			actualAmount, err = Atoi64(params[26])
			if err != nil {
				log.Printf("failed to convert damage event, field actual amount. got: %v", params[26])
				return err
			}

			spellID = 42013370310 // just a random number
			// TODO: use different name if a pet auto attacks
			spellName = "Auto Attack"
		}

		casterID = params[1]
		casterName = trimQuotes(params[2])
		casterType = trimQuotes(params[3])
	}

	currentTimeInSeconds = time.Now().Unix()
	rand.Seed(time.Now().UnixNano())

	key := strconv.Itoa(spellID)

	_, exists = rec[*combatlogUUID][key]
	if exists == true {
		var tmp []*timestreamwrite.WriteRecordsInput
		tmp = make([]*timestreamwrite.WriteRecordsInput, len(rec[*combatlogUUID][key]))
		copy(tmp, rec[*combatlogUUID][key])

		// I only care about the last element, because all other are already at 100 records
		last := len(tmp) - 1

		if len(tmp[last].Records) < 100 {
			rec[*combatlogUUID][key][last].Records = append(rec[*combatlogUUID][key][last].Records, &timestreamwrite.Record{
				// TODO: extra dimensions and common attribute to vars and reuse
				Dimensions: []*timestreamwrite.Dimension{
					{
						Name:  aws.String("caster_id"),
						Value: aws.String(casterID),
					},
					{
						Name:  aws.String("caster_name"), //
						Value: aws.String(casterName),
					},
					{
						Name:  aws.String("caster_type"),
						Value: aws.String(casterType),
					},
					{
						Name:  aws.String("rnd"),
						Value: aws.String(strconv.Itoa(rand.Int())),
						// replace with time from log
					},
				},
				MeasureValue: aws.String(strconv.FormatInt(actualAmount, 10)),
			})
		} else {
			writeInput := &timestreamwrite.WriteRecordsInput{
				CommonAttributes: &timestreamwrite.Record{
					Dimensions: []*timestreamwrite.Dimension{
						{
							Name:  aws.String("spell_id"),
							Value: aws.String(strconv.Itoa(spellID)),
						},
						{
							Name:  aws.String("spell_name"),
							Value: aws.String(spellName),
						},
						{
							Name:  aws.String("upload_uuid"),
							Value: aws.String(*uploadUUID),
						},
						{
							Name:  aws.String("combatlog_uuid"),
							Value: aws.String(*combatlogUUID),
						},
					},
					MeasureName:      aws.String("damage"),
					MeasureValueType: aws.String("BIGINT"),
					TimeUnit:         aws.String("SECONDS"), // can specify seconds for timestream instead of ms!
					Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
					// I don't care about this time, it just the time we create this entry, not the time of the combatlog event
					// I also don't care about the exact time this is written, so I always use the time the first record is created
					// and reuse it for the subsequent ones
				},
				Records: []*timestreamwrite.Record{
					{
						Dimensions: []*timestreamwrite.Dimension{
							{
								Name:  aws.String("caster_id"),
								Value: aws.String(casterID),
							},
							{
								Name:  aws.String("caster_name"), //
								Value: aws.String(casterName),
							},
							{
								Name:  aws.String("caster_type"),
								Value: aws.String(casterType),
							},
							{
								Name:  aws.String("rnd"),
								Value: aws.String(strconv.Itoa(rand.Int())),
								// replace with time from log
							},
						},
						MeasureValue: aws.String(strconv.FormatInt(actualAmount, 10)),
					},
				},
			}

			rec[*combatlogUUID][key] = append(rec[*combatlogUUID][key], writeInput)
		}
		return nil
	}

	writeRecordsInputs := []*timestreamwrite.WriteRecordsInput{
		{
			CommonAttributes: &timestreamwrite.Record{
				Dimensions: []*timestreamwrite.Dimension{
					{
						Name:  aws.String("spell_id"),
						Value: aws.String(strconv.Itoa(spellID)),
					},
					{
						Name:  aws.String("spell_name"),
						Value: aws.String(spellName),
					},
					{
						Name:  aws.String("upload_uuid"),
						Value: aws.String(*uploadUUID),
					},
					{
						Name:  aws.String("combatlog_uuid"),
						Value: aws.String(*combatlogUUID),
					},
				},
				MeasureName:      aws.String("damage"),
				MeasureValueType: aws.String("BIGINT"),
				TimeUnit:         aws.String("SECONDS"), // can specify seconds as rec for timestream instead of ms!
				Time:             aws.String(strconv.FormatInt(currentTimeInSeconds, 10)),
				// I don't care about this time, it just the time we create this entry, not the time of the combatlog event
				// I also don't care about the exact time this is written, so I always use the time the first record is created
				// and reuse it for the subsequent ones
			},
			Records: []*timestreamwrite.Record{
				{
					Dimensions: []*timestreamwrite.Dimension{
						{
							Name:  aws.String("caster_id"),
							Value: aws.String(casterID),
						},
						{
							Name:  aws.String("caster_name"),
							Value: aws.String(casterName),
						},
						{
							Name:  aws.String("caster_type"),
							Value: aws.String(casterType),
						},
						{
							Name:  aws.String("rnd"),
							Value: aws.String(strconv.Itoa(rand.Int())),
							// replace with time from log
						},
					},
					MeasureValue: aws.String(strconv.FormatInt(actualAmount, 10)),
				},
			},
		},
	}
	rec[*combatlogUUID][key] = writeRecordsInputs

	return nil
}
