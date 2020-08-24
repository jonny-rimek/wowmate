<template>
  <div>
    <Navbar />
    <div class="max-w-7xl mx-auto md:px-6 lg:px-9">
      <select
        v-model="$colorMode.preference"
        class="border w-24 h-8 dark:bg-gray-900 dark:border-gray-700"
      >
        <option value="system">System</option>
        <option value="light">Light</option>
        <option value="dark">Dark</option>
      </select>
      <nuxt />
    </div>
  </div>
</template>

<script>
import Navbar from '~/components/Navbar.vue'

export default {
  components: {
    Navbar
  },
  head() {
    return {
      title: 'wowmate'
    }
  }
}
</script>

<style></style>
