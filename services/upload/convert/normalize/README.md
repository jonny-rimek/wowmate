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
- [ ] fix time conversion https://github.com/jonny-rimek/wowmate/issues/133 and new year edge case
- [ ] 
- [ ] 
- [ ] 
- [ ] upload uuid should be file name without the date stuff
- [x] add remaining unsupported events to readme

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
- [ ] COMBATANT_INFO
- [ ] DAMAGE_SPLIT
- [ ] EMOTE
- [ ] ENVIRONMENTAL_DAMAGE
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
- [ ] SPELL_ABSORBED
- [ ] SPELL_DISPEL
- [ ] SPELL_HEAL_ABSORBED
- [ ] SPELL_INSTAKILL
- [ ] SPELL_INTERRUPT
- [ ] SPELL_MISSED
- [ ] SPELL_PERIODIC_DAMAGE
- [ ] SPELL_PERIODIC_ENERGIZE
- [ ] SPELL_PERIODIC_MISSED
- [ ] SPELL_RESURRECT
- [ ] SPELL_ENERGIZE
- [ ] SWING_DAMAGE
- [ ] SWING_DAMAGE_LANDED
- [ ] SWING_MISSED
- [ ] PARTY_KILL
- [ ] UNIT_DIED
- [ ] ZONE_CHANGE