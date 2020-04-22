<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 lg:pt-5">
    <div class="flex flex-col">
      <div class="-my-2 py-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div
          class="align-middle inline-block min-w-full shadow overflow-hidden sm:rounded-lg border-b border-gray-200"
        >
          <!-- TODO: column ilvl dungeon level bar for damage -->
          <table class="min-w-full">
            <thead>
              <tr>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                >
                  Name
                </th>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                >
                  M+ stufe
                </th>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                ></th>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                >
                  Damage
                </th>
              </tr>
            </thead>
            <tbody>
              <tr class="bg-white">
                <td
                  class="px-6 py-4 whitespace-no-wrap text-sm leading-5 font-medium text-gray-900"
                >
                  Reinhart
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap text-sm leading-5 text-gray-500"
                >
                  22
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
    <!-- </div> -->
    <!-- Content goes here -->
  </div>
</template>

<script>
export default {
  layout: 'mythicplus'
}
</script>

<style></style>
