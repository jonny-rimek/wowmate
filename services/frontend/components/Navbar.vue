<template>
  <div class="">
    <nav class="border-b dark:border-gray-600">
      <div class="max-w-7xl mx-auto px-2 sm:px-4 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex px-2 lg:px-0">
            <!-- ICONS -->
            <div class="flex-shrink-0 flex items-center">
              <nuxt-link
                :to="{ name: 'index' }"
                :prefetch="false"
                class="block focus:outline-none lg:hidden"
              >
                <img
                  class="h-12 focus:outline-none w-auto"
                  src="/wowmate_logo_only_transparent.png"
                  alt=""
                />
              </nuxt-link>
              <!-- TODO: make a icon only version -->
              <nuxt-link
                :to="{ name: 'index' }"
                :prefetch="false"
                class="hidden focus:outline-none lg:block"
              >
                <img
                  class="h-8 w-auto"
                  src="/wowmate_logo_only_transparent.png"
                  alt=""
                />
              </nuxt-link>
              <nuxt-link
                :to="{ name: 'index' }"
                :prefetch="false"
                class="hidden focus:outline-none lg:block"
              >
                <span class="pl-2 text-red-800 text-l font-sans font-semibold">
                  wowmate
                </span>
              </nuxt-link>
            </div>
            <!-- NAV LINKS -->
            <div class="hidden lg:ml-6 lg:flex">
              <NavbarLink
                prefetch="true"
                link="mythicplus"
                class="inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5"
              >
                Mythic+
              </NavbarLink>
              <NavbarLink
                link="news"
                class="ml-8 inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5"
              >
                News
              </NavbarLink>
              <NavbarLink
                link="upload"
                class="ml-8 inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5"
              >
                Upload
              </NavbarLink>
            </div>
          </div>
          <!-- SEARCH -->
          <div
            class="flex-1 flex items-center justify-center px-2 lg:ml-6 lg:justify-end"
          >
            <div class="max-w-lg w-full lg:max-w-xs">
              <!--
              <label for="search" class="sr-only">Search</label>
              <div class="relative">
                <div
                  class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                >
                  <svg
                    class="h-5 w-5 text-gray-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <input
                  id="search"
                  class="block w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-500 rounded-md leading-5 bg-white dark:bg-gray-700 placeholder-gray-300 focus:outline-none focus:placeholder-gray-100 focus:border-gray-200 focus:shadow-outline-red sm:text-sm transition duration-150 ease-in-out"
                  placeholder="Search"
                  type="search"
                />
              </div>
			  -->
            </div>
          </div>
          <!-- MOBILE MENU BUTTON -->
          <!--  -->
          <div
            v-on-clickaway="hideMobileMenu"
            @click="toggleMobileMenu"
            class="flex items-center lg:hidden"
          >
            <button
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 dark-hover:bg-gray-700 dark-hover:text-gray-300 focus:outline-none focus:bg-gray-100 focus:text-gray-500 transition duration-150 ease-in-out"
            >
              <!-- Icon when menu is closed. -->
              <!-- Menu open: "hidden", Menu closed: "block" -->
              <svg
                class="block h-6 w-6"
                stroke="currentColor"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
              <!-- Icon when menu is open. -->
              <!-- Menu open: "block", Menu closed: "hidden" -->
              <svg
                class="hidden h-6 w-6"
                stroke="currentColor"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!--MOBILE MENU-->
      <div :class="mobileMenu ? 'block' : 'hidden'" class="lg:hidden">
        <div class="pt-2 pb-3">
          <MobileNavbarLink link="mythicplus-mobile">
            Mythic+
            <!--TODO: Add all m+ links and expand on click, add overall option  -->
          </MobileNavbarLink>
          <MobileNavbarLink link="news">
            News
          </MobileNavbarLink>
          <MobileNavbarLink link="upload">
            Upload
          </MobileNavbarLink>
        </div>
      </div>
    </nav>
  </div>
  <!-- END TOP NAVBARR -->
</template>

<script>
import NavbarLink from '~/components/NavbarLink.vue'
import MobileNavbarLink from '~/components/MobileNavbarLink.vue'
import { mixin as clickaway } from 'vue-clickaway' //IMPROVE: find a native solution

export default {
  mixins: [clickaway],
  components: {
    NavbarLink,
    MobileNavbarLink
  },
  data() {
    return {
      mobileMenu: false
    }
  },
  methods: {
    toggleMobileMenu() {
      this.mobileMenu = !this.mobileMenu
    },
    hideMobileMenu() {
      this.mobileMenu = false
    }
  }
}
</script>

<style></style>
