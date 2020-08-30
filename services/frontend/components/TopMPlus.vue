<template>
  <div class="min-w-full sm:min-w-0 max-w-7xl mx-auto sm:px-6 lg:px-8 sm:pt-5">
    <div
      class="lg:hidden md:flex md:items-center md:justify-between pt-4 lg:pt-8 lg:pb-6 pl-4 lg:pl-0"
    >
      <div class="flex-1 min-w-0">
        <h2
          class="text-l font-bold leading-7 text-gray-900 dark:text-gray-200 sm:text-2xl sm:leading-9 sm:truncate"
        >
          <!-- TODO: dynamic depending on page -->
          {{ title }}
        </h2>
      </div>
    </div>
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
                  <span class="hidden sm:inline">
                    Overall Damage
                  </span>
                  <span class="sm:hidden">
                    Damage
                  </span>
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-b border-gray-200 dark:border-gray-500  bg-gray-100 dark:bg-gray-700 "
                ></th>
              </tr>
            </thead>
            <tbody class="dark:bg-gray-700">
              <tr v-for="log in logs" :key="log.id">
                <td
                  class="px-4 md:px-6 py-4 whitespace-no-wrap border-b border-gray-200 dark:border-gray-500 "
                >
                  <div class="flex items-center">
                    <div class="">
                      <div
                        class="text-sm leading-5 font-medium text-gray-900 dark:text-gray-200"
                      >
                        <span class="sm:hidden">
                          {{ log.dungeonShort }}
                        </span>
                        <span class="hidden sm:inline">
                          {{ log.dungeonName }}
                        </span>
                        <span class="text-gray-500 dark:text-gray-400">{{
                          log.keyLevel
                        }}</span>
                      </div>
                      <div
                        class="text-sm leading-5 text-gray-500 dark:text-gray-400"
                      >
                        <div>{{ log.affixes }}</div>
                        <div>{{ log.duration }}</div>
                        <div>{{ log.deaths }} death</div>
                      </div>
                    </div>
                  </div>
                </td>
                <td
                  class="px-2 md:px-6 py-4 whitespace-no-wrap border-b border-gray-200 dark:border-gray-500 "
                >
                  <div
                    v-for="player in log.player"
                    :key="player.id"
                    class="text-sm capitalize leading-5 text-gray-900 dark:text-gray-200"
                  >
                    {{ player.name }}
                  </div>
                </td>
                <td
                  class="px-2 md:px-6 py-4 whitespace-no-wrap border-b border-gray-200 dark:border-gray-500 "
                >
                  <div
                    v-for="player in log.player"
                    :key="player.id"
                    class="block overflow-hidden text-sm dark:text-gray-400"
                  >
                    <span class="hidden sm:inline">
                      {{ player.damageGraph }}
                    </span>
                    <span
                      class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-900 dark:text-gray-200"
                      >{{ player.damage }}</span
                    >
                  </div>
                </td>
                <td
                  class="pr-4 sm:px-4 md:px-6 py-4 whitespace-no-wrap text-right border-b border-gray-200 dark:border-gray-500  text-sm leading-5 font-medium"
                >
                  <nuxt-link
                    :to="{
                      name: 'mythicplus-log-id',
                      params: { id: log.id }
                    }"
                    class="text-2xl text-red-600 hover:text-red-800"
                  >
                    >
                  </nuxt-link>
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
  props: {
    logs: {
      type: Array,
      required: true
    },
    title: {
      type: String,
      required: true
    }
  }
}
</script>
<style></style>
