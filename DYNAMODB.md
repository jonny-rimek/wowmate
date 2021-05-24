# DYNAMODB SCHEMA

this document describes all the access patterns in dynamodb and how they are modelled.

#### ERD:

GitHub has no mermaid support =(

```mermaid
erDiagram
    Mythicplus-Log }o--o{ Player : has
    Player ||..|{ Class : has
    Class ||..|{ Active-Specc : has
    Mythicplus-Log ||--|{ Mythicplus-Dungeon : has
    Mythicplus-Dungeon ||..o{ Level-2-Affix : has
    Mythicplus-Dungeon |o..o{ Level-5-Affix : has
    Mythicplus-Dungeon |o..o{ Level-7-Affix : has
    Mythicplus-Dungeon |o..o{ Level-10-Affix : has
    Level-2-Affix {
      int id
      string name
    }
```

#### access patterns:

1. access m+ log overall damage/overall healing by CombatlogHash
   - PK: LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE
   - SK: LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE
   - healing  
   - PK: LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_HEALING
   - SK: LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_HEALING
   - same pattern for every other type
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
   - PK: DUNGEON#PLAYER#__<player_id>__
   - SK: __<dungeon_id>__#__<time_per_cent>__
   - limit 1 and do for every dungeon
   - do it in parallel in go and fuse result back together
1. get all m+ logs for a dungeon for a player per season
   - PK: DUNGEON#PLAYER#__<player_id>__
   - SK: __<dungeon_id>__#__<time_per_cent>__
   - paginate to get all keys in 1 request 
1. sort player m+ logs by most recent
   - PK: PLAYER_ID#__<player_id>__
   - SK: __<created_at>__
1. get player via player_id
   - PK: PLAYER_ID#__<player_id>__
   - SK: PLAYER_ID#__<player_id>__
1. search by player name
   - PK: PLAYER#SEARCH
   - SK: __<player_name>__
   - add player id as attribute, so the follow-up requests can search by player id and not player name, to avoid rename problems

| Patterns | PK | SK | GSI1PK | GSI1SK | GSI2PK | GSI2SK | GSI3PK | GSI3SK | GSI4PK | GSI4SK |
--- | --- | --- | --- | --- | --- | --- | --- | ---| --- | ---
| 1. | LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE| LOG#KEY#__<combatlog_hash>__#OVERALL_PLAYER_DAMAGE|  |  | - |  |  |  |  |
| 2-4 | LOG#KEY#__season__ | __<key_level>__#__<time_as_percent>__#__<combatlog_hash>__ | LOG#KEY#__<season>__#__<dungeon_id>__ | __<key_level>__#__<time_as_percent>__#__<combatlog_hash>__ | - |  |  |  |  |
| 6. | DEDUP#__<combatlog_hash>__ | DEDUP#__<combatlog_hash>__ |  |  | - |  |  |  |  |
| 7+8 | DUNGEON#PLAYER#__<player_id>__ | DUNGEON#PLAYER#__<player_id>__ |  |  | - |  |  |  |  |
| 9. | PLAYER_ID#__<player_id>__ | __<created_at>__ |  |  | - |  |  |  |  |
| 10. | PLAYER_ID#__<player_id>__ | PLAYER_ID#__<player_id>__ |  |  | - |  |  |  |  |
| 11. | PLAYER | __<player_name>__  |  |  | - |  |  |  |  |

