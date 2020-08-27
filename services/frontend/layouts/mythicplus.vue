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
                  pathName="mythicplus-id"
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
import MobileNavbarLink from '~/components/MobileNavbarLink.vue'
import SidebarLink from '~/components/SidebarLink.vue'
import Banner from '~/components/Banner.vue'
import Footer from '~/components/Footer.vue'

export default {
  components: {
    Navbar,
    MobileNavbarLink,
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
        { name: "Atal'Dazar", id: 2144 },
        { name: 'Freehold', id: 2145 },
        { name: "Kings'Rest", id: 2146 },
        { name: 'Shrine of the Storm', id: 2147 },
        { name: 'Siege of Boralus', id: 2148 },
        { name: 'Temple of Sethraliss', id: 2149 },
        { name: 'The MOTHERLODE!!', id: 2154 },
        { name: 'The Underrot', id: 2164 },
        { name: 'Tol Dagor', id: 2174 },
        { name: 'Waycrest Manor', id: 2184 },
        { name: 'Operation: Mechagon - Junkyard', id: 2194 },
        { name: 'Operation: Mechagon - Workshop', id: 2104 }
      ]
    }
  }
}
</script>

<style></style>
