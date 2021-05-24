# DYNAMODB SCHEMA

this document describes all the access patterns in dynamodb and how they are modelled.

[comment]: <> (TODO: ERD)

#### access patterns:

1. access m+ log over all damage by CombatlogHash
   - PK: LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE
   - SK: LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE
1. sort m+ log by highest key and by time
   - PK: LOG#KEY#__<season>__
   - SK: __<key_level>__#__<time_as_percent>__#__<combatlog_hash>__
1. sort m+ log by highest key per dungeon and sorted by time
   - GSI1PK: LOG#KEY#__<season>__#__<dungeon_id>__
   - GSI1SK: __<key_level>__#__<time_as_percent>__#__<combatlog_hash>__
1. sort m+ log by highest key per dungeon per affix and sorted by time
   - same, but filter with filter expression. because it is a rare pattern
1. sort m+ log by highest key per dungeon per specc/class contained and sorted by time
   - same, but filter with filter expression. because it is a rare pattern
1. check for duplicate combatlog by combatlog hash
   - PK: DEDUP#__<combatlog_hash>__
   - SK: DEDUP#__<combatlog_hash>__
1. get best m+ log from each dungeon for a player
1. sort player m+ logs by most recent
1. get all m+ logs for a dungeon for a player per season
1. get player via player_id
1. search by player name

| Patterns | PK | SK | GSI1PK | GSI1SK | GSI2PK | GSI2SK | GSI3PK | GSI3SK | GSI4PK | GSI4SK |
--- | --- | --- | --- | --- | --- | --- | --- | ---| --- | ---
| 1. | LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE| LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE|  |  | - |  |  |  |  |
| 2-4 | LOG#KEY#__season__ | __<key_level>__#__<time_as_percent>__#__<combatlog_hash>__ | LOG#KEY#__<season>__#__<dungeon_id>__ | __<key_level>__#__<time_as_percent>__#__<combatlog_hash>__ | - |  |  |  |  |
| 4. | DEDUP#__<combatlog_hash>__ | DEDUP#__<combatlog_hash>__ |  |  | - |  |  |  |  |
| 5. |  |  |  |  | - |  |  |  |  |
| 6. |  |  |  |  | - |  |  |  |  |
| 7. |  |  |  |  | - |  |  |  |  |
| 8. |  |  |  |  | - |  |  |  |  |
| 9. |  |  |  |  | - |  |  |  |  |    
