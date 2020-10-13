package normalize

//Event contains all events that are relevant for the combat itself.
//It does not contain things like COMBAT_LOG_VERSION
//
//example:
//6/30 21:46:57.014  SPELL_HEAL,Player-970-000BF9AC,"Eluxbeta-Sylvanas",0x511,0x0,Player-970-000BF9AC,"Eluxbeta-Sylvanas",0x511,0x0,194509,"Power Word:         Radiance",0x2,Player-970-000BF9AC,0000000000000000,92811,110009,368,4868,801,0,100000,100000,0,-922.47,2149.70,3.0913,326,5101,5101,0,0,nil
type Event struct {
	//START_BASE_PARAMS
	// UploadUUID     string
	// Unsupported    bool //DEBUGGING PARAM
	// CombatlogUUID  string
	// BossFightUUID  string
	MythicplusUUID string
	// ColumnUUID     string
	Timestamp      int64  //6/30 21:46:57.014
	// EventType      string //SPELL_HEAL

	//START COMBAT_LOG_VERSION
	// Version int32 //8
	//Type string [...]                                //ADVANCED_LOG_ENABLED, COLUMN dropped as the value is always the same
	// AdvancedLogEnabled int32 //1 or 0 for on or off
	//END COMBAT_LOG_VERSION

	//START CHALLANGE_MODE_START
	//11/3 09:00:00.760  CHALLENGE_MODE_START,"Atal'Dazar",1763,244,10,[10,11,14,16]
	//11/3 09:34:07.310  CHALLENGE_MODE_END,1763,1,10,2123441
	// DungeonName string //"Atal'dazar"
	// DungeonID   int32  //1763 it's only a guess tho
	// KeyUnkown1  int32  //244, dunno what this is
	// KeyLevel    int32  //10
	// KeyArray    string //[10,11,14,16] no idea....
	// KeyDuration int64  //2123441 my guess it. that this is amount of milliseconds the key took, would be about 35min
	//END CHALLANGE_MODE_START

	//START ENCOUNTER_START
	//11/3 09:00:22.354  ENCOUNTER_START,2086,"Rezan",8,5,1763
	//11/3 09:01:58.364  ENCOUNTER_END,2086,"Rezan",8,5,1
	// EncounterID       int32  //2086
	// EncounterName     string //"Rezan"
	// EncounterUnknown1 int32  //8
	// EncounterUnknown2 int32  //5
	//DungeonID    int32  `parquet:"name=key_level, type=INT32"`    //1763 column already exists is only in encounter start event
	// Killed int32 //1 true 0 false
	//END ENCOUNTER_END

	CasterID   string //Player-970-000BF9AC
	CasterName string //"Eluxbeta-Sylvanas"
	CasterType string //0x511 its always 511 for me and 512 for other grp members and other stuff for enemy trash
	SourceFlag string //0x0
	TargetID   string //Player-970-000BF9AC
	TargetName string //"Eluxbeta-Sylvanas"
	TargetType string //0x511
	DestFlag   string //0x0
	//END_BASE_PARAMS

	//START_BASE_SPELL_PARAMS
	SpellID   int32  //194509
	SpellName string //"Power Word: Radiance"
	SpellType string //0x2 //holy i guess
	//END_BASE_SPELL_PARAMS

	//START_DISPEL_PARAMS
	// ExtraSpellID   int32
	// ExtraSpellName string
	// ExtraSchool    string
	//START_SPELL_AURAS_PARAMS
	// AuraType string //BUFF
	//END_DISPELL_PARAMS
	//END_SPELL_AURAS_PARAMS

	//START_ADVANCED_COMBAT_LOGGING_PARAMS
	// AnotherPlayerID string //Player-970-000BF9AC in case of pets this is always the target player_id not the summoner
	// D0              string //0000000000000000
	// D1              int64  //89449
	// D2              int64  //93932
	// D3              int64  //5637
	// D4              int64  //998
	// D5              int64  //2599
	// D6              int64  //1
	// D7              int64  //596
	// D8              int64  //1000
	// D9              string //0
	// D10             string //-967.46 coordinates?
	// D11             string //2171.14 ^
	// D12             string //0.4005  ^
	// D13             string //313
	// DamageUnkown14  string //Added with combatlog version 9?
	//END_ADVANCED_COMBAT_LOGGING_PARAMS
	//START_HEAL_PARAMS (SPELL_HEAL, SPELL_PERIODIC_HEAL)
	//START_DAMAGE_PARAMS e.g. 3724,5319,-1,1,0,0,0,nil,nil,nil
	ActualAmount int64 //reduced by amor, 2x for crit, reduced by absorb
	// BaseAmount   int64 //before reduction, before crit
	//--
	// Overhealing int64 //0 HEAL events only
	//--
	// Overkill string //0 DMG events only
	// School   string //0 DMG events only - confirmed as spell school
	// Resisted string //0 DMG events only
	// Blocked  string //0 DMG events only
	//--
	// Absorbed int64  //0
	// Critical string //1 = crit, nil = noncrit
	//END_HEAL_PARAMS
	// Glancing  string //0
	// Crushing  string //0
	// IsOffhand string
	//END_DAMAGE_PARAMS
}

