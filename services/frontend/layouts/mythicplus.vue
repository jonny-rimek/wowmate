<template>
  <div class="flex flex-col h-screen justify-between">
    <div>
      <Navbar> </Navbar>
      <div class="max-w-7xl mx-auto md:px-6 lg:px-9">
        <div class="flex flex-shrink-0">
          <!-- DESKTOP ONLY SIDEBAR -->
          <!-- THE MOBILE VERSION OF THE SIDEBAR IS PART OF NAVBAR COMPONENT -->
          <div class="hidden lg:flex flex-col w-64">
            <div class="h-0 flex-1 flex flex-col pt-5 pb-4 overflow-y-auto">
              <nav>
                <SidebarLink
                  v-for="item in dungeons"
                  :param="item.id"
                  :pathName="item.pathName"
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
      <Banner />
    </div>
    <Footer />
  </div>
</template>

<script>
import Navbar from '~/components/Navbar.vue'
import SidebarLink from '~/components/SidebarLink.vue'
import Banner from '~/components/Banner.vue'
import Footer from '~/components/Footer.vue'

export default {
  components: {
    Navbar,
    SidebarLink,
    Banner,
    Footer
  },
  head() {
    return {
      title: 'wowmate'
    }
  },
  data() {
    return {
      dungeons: [
        { name: 'All Dungeons', pathName: 'mythicplus' },
        { name: "De Other Site", pathName: 'mythicplus-id', id: 2291 },
        { name: "Plaguefall", pathName: 'mythicplus-id', id: 2289 },
        { name: "Sanguine Depths", pathName: 'mythicplus-id', id: 2284 },
        { name: "Halls of Atonement", pathName: 'mythicplus-id', id: 2287 },
        { name: "Mists of Tirna Scithe", pathName: 'mythicplus-id', id: 2290 },
        { name: "Spires of Ascension", pathName: 'mythicplus-id', id: 2285 },
        { name: "Necrotic Wake", pathName: 'mythicplus-id', id: 2286 },
        { name: "Theater of Pain", pathName: 'mythicplus-id', id: 2293 },
      ]
    }
  }
}
</script>

<style></style>
