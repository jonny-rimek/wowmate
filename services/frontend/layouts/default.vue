<template>
  <div>
    <Navbar />
    <div class="max-w-7xl mx-auto md:px-6 lg:px-9">
      <nuxt />
    </div>
    <Banner />
  </div>
</template>

<script>
import Navbar from '~/components/Navbar.vue'
import Banner from '~/components/Banner.vue'

export default {
  components: {
    Navbar,
    Banner
  },
  head() {
    return {
      title: 'wowmate'
    }
  }
}
</script>

<style></style>