//SpellAuraRemoved is a event that show when and which buff is fading on a target
//NO advanced version
//If there is no amount, the field is not even added apparently.
//According to the "docu" amount should only exist if the event has the _DOSE suffix
//
//SpellAuraApplied is simply the reverse of SpellAuraRemoved
//
//example:
//6/30 21:46:25.139  SPELL_AURA_REMOVED,Player-970-000BF9AC,"Eluxbeta-Sylvanas",0x511,0x0,Player-970-000BF9AC,"Eluxbeta-Sylvanas",0x511,0x0,259161,"Speed of Gonk",0x1,BUFF
//amount example:
//6/30 21:46:28.598  SPELL_AURA_REMOVED,Player-970-00326DAB,"Maccounet-Sylvanas",0x518,0x0,Player-970-00326DAB,                 "Maccounet-Sylvanas",0x518,0x0,269279,"Resounding Protection",0x1,BUFF,1376

//------------------------------------------------------------------------------------------
//------------------------------------------------------------------------------------------
//ENVIROMENTAL_DAMAGE
//same as SWING_DAMAGE reduced_amount is here the Environmental Type:
// - Drowning
// - Falling
// - Fatigue
// - Lava
// - Slime
//
//example:
//6/30 21:46:29.856  ENVIRONMENTAL_DAMAGE,0000000000000000,nil,0x80000000,0x80000000,Player-970-000BD9D0,"Mehnari-Anduin",0x518,0x0,Player-970-000BD9D0,0000000000000000,27106,27940,1401,297,544,2,120,120,0,-1092.51,807.05,3.7291,196,Falling,834,834,0,1,0,0,0,nil,nil,nil
//------------------------------------------------------------------------------------------
//------------------------------------------------------------------------------------------

//------------------------------------------------------------------------------------------
//------------------------------------------------------------------------------------------
//SpellCastSuccess is the same as SpellHeal or SpellDamage, but for instant spells and have no SpellCastStart
//this example has no dmg components is it because of the event type or the spell type (utility spell)
//
//example:
//6/30 21:46:32.394  SPELL_CAST_SUCCESS,Player-970-003050DB,"Justbones-Sylvanas",0x518,0x0,0000000000000000,nil,0x80000000,0x80000000,115008,"Chi Torpedo",0x1,Player-970-003050DB,0000000000000000,50940,50940,2198,621,519,3,100,100,0,-1048.48,803.62,3.1439,209

//SPELL_PERIODIC_DAMAGE
//is and event for every dot tick
//
//example:
// 6/30 21:46:38.902  SPELL_PERIODIC_DAMAGE,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,124255,"Stagger",0x1,Player-970-00307C5B,0000000000000000,101295,135424,4706,1467,1455,3,100,100,0,-848.36,2082.47,1.5708,307,552,552,-1,1,0,0,0,nil,nil,nil

//SpellPeriodicHeal hots
//
//example:
//6/30 21:46:40.174  SPELL_PERIODIC_HEAL,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,196608,"Eye of the Tiger",0x8,Player-970-00307C5B,0000000000000000,101687,135424,4706,1467,1455,3,100,100,0,-848.36,2082.47,1.5708,307,392,392,0,0,nil

//UnitDied only gives the target that died, not the one that made the killing blow,
//for that there is a different event, but it doesn't show for everybody
//
//example:
//6/30 21:46:51.080  UNIT_DIED,0000000000000000,nil,0x80000000,0x80000000,Creature-0-4160-1763-15940-135989-0000B7FA28,"Shieldbearer of Zul",0xa48,0x0

//SpellDamage is an event for casted spells (non instant, atleast for casters, dunno if melees always use SPELL_DAMAGE)
//
//example:
//6/30 21:46:54.698  SPELL_DAMAGE,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,Creature-0-4160-1763-15940-122971-000237FA28,"Dazar'ai Juggernaut",0xa48,0x0,121253,"Keg Smash",0x1,Creature-0-4160-1763-15940-122971-000237FA28,0000000000000000,265566,269290,0,0,2700,1,0,0,0,-938.49,2157.20,5.5871,120,3724,5319,-1,1,0,0,0,nil,nil,nil

//SWING_DAMAGE
//and SWING_DAMAGE_LANDED are the same as SPELL_DAMAGE
//
//example:
//6/30 21:46:55.218  SWING_DAMAGE,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,Creature-0-4160-1763-15940-127799-0000B7FA28,"Dazar'ai Honor Guard",0xa48,0x0,Player-970-00307C5B,0000000000000000,112044,135424,4706,1467,1455,3,72,100,0,-930.61,2149.90,3.1377,307,2404,3270,-1,1,0,0,0,nil,nil,nil
//
//SWING_DAMAGE_LANDED
//
//example
//6/30 21:46:55.218  SWING_DAMAGE_LANDED,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,Creature-0-4160-1763-15940-127799-0000B7FA28,"Dazar'ai Honor Guard",0xa48,0x0,Creature-0-4160-1763-15940-127799-0000B7FA28,0000000000000000,226059,234165,0,0,2700,1,0,0,0,-937.74,2149.08,1.2436,120,2404,3270,-1,1,0,0,0,nil,nil,nil

//SPELL_INTERRUPT, SPELL_DISPEL, SPELL_DISPEL_FAILED, SPELL_STOLEN
//
//_DISPEL 	extraSpellId 	extraSpellName 	extraSchool
//dispell und stolen/purge also have the auraType param
//
//example:
//6/30 21:54:37.112  SPELL_INTERRUPT,Player-970-00307C5B,"Brimidreki-Sylvanas",0x512,0x0,Creature-0-4160-1763-15940-128434-0000B7FA28,"Feasting Skyscreamer",0x10a48,0x0,116705,"Spear Hand Strike",0x1,255041,"Terrifying Screech",32
