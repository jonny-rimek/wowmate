<template>
  <div>
    <nav class="bg-white shadow">
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
                  src="/wowmate_logo_full_transparent.png"
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
                  class="h-12 w-auto"
                  src="/wowmate_logo_full_transparent.png"
                  alt=""
                />
              </nuxt-link>
            </div>
            <!-- NAV LINKS -->
            <div class="hidden lg:ml-6 lg:flex">
              <NavbarLink
                link="mythicplus"
                class="inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5"
              >
                Mythic+
              </NavbarLink>
              <NavbarLink
                link="nyalotha"
                class="ml-8 inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium leading-5"
              >
                Ny'alotha
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
                  class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:border-red-400 focus:shadow-outline-red sm:text-sm transition duration-150 ease-in-out"
                  placeholder="Search"
                  type="search"
                />
              </div>
            </div>
          </div>
          <!-- MOBILE MENU BUTTON -->
          <!-- @click="toggleMobileMenu()" -->
          <div
            @click="showMobileMenu"
            v-on-clickaway="toggleMobileMenu"
            class="flex items-center lg:hidden"
          >
            <button
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 focus:text-gray-500 transition duration-150 ease-in-out"
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
          <!-- NOTIFICATIONS -->
          <!--
        <div class="hidden lg:ml-4 lg:flex lg:items-center">
          <button
            class="flex-shrink-0 p-1 border-2 border-transparent text-gray-400 rounded-full hover:text-gray-500 focus:outline-none focus:text-gray-500 focus:bg-gray-100 transition duration-150 ease-in-out"
          >
            <svg
              class="h-6 w-6"
              stroke="currentColor"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
              />
            </svg>
          </button>
-->

          <!-- Profile dropdown -->
          <!--
          <div class="ml-4 relative flex-shrink-0">
            <div>
              <button
                class="flex text-sm border-2 border-transparent rounded-full focus:outline-none focus:border-gray-300 transition duration-150 ease-in-out"
              >
                <img
                  class="h-8 w-8 rounded-full"
                  src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
                  alt=""
                />
              </button>
            </div>
			-->

          <!--
            Profile dropdown panel, show/hide based on dropdown state.

            Entering: "transition ease-out duration-100"
              From: "transform opacity-0 scale-95"
              To: "transform opacity-100 scale-100"
            Leaving: "transition ease-in duration-75"
              From: "transform opacity-100 scale-100"
              To: "transform opacity-0 scale-95"
          -->
          <!--
            <div
              class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg"
            >
              <div class="py-1 rounded-md bg-white shadow-xs">
                <a
                  href="#"
                  class="block px-4 py-2 text-sm leading-5 text-gray-700 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 transition duration-150 ease-in-out"
                  >Your Profile</a
                >
                <a
                  href="#"
                  class="block px-4 py-2 text-sm leading-5 text-gray-700 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 transition duration-150 ease-in-out"
                  >Settings</a
                >
                <a
                  href="#"
                  class="block px-4 py-2 text-sm leading-5 text-gray-700 hover:bg-gray-100 focus:outline-none focus:bg-gray-100 transition duration-150 ease-in-out"
                  >Sign out</a
                >
              </div>
            </div>
          </div>
        </div>
		-->
        </div>
      </div>

      <!--
    Mobile menu, toggle classes based on menu state.

    Menu open: "block", Menu closed: "hidden"
  -->
      <div
        v-if="mobileMenu"
        v-on-clickaway="toggleMobileMenu"
        class="lg:hidden"
      >
        <div class="pt-2 pb-3">
          <NavbarLink
            link="mythicplus"
            class="block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            Mythic+
          </NavbarLink>
          <NavbarLink
            link="nyalotha"
            class="block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            Ny'alotha
          </NavbarLink>
          <NavbarLink
            link="news"
            class="block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            News
          </NavbarLink>
          <NavbarLink
            link="upload"
            class="block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            Upload
          </NavbarLink>
        </div>
        <!--
      <div class="pt-4 pb-3 border-t border-gray-200">
        <div class="flex items-center px-4">
          <div class="flex-shrink-0">
            <img
              class="h-10 w-10 rounded-full"
              src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80"
              alt=""
            />
          </div>
          <div class="ml-3">
            <div class="text-base font-medium leading-6 text-gray-800">
              Tom Cook
            </div>
            <div class="text-sm font-medium leading-5 text-gray-500">
              tom@example.com
            </div>
          </div>
        </div>
        <div class="mt-3">
          <a
            href="#"
            class="block px-4 py-2 text-base font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-100 focus:outline-none focus:text-gray-800 focus:bg-gray-100 transition duration-150 ease-in-out"
            >Your Profile</a
          >
          <a
            href="#"
            class="mt-1 block px-4 py-2 text-base font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-100 focus:outline-none focus:text-gray-800 focus:bg-gray-100 transition duration-150 ease-in-out"
            >Settings</a
          >
          <a
            href="#"
            class="mt-1 block px-4 py-2 text-base font-medium text-gray-600 hover:text-gray-800 hover:bg-gray-100 focus:outline-none focus:text-gray-800 focus:bg-gray-100 transition duration-150 ease-in-out"
            >Sign out</a
          >
        </div>
      </div>
	  -->
      </div>
    </nav>
    <nuxt />
  </div>
</template>

<script>
import NavbarLink from '~/components/NavbarLink.vue'
import { mixin as clickaway } from 'vue-clickaway' //IMPROVE: find a native solution

export default {
  mixins: [clickaway],
  components: {
    NavbarLink
  },
  data() {
    return {
      mobileMenu: false
    }
  },
  methods: {
    showMobileMenu() {
      this.mobileMenu = true
    },
    toggleMobileMenu() {
      this.mobileMenu = false
    }
  }
}
</script>

<style></style>
