# Convert lambda:

is taking the WoWCombatlog.txt format of displaying the events and normalizes it into a csv.
It can handle raw .txt file, archived as .zip or compressed as gzip.

The data has to be normalized in order to be able to get loaded into a relational db, where
things like overall damage is calculated.

### TODO:

general:
- [x] check combatlog version
- [x] check advanced combatlog enabled
- [x] fix encounter start/end, ignore commons on split if it is inside ""
- [x] count columns per event and check against expected value
- [x] split up encounter start & end
- [ ] add tests
- [ ] reimplement strings.Split that accepts strings as a pointer
- [ ] 
- [ ] 
- [ ] add remaining unsupported events
- [ ] write tests for normalize
	- [ ] test that no code is added outside of a m+
	- [ ] test that boss fight uuid is empty after the bossfight
	- [ ] test that boss fight uuid is generated after a bossfight starts
	- [ ] combatlog uuid can't be empty
	- [ ] 
	- [ ] 
	- [ ] 

table changes:
- [ ] rename KeyUnkown1 to KeyChests
- [ ] drop advanced combat logging field and column, it has to be 1
- [ ] add column in damage event
- [ ] 
- [ ] 

events: 
- [x] COMBAT_LOG_VERSION
- [x] SPELL_DAMAGE
- [x] CHALLENGE_MODE_END
- [x] CHALLENGE_MODE_START
- [x] ENCOUNTER_END
- [x] ENCOUNTER_START
- [ ] ZONE_CHANGE
- [ ] SPELL_CAST_SUCCESS
- [ ] SPELL_CAST_START
- [ ] SPELL_CAST_FAILED
- [ ] SPELL_AURA_APPLIED
- [ ] SPELL_AURA_REFRESH
- [ ] SPELL_SUMMON
- [ ] SPELL_PERIODIC_HEAL
- [ ] SPELL_AURA_REMOVED
- [ ] SPELL_HEAL
- [ ] SPELL_AURA_APPLIED_DOSE
- [ ] SPELL_CREATE
- [ ] SPELL_AURA_REMOVED_DOSE
- [ ] CHALLENGE_MODE_END
- [ ] CHALLENGE_MODE_START
- [ ] COMBATANT_INFO
- [ ] SPELL_ENERGIZE
- [ ] SWING_DAMAGE
- [ ] SWING_DAMAGE_LANDED
- [ ] 
- [ ] 
- [ ] 