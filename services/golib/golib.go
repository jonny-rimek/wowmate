package golib

//DamageSummary is the DynamoDB schema for all damage summaries
type DamageSummary struct {
	PK            string `json:"pk"`
	Damage        int64  `json:"sk"`
	EncounterID   int    `json:"gsi1pk"`
	BossFightUUID string `json:"gsi2pk"`
	CasterID      string `json:"gsi3pk"`
	CasterName    string `json:"caster_name"`
}