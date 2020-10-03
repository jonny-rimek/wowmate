# Convert lambda:

is taking the WoWCombatlog.txt format of displaying the events and normalizes it into a csv.
It can handle raw .txt file, archived as .zip or compressed as gzip.

The data has to be normalized in order to be able to get loaded into a relational db, where
things like overall damage is calculated.

### TODO:

- [x] check combatlog version
- [x] check advanced combatlog enabled
- [x] fix encounter start/end, ignore commons on split if it is inside ""
- [ ] count columns per event and check against expected value
- [ ] split up encounter start & end
- [ ] check expected content of each cell with regex

events: 
- [ ] COMBAT_LOG_VERSION
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
- [ ] SPELL_DAMAGE
- [ ] SWING_DAMAGE_LANDED
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 
- [ ] 