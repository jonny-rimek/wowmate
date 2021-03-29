<template>
  <div class="flex flex-col h-screen justify-between">
    <div>
      <Navbar />
      <div class="mb-auto max-w-7xl mx-auto md:px-6 lg:px-9">
        <nuxt />
      </div>
      <Banner />
    </div>
    <Footer />
  </div>
</template>

<script>
import Navbar from '~/components/Navbar.vue'
import Banner from '~/components/Banner.vue'
import Footer from '~/components/Footer.vue'

export default {
  components: {
    Navbar,
    Footer,
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
