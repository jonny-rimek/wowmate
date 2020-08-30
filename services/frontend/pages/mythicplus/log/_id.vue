<template>
  <div>
    view for a specific log
  </div>
</template>

<script>
export default {
  layout: 'mythicplus',
  data() {
    return {
      logs: [
        {
          id: 1,
          dungeonName: 'Freehold',
          dungeonShort: 'FH',
          keyLevel: '+23*',
          affixes: 'explosive, teeming, fortified',
          duration: '34:59 +0:01',
          deaths: 1,
          player: [
            {
              playerId: 1,
              name: 'terra',
              class: 'paladin',
              specc: 'retribution',
              damage: 56123,
              damageGraph: '||||||||||||||||||||||||||||||||||||||||||||||'
            },
            {
              playerId: 2,
              name: 'xava',
              class: 'hunter',
              specc: 'beastmaster',
              damageGraph: '|||||||||||||||||||||||||||||||||||||||',
              damage: 46123
            },
            {
              playerId: 3,
              name: 'micha',
              class: 'monk',
              specc: 'windwalker',
              damageGraph: '|||||||||||||||||||||||||||||||',
              damage: 36123
            },
            {
              playerId: 4,
              name: 'holytank',
              class: 'paladin',
              specc: 'protection',
              damageGraph: '|||||||||||||||||',
              damage: 26123
            },
            {
              playerId: 5,
              name: 'tova',
              class: 'druid',
              specc: 'restoration',
              damageGraph: '|||||||',
              damage: 16123
            }
          ]
        }
      ]
    }
  }
}
</script>

<style></style>
