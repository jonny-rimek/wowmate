<template>
  <div>
    <Navbar />
    <div class="max-w-7xl mx-auto md:px-6 lg:px-9">
      <div class="flex flex-shrink-0">
        <div class="hidden md:flex flex-col w-64">
          <div class="h-0 flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
            <nav>
              <SidebarLink
                v-for="item in nyalotha"
                :link="item.id"
                :pathName="pathName"
                :key="item.id"
              >
                {{ item.name }}
              </SidebarLink>
            </nav>
          </div>
        </div>

        <nuxt />
      </div>
    </div>
  </div>
</template>

<script>
import Navbar from '~/components/Navbar.vue'
import SidebarLink from '~/components/SidebarLink.vue'

export default {
  components: {
    Navbar,
    SidebarLink
  },
  head() {
    return {
      title: 'wowmate'
    }
  },
  data() {
    return {
      nyalotha: [
        { name: 'Wrathion', id: 2329 },
        { name: 'Maut', id: 2327 },
        { name: 'Prophet Skitra', id: 2334 },
        { name: 'Dark Inquisitor Xanesh', id: 2328 },
        { name: 'The Hivemind', id: 2333 },
        { name: "Shad'har the Insatiable", id: 2335 },
        { name: "Drest'agath", id: 2343 },
        { name: 'Vexiona', id: 2336 },
        { name: 'Ra-den the Despoiled', id: 2331 },
        { name: "Il'ganoth, Corruption Reborn", id: 2345 },
        { name: "Carapace of N'Zoth", id: 2337 },
        { name: "N'Zoth the Corruptor", id: 2344 }
      ],
      pathName: 'nyalotha-id'
    }
  }
}
</script>

<style></style>
