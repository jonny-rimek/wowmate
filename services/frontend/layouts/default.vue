<template>
  <div>
    <Navbar />
    <div class="max-w-7xl mx-auto md:px-6 lg:px-9">
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
