<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-5">
    <div class="flex flex-col">
      <div class="-my-2 py-2 overflow-x-auto sm:-mx-6 sm:px-6 lg:-mx-8 lg:px-8">
        <div
          class="align-middle inline-block min-w-full shadow overflow-hidden sm:rounded-lg border-b border-gray-200"
        >
          <table class="min-w-full">
            <thead>
              <tr>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                >
                  Dungeon
                </th>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                >
                  Name
                </th>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                ></th>
                <th
                  class="px-6 py-3 border-b border-gray-200 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider"
                >
                  Damage
                </th>
                <th class="px-6 py-3 border-b border-gray-200 bg-gray-50"></th>
              </tr>
            </thead>
            <tbody class="bg-white">
              <tr>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="flex items-center">
                    <div class="ml-4">
                      <div class="text-sm leading-5 font-medium text-gray-900">
                        Tempel of Sethraliss
                        <span class="text-gray-500">+23*</span>
                      </div>
                      <div class="text-sm leading-5 text-gray-500">
                        <div>affixes:</div>
                        <div>34:23 +3:12</div>
                        <div>4 deaths</div>
                      </div>
                    </div>
                  </div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="text-sm leading-5 text-gray-900">Terra</div>
                  <div class="text-sm leading-5 text-gray-900">Xava</div>
                  <div class="text-sm leading-5 text-gray-900">Micha</div>
                  <div class="text-sm leading-5 text-gray-900">Holytank</div>
                  <div class="text-sm leading-5 text-gray-900">Tova</div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="text-sm">
                    ||||||||||||||||||||||||||||||||||||||||||||||
                  </div>
                  <div class="text-sm">
                    |||||||||||||||||||||||||||||||||||||||
                  </div>
                  <div class="text-sm">|||||||||||||||||||||||||||||</div>
                  <div class="text-sm">||||||||||||||||||</div>
                  <div class="text-sm">|||||||||||</div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="text-sm leading-5 text-gray-900">56.535</div>
                  <div class="text-sm leading-5 text-gray-900">46.535</div>
                  <div class="text-sm leading-5 text-gray-900">36.535</div>
                  <div class="text-sm leading-5 text-gray-900">26.535</div>
                  <div class="text-sm leading-5 text-gray-900">16.535</div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap text-right border-b border-gray-200 text-sm leading-5 font-medium"
                >
                  <a href="#" class="text-indigo-600 hover:text-indigo-900"
                    >more</a
                  >
                </td>
              </tr>
              <tr>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="flex items-center">
                    <div class="ml-4">
                      <div class="text-sm leading-5 font-medium text-gray-900">
                        Atal'Dazar <span class="text-gray-500">+23*</span>
                      </div>
                      <div class="text-sm leading-5 text-gray-500">
                        <div>affixes:</div>
                        <div>34:23 +3:12</div>
                        <div>4 deaths</div>
                      </div>
                    </div>
                  </div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="text-sm leading-5 text-gray-900">Terra</div>
                  <div class="text-sm leading-5 text-gray-900">Xava</div>
                  <div class="text-sm leading-5 text-gray-900">Micha</div>
                  <div class="text-sm leading-5 text-gray-900">Holytank</div>
                  <div class="text-sm leading-5 text-gray-900">Tova</div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <!-- dont show on small -->
                  <div class="text-sm">
                    ||||||||||||||||||||||||||||||||||||||||||||||
                  </div>
                  <div class="text-sm">
                    |||||||||||||||||||||||||||||||||||||||
                  </div>
                  <div class="text-sm">|||||||||||||||||||||||||||||</div>
                  <div class="text-sm">||||||||||||||||||</div>
                  <div class="text-sm">|||||||||||</div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap border-b border-gray-200"
                >
                  <div class="text-sm leading-5 text-gray-900">56.535</div>
                  <div class="text-sm leading-5 text-gray-900">46.535</div>
                  <div class="text-sm leading-5 text-gray-900">36.535</div>
                  <div class="text-sm leading-5 text-gray-900">26.535</div>
                  <div class="text-sm leading-5 text-gray-900">16.535</div>
                </td>
                <td
                  class="px-6 py-4 whitespace-no-wrap text-right border-b border-gray-200 text-sm leading-5 font-medium"
                >
                  <a href="#" class="text-indigo-600 hover:text-indigo-900"
                    >more</a
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
  layout: 'mythicplus'
}
</script>

<style></style>
