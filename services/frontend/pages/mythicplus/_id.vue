<template>
  <TopMPlus :logs="logs" :title="title"></TopMPlus>
</template>

<script>
import TopMPlus from '~/components/TopMPlus.vue'

export default {
  components: {
    TopMPlus
  },
  layout: 'mythicplus',
  data() {
    return {
      title: 'Not All dungeons',
      dungeons: [
        { name: 'All dungeons', pathName: 'mythicplus' },
        { name: "Atal'Dazar", pathName: 'mythicplus-id', id: 2144 },
        { name: 'Freehold', pathName: 'mythicplus-id', id: 2145 },
        { name: "Kings'Rest", pathName: 'mythicplus-id', id: 2146 },
        { name: 'Shrine of the Storm', pathName: 'mythicplus-id', id: 2147 },
        { name: 'Siege of Boralus', pathName: 'mythicplus-id', id: 2148 },
        { name: 'Temple of Sethraliss', pathName: 'mythicplus-id', id: 2149 },
        { name: 'The MOTHERLODE!!', pathName: 'mythicplus-id', id: 2154 },
        { name: 'The Underrot', pathName: 'mythicplus-id', id: 2164 },
        { name: 'Tol Dagor', pathName: 'mythicplus-id', id: 2174 },
        { name: 'Waycrest Manor', pathName: 'mythicplus-id', id: 2184 },
        {
          name: 'Operation: Mechagon - Junkyard',
          pathName: 'mythicplus-id',
          id: 2194
        },
        {
          name: 'Operation: Mechagon - Workshop',
          pathName: 'mythicplus-id',
          id: 2104
        }
      ],
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
        },
        {
          id: 2,
          dungeonName: 'Template of Sethralis',
          dungeonShort: 'ToS',
          keyLevel: '+24*',
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
