<template>
  <div class="container">
    wowmate
  </div>
</template>

<script>
// import Logo from '~/components/Logo.vue'

export default {
  layout: 'navbar',
  layout(context) {
    return 'navbar'
  },
  components: {
    // Logo
  }
}
</script>

<style>
/* Sample `apply` at-rules with Tailwind CSS
.container {
  @apply min-h-screen flex justify-center items-center text-center mx-auto;
}
*/
</style>
