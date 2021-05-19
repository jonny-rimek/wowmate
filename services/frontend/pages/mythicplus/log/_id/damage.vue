<template>
  <div class="sm:ml-4 pt-2 sm:p-4 w-full">
    <div class="p-2 pb-4 lg:pt-6 lg:flex lg:items-center lg:justify-between">
      <div class="flex-1 min-w-0">
        <div class="">
          <div class="sm:flex sm:items-baseline leading-10">
            <!-- class="space-y-4 sm:flex sm:items-baseline sm:space-y-0 sm:space-x-10" -->
            <div></div>
            <h3
              class="text-lg lg:text-xl leading-6 font-medium text-gray-200"
            >
              {{ log.dungeon_name }}
              <span
                class="text-base lg:text-lg pr-3 lg:pr-6 text-gray-400"
              >
                {{ log.keylevel }} <!--THIS WILL BE DIFFERENT THAN IN OVERVIEW, CUS STUBBED DATA -->
              </span>
            </h3>
            <div>
              <nav class="-mb-px flex space-x-8">
                <a
                  href="#"
                  class="whitespace-no-wrap pb-1 px-1 border-b-2 font-bold text-sm leading-3 text-gray-300 border-red-800 hover:border-gray-400 focus:outline-none"
                  aria-current="page"
                >
                  Damage
                </a>

                <a
                  href="#"
                  class="line-through whitespace-no-wrap pb-1 px-1 border-b-2 border-transparent font-medium text-sm leading-5 text-gray-300 hover:border-gray-400 focus:outline-none"
                >
                  Healing
                </a>
              </nav>
            </div>
          </div>
        </div>

        <div
          class="flex flex-col sm:mt-0 sm:flex-row sm:flex-wrap text-gray-400"
        >
          <div class="mt-2 flex items-center text-sm leading-5 sm:mr-6">
            <svg
              class="flex-shrink-0 mr-1.5 h-5 w-5 "
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fill-rule="evenodd"
                d="M6 6V5a3 3 0 013-3h2a3 3 0 013 3v1h2a2 2 0 012 2v3.57A22.952 22.952 0 0110 13a22.95 22.95 0 01-8-1.43V8a2 2 0 012-2h2zm2-1a1 1 0 011-1h2a1 1 0 011 1v1H8V5zm1 5a1 1 0 011-1h.01a1 1 0 110 2H10a1 1 0 01-1-1z"
                clip-rule="evenodd"
              />
              <path
                d="M2 13.692V16a2 2 0 002 2h12a2 2 0 002-2v-2.308A24.974 24.974 0 0110 15c-2.796 0-5.487-.46-8-1.308z"
              />
            </svg>
            {{ log.affixes }}
          </div>
          <div class="mt-2 flex items-center text-sm leading-5 sm:mr-6">
            <svg
              class="flex-shrink-0 mr-1.5 h-5 w-5 "
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fill-rule="evenodd"
                d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z"
                clip-rule="evenodd"
              />
            </svg>
            {{ log.duration }}
          </div>
          <div class="mt-2 flex items-center text-sm leading-5 sm:mr-6">
            <svg
              class="flex-shrink-0 mr-1.5 h-5 w-5 "
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                d="M8.433 7.418c.155-.103.346-.196.567-.267v1.698a2.305 2.305 0 01-.567-.267C8.07 8.34 8 8.114 8 8c0-.114.07-.34.433-.582zM11 12.849v-1.698c.22.071.412.164.567.267.364.243.433.468.433.582 0 .114-.07.34-.433.582a2.305 2.305 0 01-.567.267z"
              />
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-13a1 1 0 10-2 0v.092a4.535 4.535 0 00-1.676.662C6.602 6.234 6 7.009 6 8c0 .99.602 1.765 1.324 2.246.48.32 1.054.545 1.676.662v1.941c-.391-.127-.68-.317-.843-.504a1 1 0 10-1.51 1.31c.562.649 1.413 1.076 2.353 1.253V15a1 1 0 102 0v-.092a4.535 4.535 0 001.676-.662C13.398 13.766 14 12.991 14 12c0-.99-.602-1.765-1.324-2.246A4.535 4.535 0 0011 9.092V7.151c.391.127.68.317.843.504a1 1 0 101.511-1.31c-.563-.649-1.413-1.076-2.354-1.253V5z"
                clip-rule="evenodd"
              />
            </svg>
            {{ log.deaths }}
          </div>
        </div>
      </div>
    </div>
    <div class="flex flex-col">
      <div class="overflow-x-auto">
        <div
          class="py-2 lg:py-4 align-middle inline-block min-w-full sm:px-8 lg:px-2"
        >
          <div
            class="shadow overflow-hidden border-b border-gray-500 sm:rounded-lg"
          >
            <table class="bg-gray-700 min-w-full divide-y">
              <thead class="text-gray-300">
                <tr>
                  <th
                    class="px-2 sm:px-6 py-3 text-left text-xs leading-4 font-medium uppercase tracking-wider"
                  >
                    Name
                  </th>
                  <th
                    class="px-2 sm:px-6 py-3 text-left text-xs leading-4 font-medium uppercase tracking-wider"
                  >
                    Damage
                  </th>
                  <th class="px-2 sm:px-6 py-3"></th>
                </tr>
              </thead>
              <tbody class="divide-y border-gray-500 ">
                <tr
                  class="border-gray-500 "
                  v-for="entry in log.player_damage"
                  :key="entry.player_id"
                >
                  <!-- TODO: add dmg breakdown here somewhere -->
                  <td
                    class="px-2 sm:px-6 py-4 whitespace-no-wrap text-sm capitalize leading-5 font-medium text-gray-200 "
                  >
                    {{ entry.player_name }}
                  </td>
                  <td
                    class="sm:px-6 py-4 whitespace-no-wrap text-sm leading-5 text-gray-400"
                  >
                    TODO: graph
                    <span class="text-gray-200 float-right">
                      {{ entry.damage }}</span
                    >
                  </td>
                  <td
                    class="px-2 sm:px-6 pr-4 py-4 whitespace-no-wrap text-right text-sm leading-5 font-medium"
                  >
                    <a href="#" class="text-red-600 hover:text-red-800">></a>
                  </td>
                </tr>

                <!-- More rows... -->
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'mythicplus',
  async asyncData(context) {
    const log = await context.app.$wowmateApi.getMythicplusPlayerOverallDamage()
    return { log }
  },
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
          players: [
            {
              id: 1,
              name: 'terra',
              class: 'paladin',
              specc: 'retribution',
              damage: 56123,
              damageGraph: '||||||||||||||||||||||||||||||||||||||||||||||'
            },
            {
              id: 2,
              name: 'xava',
              class: 'hunter',
              specc: 'beastmaster',
              damageGraph: '|||||||||||||||||||||||||||||||||||||||',
              damage: 46123
            },
            {
              id: 3,
              name: 'micha',
              class: 'monk',
              specc: 'windwalker',
              damageGraph: '|||||||||||||||||||||||||||||||',
              damage: 36123
            },
            {
              id: 4,
              name: 'holytank',
              class: 'paladin',
              specc: 'protection',
              damageGraph: '|||||||||||||||||',
              damage: 26123
            },
            {
              id: 5,
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
