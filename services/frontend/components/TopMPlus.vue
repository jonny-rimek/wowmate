<template>
  <div class="min-w-full sm:min-w-0 max-w-7xl mx-auto sm:px-6 lg:px-8 sm:pt-5">
    <div
      class="lg:hidden md:flex md:items-center md:justify-between pt-4 lg:pt-8 lg:pb-6 pl-4 lg:pl-0"
    >
      <div class="flex-1 min-w-0">
        <h2
          class="text-xl font-bold leading-7 text-gray-200 sm:text-2xl sm:leading-9 sm:truncate"
        >
          <!-- TODO: dynamic depending on page -->
          {{ title }}
        </h2>
      </div>
    </div>
    <span v-if='typeof logs.data === "undefined"'>no data availabe, upload a log for that dungeon =)</span>
    <div v-else class="flex flex-col">
      <div class="my-2 py-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div
          class="align-middle inline-block min-w-full shadow overflow-hidden sm:rounded-lg border-b border-gray-500 "
        >
          <table class="min-w-full bg-gray-700">
            <thead class='text-gray-300 bg-gray-700 text-left text-xs leading-4 uppercase tracking-wider'>
              <tr>
                <th
                  class="px-4 md:px-6 py-3 border-b border-gray-500 font-medium "
                >
                  Dungeon
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-b border-gray-500 font-medium "
                >
                  Name
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-b border-gray-500 font-medium "
                >
                  <span class="hidden sm:inline">
                    Overall Damage
                  </span>
                  <span class="sm:hidden">
                    Damage
                  </span>
                </th>
                <th
                  class="px-2 md:px-6 py-3 border-gray-500 border-b font-medium "
                ></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="log in logs.data" :key="log.id">
                <td
                  class="px-4 md:px-6 py-4 whitespace-no-wrap border-b border-gray-500 "
                >
                  <div class="flex items-center">
                    <div class="">
                      <div
                        class="text-sm leading-5 font-medium text-gray-200"
                      >
                        <!--
                        TODO:
                        <span class="sm:hidden">
                          {{ log.dungeonShort }}
                        </span>
                        -->
                        <span class="hidden sm:inline">
                          {{ log.dungeon_name }}
                        </span>
                        <span class="text-gray-400">
                            +{{ log.keylevel }}
                        </span>
                      </div>
                      <div
                        class="text-sm leading-5 text-gray-400"
                      >
                        <div>{{ log.affixes }}</div>
                        <div>{{ log.duration }}</div>
                        <div>{{ log.deaths }} death</div>
                      </div>
                    </div>
                  </div>
                </td>
                <td
                  class="px-2 md:px-6 py-4 whitespace-no-wrap border-b border-gray-500 "
                >
                  <div
                    v-for="player in log.player_damage"
                    :key="player.player_name"
                    class="text-sm capitalize leading-5 text-gray-200"
                  >
                    {{ player.player_name }}
                  </div>
                </td>
                <td
                  class="px-2 md:px-6 py-4 whitespace-no-wrap border-b border-gray-500 "
                >
                  <div
                    v-for="player in log.player_damage"
                    :key="player.player_name"
                    class="block overflow-hidden text-sm text-gray-400"
                  >
                    <!--
                    <span class="hidden sm:inline">
                      {{ player.damageGraph }}
                    </span>
                    -->
                    <span class="text-sm sm:float-right sm:pl-2 leading-5 text-gray-200">{{ player.damage }}</span>
                  </div>
                </td>
                <td
                  class="pr-4 sm:px-4 md:px-6 py-4 whitespace-no-wrap text-right border-b border-gray-500  text-sm leading-5 font-medium"
                >
                  <!-- TODO: add real combatlog uuid -->
                  <nuxt-link
                    :to="{
                      name: 'mythicplus-log-id-damage',
                      params: { id: log.combatlog_uuid }
                    }"
                  class="text-2xl text-red-600 hover:text-red-800"
                  >
                    > <!-- TODO: replace with icon^^ -->
                  </nuxt-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <nav class="px-4 py-3 flex items-center justify-between sm:px-6" aria-label="Pagination">
          <div class="flex-1 flex justify-between sm:justify-end">
            <!-- dont display the previous button if it is the first page -->
            <nuxt-link
              v-if='!logs.first_page'
              :to="{path: this.$route.path, query: { prev: this.$store.state.pagination.prev }}"
              :prefetch="false"
              active-class=""
              class="relative inline-flex items-center px-4 py-2 border text-sm font-medium rounded-md border-gray-300 text-gray-700 bg-gray-100 hover:bg-gray-200"
            >
              Previous
            </nuxt-link>
            <nuxt-link
              v-if=!logs.last_page
              :to="{path: this.$route.path, query: { next: this.$store.state.pagination.next }}"
              :prefetch="false"
              active-class=""
              class="ml-3 relative inline-flex items-center px-4 py-2 border text-sm font-medium rounded-md border-gray-300 text-gray-700 bg-gray-100 hover:bg-gray-200"
            >
              Next
            </nuxt-link>
          </div>
        </nav>
      </div>
    </div>
  </div>
</template>
<script>
export default {
  watchQuery: true,
  props: {
    logs: {
      type: Object,
      required: true
    },
    title: {
      type: String,
      required: true
    },
  }
}
</script>
<style></style>
