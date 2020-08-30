<template>
  <div class="min-w-full sm:min-w-0 max-w-7xl mx-auto sm:px-6 lg:px-8 sm:pt-5">
    <div class="flex flex-col">
      <div class="my-2 py-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div
          class="align-middle inline-block min-w-full shadow overflow-hidden sm:rounded-lg border-b border-gray-200 dark:border-gray-500 "
        >
          <table class="min-w-full">
            <thead>
              <tr>
                <th
                  class="px-4 md:px-6 py-3 border-b border-gray-200  dark:border-gray-500 bg-gray-100 dark:bg-gray-700 text-left text-xs leading-4 font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  Dungeon
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-b border-gray-200 dark:border-gray-500  bg-gray-100 dark:bg-gray-700  text-left text-xs leading-4 font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  Name
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-b border-gray-200 dark:border-gray-500  bg-gray-100 dark:bg-gray-700  text-left text-xs leading-4 font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider"
                >
                  Damage
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-b border-gray-200 dark:border-gray-500  bg-gray-100 dark:bg-gray-700 "
                ></th>
              </tr>
            </thead>
            <tbody class="dark:bg-gray-700">
              <tr>
                <td
                  class="px-4 md:px-6 py-4 whitespace-no-wrap border-b border-gray-200 dark:border-gray-500 "
                >
                  <div class="flex items-center">
                    <div class="">
                      <div
                        class="text-sm leading-5 font-medium text-gray-900 dark:text-gray-200"
                      >
                        <span class="sm:hidden">
                          ToS
                        </span>
                        <span class="hidden sm:inline">
                          Tempel of Sethraliss
                        </span>
                        <span class="text-gray-500">+23*</span>
                      </div>
                      <div
                        class="text-sm leading-5 text-gray-500 dark:text-gray-400"
                      >
                        <div>affixes:</div>
                        <div>34:23 +3:12</div>
                        <div>4 deaths</div>
                      </div>
                    </div>
                  </div>
                </td>
                <td
                  class="px-2 md:px-6 py-4 whitespace-no-wrap border-b border-gray-200 dark:border-gray-500 "
                >
                  <div
                    class="text-sm leading-5 text-gray-900 dark:text-gray-200"
                  >
                    Terra
                  </div>
                  <div
                    class="text-sm leading-5 text-gray-900 dark:text-gray-200"
                  >
                    Xava
                  </div>
                  <div
                    class="text-sm leading-5 text-gray-900 dark:text-gray-200"
                  >
                    Micha
                  </div>
                  <div
                    class="text-sm leading-5 text-gray-900 dark:text-gray-200"
                  >
                    Holytank
                  </div>
                  <div
                    class="text-sm leading-5 text-gray-900 dark:text-gray-200"
                  >
                    Tova
                  </div>
                </td>
                <td
                  class="px-2 md:px-6 py-4 whitespace-no-wrap border-b border-gray-200 dark:border-gray-500 "
                >
                  <div class="block overflow-hidden text-sm dark:text-gray-400">
                    <span class="hidden sm:inline">
                      ||||||||||||||||||||||||||||||||||||||||||||||
                    </span>
                    <span
                      class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-900 dark:text-gray-200"
                      >56.535</span
                    >
                  </div>
                  <div class="block overflow-hidden text-sm dark:text-gray-400">
                    <span class="hidden sm:inline">
                      |||||||||||||||||||||||||||||||||||||||
                    </span>
                    <span
                      class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-900 dark:text-gray-200"
                      >46.535</span
                    >
                  </div>
                  <div class="block overflow-hidden text-sm dark:text-gray-400">
                    <span class="hidden sm:inline">
                      |||||||||||||||||||||||||||||
                    </span>

                    <span
                      class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-900 dark:text-gray-200"
                      >36.535</span
                    >
                  </div>
                  <div class="block overflow-hidden text-sm dark:text-gray-400">
                    <span class="hidden sm:inline">
                      ||||||||||||||||||
                    </span>

                    <span
                      class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-900 dark:text-gray-200"
                      >26.535</span
                    >
                  </div>
                  <div class="block overflow-hidden text-sm dark:text-gray-400">
                    <span class="hidden sm:inline">
                      |||||||||||
                    </span>
                    <span
                      class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-900 dark:text-gray-200"
                      >16.535</span
                    >
                  </div>
                </td>
                <td
                  class="pr-4 sm:px-4 md:px-6 py-4 whitespace-no-wrap text-right border-b border-gray-200 dark:border-gray-500  text-sm leading-5 font-medium"
                >
                  <a href="#" class="text-2xl text-red-600 hover:text-red-800"
                    >></a
                  >
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'mythicplus',
  data() {
    return {
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
          affixes: ['explosive', 'teeming', 'fortified'],
          duration: '34:59 +0:01',
          deaths: 1,
          playerDamage: [
            {
              player: 'terra',
              class: 'paladin',
              specc: 'retribution',
              damage: 56123
            },
            {
              player: 'xava',
              class: 'hunter',
              specc: 'beastmaster',
              damage: 56123
            },
            {
              player: 'micha',
              class: 'monk',
              specc: 'windwalker',
              damage: 56123
            },
            {
              player: 'holytank',
              class: 'paladin',
              specc: 'protection',
              damage: 56123
            },
            {
              player: 'tova',
              class: 'druid',
              specc: 'restoration',
              damage: 56123
            }
          ]
        },
        {
          id: 2,
          dungeonName: 'Freehold',
          affixes: ['explosive', 'teeming', 'fortified'],
          duration: '34:59 +0:01',
          deaths: 1,
          playerDamage: [
            {
              player: 'terra',
              damage: 56123
            },
            {
              player: 'xava',
              damage: 56123
            },
            {
              player: 'micha',
              damage: 56123
            },
            {
              player: 'holytank',
              damage: 56123
            },
            {
              player: 'tova',
              damage: 56123
            }
          ]
        },
        {
          id: 3,
          dungeonName: 'Freehold',
          affixes: ['explosive', 'teeming', 'fortified'],
          duration: '34:59 +0:01',
          deaths: 1,
          playerDamage: [
            {
              player: 'terra',
              damage: 56123
            },
            {
              player: 'xava',
              damage: 56123
            },
            {
              player: 'micha',
              damage: 56123
            },
            {
              player: 'holytank',
              damage: 56123
            },
            {
              player: 'tova',
              damage: 56123
            }
          ]
        },
        {
          id: 4,
          dungeonName: 'Freehold',
          affixes: ['explosive', 'teeming', 'fortified'],
          duration: '34:59 +0:01',
          deaths: 1,
          playerDamage: [
            {
              player: 'terra',
              damage: 56123
            },
            {
              player: 'xava',
              damage: 56123
            },
            {
              player: 'micha',
              damage: 56123
            },
            {
              player: 'holytank',
              damage: 56123
            },
            {
              player: 'tova',
              damage: 56123
            }
          ]
        },
        {
          id: 5,
          dungeonName: 'Freehold',
          affixes: ['explosive', 'teeming', 'fortified'],
          duration: '34:59 +0:01',
          deaths: 1,
          playerDamage: [
            {
              player: 'terra',
              damage: 56123
            },
            {
              player: 'xava',
              damage: 56123
            },
            {
              player: 'micha',
              damage: 56123
            },
            {
              player: 'holytank',
              damage: 56123
            },
            {
              player: 'tova',
              damage: 56123
            }
          ]
        }
      ]
    }
  }
}
</script>

<style></style>
