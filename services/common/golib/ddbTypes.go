package golib

// Structs for the PlayerDamageDone ddb type

// DynamoDBPlayerDamageDone is used to save player damage done to dynamodb, log specific view
type DynamoDBPlayerDamageDone struct {
	Pk            string             `json:"pk"`
	Sk            string             `json:"sk"`
	Damage        []PlayerDamageDone `json:"player_damage"`
	Duration      string             `json:"duration"`
	Deaths        int                `json:"deaths"`
	Affixes       string             `json:"affixes"`
	Keylevel      int                `json:"keylevel"`
	DungeonName   string             `json:"dungeon_name"`
	DungeonID     int                `json:"dungeon_id"`
	CombatlogUUID string             `json:"combatlog_uuid"`
	Finished      bool               `json:"finished"`
	Intime        int                `json:"intime"`
	Date          int64              `json:"date"`
	CreatedAt     string             `json:"created_at"`
}

// DamagePerSpell is a part of PlayerDamageDone and contains the breakdown of damage per spell
type DamagePerSpell struct {
	SpellID   int    `json:"spell_id"`
	SpellName string `json:"spell_name"`
	Damage    int64  `json:"damage"`
}

// PlayerDamageDone contains player and damage per spell info for the log specific view
type PlayerDamageDone struct {
	Damage         int64            `json:"damage"`
	DamagePerSpell []DamagePerSpell `json:"damage_per_spell"`
	Name           string           `json:"player_name"`
	PlayerID       string           `json:"player_id"`
	Class          string           `json:"class"`
	Specc          string           `json:"specc"`
	/*
	   TODO:
	   	ItemLevel
	   	Covenant
	   	Traits
	   	Conduits
	   	Legendaries
	   	Trinkets
	   	Talents
	*/
}

// Structs for the Keys ddb type

// PlayerDamage contains player and damage info for the top keys view etc.
type PlayerDamage struct {
	Damage   int    `json:"damage"` // TODO: convert to int64
	Name     string `json:"player_name"`
	PlayerID string `json:"player_id"`
	Class    string `json:"class"`
	Specc    string `json:"specc"`
}

// DynamoDBKeys is used to display the top keys and top keys per dungeon
type DynamoDBKeys struct {
	Pk            string         `json:"pk"`
	Sk            string         `json:"sk"`
	Damage        []PlayerDamage `json:"player_damage"`
	Gsi1pk        string         `json:"gsi1pk"`
	Gsi1sk        string         `json:"gsi1sk"`
	Duration      string         `json:"duration"`
	Deaths        int            `json:"deaths"`
	Affixes       string         `json:"affixes"`
	Keylevel      int            `json:"keylevel"`
	DungeonName   string         `json:"dungeon_name"`
	DungeonID     int            `json:"dungeon_id"`
	CombatlogUUID string         `json:"combatlog_uuid"`
	Finished      bool           `json:"finished"`
	Intime        int            `json:"intime"`
	Date          int64          `json:"date"`
	CreatedAt     string         `json:"created_at"`
}

// DynamodbDedup is used to save if a Combatlog has already been uploaded to timestream
// Pk:            fmt.Sprintf("DEDUP#%d", hash),
// Sk:            fmt.Sprintf("DEDUP#%d", hash),
type DynamodbDedup struct {
	Pk            string `json:"pk"`
	Sk            string `json:"sk"`
	CombatlogUUID string `json:"combatlog_uuid"`
	CreatedAt     string `json:"created_at"`
}
