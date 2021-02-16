package normalize

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/timestreamwrite"
)

//11/3 09:00:29.792  SPELL_DAMAGE,Player-1302-09C8C064,"Hyrriuk-Archimonde",0x512,0x0,Vehicle-0-3892-1763-30316-122963-00005D638F,"Rezan",0x10a48,0x0,283810,"Reckless Flurry",0x1,Vehicle-0-3892-1763-30316-122963-00005D638F,0000000000000000,3600186,3811638,0,0,2700,1,0,0,0,-790.59,2265.96,935,0.8059,122,1287,1599,-1,1,0,0,0,nil,nil,nil

// v16
// 10/3 05:51:15.415  SPELL_DAMAGE,Player-4184-00130F03,"Unstaebl-Torghast",0x512,0x0,Creature-0-2085-2287-15092-165515-0005F81144,"Depraved Darkblade",0xa48,0x0,127802,"Touch of the Grave",0x20,Creature-0-2085-2287-15092-165515-0005F81144,0000000000000000,92482,96120,0,0,1071,0,3,100,100,0,-2206.68,5071.68,1663,2.1133,60,456,456,-1,32,0,0,0,nil,nil,nil
func spellDamage(params []string) (*timestreamwrite.Record, error) {
	if len(params) != 39 {
		return nil, fmt.Errorf("combatlog version should have 39 columns, it has %v: %v", len(params), params)
	}

	actualAmount, err := Atoi64(params[29]) //1287
	if err != nil {
		log.Printf("failed to convert damage event, field actual amount. got: %v", params[27])
		return nil, err
	}

	// currentTimeInMilliseconds := time.Now().UnixNano() / 1e6
	currentTimeInMilliseconds := time.Now().Unix()

	e := &timestreamwrite.Record{
		Dimensions: []*timestreamwrite.Dimension{
			{
				Name:  aws.String("caster_name"),
				Value: aws.String(trimQuotes(params[2])),
			},
			{
				Name:  aws.String("spell_name"),
				Value: aws.String(trimQuotes(params[10])),
			},
			{
				Name:  aws.String("target_name"),
				Value: aws.String(trimQuotes(params[6])),
			},
		},
		MeasureName:      aws.String("damage"),
		MeasureValue:     aws.String(strconv.FormatInt(actualAmount, 10)),
		MeasureValueType: aws.String("BIGINT"),
		Time:             aws.String(strconv.FormatInt(currentTimeInMilliseconds, 10)), //TODO: get time from log
		TimeUnit:         aws.String("SECONDS"),
	}

	return e, nil
	/*
		if len(params) != 39 {
			return fmt.Errorf("combatlog version should have 39 columns, it has %v: %v", len(params), params)
		}

		e.CasterID = params[1]               //Player-1302-09C8C064 ✔
		e.CasterName = trimQuotes(params[2]) //"Hyrriuk-Archimonde" ✔
		e.CasterType = params[3]             //0x512
		e.SourceFlag = params[4]             //0x0
		e.TargetID = params[5]               //Vehicle-0-3892-1763-30316-122963-00005D638F
		e.TargetName = trimQuotes(params[6]) //"Rezan" ✔
		e.TargetType = params[7]             //0x10a48
		e.DestFlag = params[8]               //0x0
		e.SpellID, err = Atoi32(params[9])   //283810
		if err != nil {
			log.Printf("failed to convert damage event, field spell id. got: %v", params[9])
			return err
		}
		e.SpellName = trimQuotes(params[10]) //"Reckless Flurry" ✔
		e.SpellType = params[11]             //0x1
		// e.AnotherPlayerID = params[12]                   //Vehicle-0-3892-1763-30316-122963-00005D638F
		// e.D0 = params[13]                                //0000000000000000
		// e.D1, err = strconv.ParseInt(params[14], 10, 64) //3600186
		if err != nil {
			log.Printf("failed to convert damage event, field d1. got: %v", params[14])
			return err
		}
		// e.D2, err = strconv.ParseInt(params[15], 10, 64) //3811638
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d2. got: %v", params[15])
		// 	return err
		// }
		// 	.D3, err = strconv.ParseInt(params[16], 10, 64) //0
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d3. got: %v", params[16])
		// 	return err
		// }
		// e.D4, err = strconv.ParseInt(params[17], 10, 64) //0
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d4. got: %v", params[17])
		// 	return err
		// }
		// e.D5, err = strconv.ParseInt(params[18], 10, 64) //2700
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d5. got: %v", params[18])
		// 	return err
		// }
		// e.D6, err = strconv.ParseInt(params[19], 10, 64) //1
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d6. got: %v", params[19])
		// 	return err
		// }
		// e.D7, err = strconv.ParseInt(params[20], 10, 64) //0
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d7. got: %v", params[20])
		// 	return err
		// }
		// e.D8, err = strconv.ParseInt(params[21], 10, 64) //0
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field d8. got: %v", params[21])
		// 	return err
		// }
		// e.D9 = params[22]                        //0
		// e.D10 = params[23]                       //-790.59
		// e.D11 = params[24]                       //2265.96
		// e.D12 = params[25]                       //935 -- mb something like a map id?
		// e.D13 = params[26]                       //0.8059
		// e.DamageUnkown14 = params[27]            //122
		e.ActualAmount, err = Atoi64(params[28]) //1287
		if err != nil {
			log.Printf("failed to convert damage event, field actual amount. got: %v", params[27])
			return err
		}
		// e.BaseAmount, err = Atoi64(params[29]) //1599
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field base amount. got: %v", params[29])
		// 	return err
		// }
		// e.Overkill = params[30]              // ✔ -1 no overkill, otherwise the dmg number it was overkilled with. TODO: convert to int64
		// e.School = params[31]                //1 ✔
		// e.Crushing = params[32]              //0 always 0 with ad10-disci TODO: double check with more data NOT CONFIRMED AS crushing
		// e.Blocked = params[33]               //0 TODO: always a number and should be converted to int64, pretty sure it is not blocked bc it is not reflected by actual_amount vs base_amount like absorbed
		// e.Absorbed, err = Atoi64(params[34]) //0 ✔
		// if err != nil {
		// 	log.Printf("failed to convert damage event, field absorbed. got: %v", params[34])
		// 	return err
		// }
		// e.Critical = params[35]  //nil ✔ fairly certain this one is crit it plays into base and actual amount, nil or 1
		// e.Glancing = params[36]  //nil always nil with ad10-disci TODO double check with more data NOT CONFIRMED AS glancing
		// e.IsOffhand = params[37] //nil always nil with ad10-disci TODO double check with more data NOT CONFIRMED AS is_offhand

		return nil
	*/
}
