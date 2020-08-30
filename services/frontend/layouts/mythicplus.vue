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
        { name: 'All dungeons', pathName: 'mythicplus-dungeon' },
        { name: "Atal'Dazar", pathName: 'mythicplus-dungeon-id', id: 2144 },
        { name: 'Freehold', pathName: 'mythicplus-dungeon-id', id: 2145 },
        { name: "Kings'Rest", pathName: 'mythicplus-dungeon-id', id: 2146 },
        {
          name: 'Shrine of the Storm',
          pathName: 'mythicplus-dungeon-id',
          id: 2147
        },
        {
          name: 'Siege of Boralus',
          pathName: 'mythicplus-dungeon-id',
          id: 2148
        },
        {
          name: 'Temple of Sethraliss',
          pathName: 'mythicplus-dungeon-id',
          id: 2149
        },
        {
          name: 'The MOTHERLODE!!',
          pathName: 'mythicplus-dungeon-id',
          id: 2154
        },
        { name: 'The Underrot', pathName: 'mythicplus-dungeon-id', id: 2164 },
        { name: 'Tol Dagor', pathName: 'mythicplus-dungeon-id', id: 2174 },
        { name: 'Waycrest Manor', pathName: 'mythicplus-dungeon-id', id: 2184 },
        {
          name: 'Operation: Mechagon - Junkyard',
          pathName: 'mythicplus-dungeon-id',
          id: 2194
        },
        {
          name: 'Operation: Mechagon - Workshop',
          pathName: 'mythicplus-dungeon-id',
          id: 2104
        }
      ]
    }
  }
}
</script>

<style></style>
