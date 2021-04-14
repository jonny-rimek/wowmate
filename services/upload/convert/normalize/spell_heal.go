package normalize

//TODO:
// func (e *Event) spellHeal(params []string) (err error) {
// 	e.SpellID = params[1]               //Player-970-00307C5B
// 	e.CasterName = trimQuotes(params[2]) //"Brimidreki-Sylvanas"
// 	e.CasterType = params[3]             //0x512
// 	e.SourceFlag = params[4]             //0x0
// 	e.TargetID = params[5]               //Player-970-00307C5B
// 	e.TargetName = trimQuotes(params[6]) //"Brimidreki-Sylvanas"
// 	e.TargetType = params[7]             //0x512
// 	e.DestFlag = params[8]               //0x0
// 	e.SpellID, err = Atoi32(params[9])   //122281
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.SpellName = trimQuotes(params[10])           //"Healing Elixir"
// 	e.SpellType = params[11]                       //0x8
// 	e.AnotherPlayerID = params[12]                 //Player-970-00307C5B
// 	e.D0 = params[13]                              //0000000000000000
// 	e.D1, _ = strconv.ParseInt(params[14], 10, 64) //132358
// 	e.D2, _ = strconv.ParseInt(params[15], 10, 64) //135424
// 	e.D3, _ = strconv.ParseInt(params[16], 10, 64) //4706
// 	e.D4, _ = strconv.ParseInt(params[17], 10, 64) //1467
// 	e.D5, _ = strconv.ParseInt(params[18], 10, 64) //1455
// 	e.D6, _ = strconv.ParseInt(params[19], 10, 64) //3
// 	e.D7, _ = strconv.ParseInt(params[20], 10, 64) //79
// 	e.D8, _ = strconv.ParseInt(params[21], 10, 64) //100
// 	e.D9 = params[22]                              // 0
// 	e.D10 = params[23]                             // -934.51
// 	e.D11 = params[24]                             //2149.50
// 	e.D12 = params[25]                             //3.4243
// 	e.D13 = params[26]                             //307
// 	e.DamageUnkown14 = params[27]                  //.
// 	e.ActualAmount, err = Atoi64(params[28])       //20314
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.Overhealing, err = Atoi64(params[29]) //0
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.Absorbed, err = Atoi64(params[30]) //0
// 	if err != nil {
// 		log.Println("failed to convert heal event")
// 		return err
// 	}
// 	e.Critical = params[31] //nil

// 	return nil
// }
