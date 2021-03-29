<template>
  <nuxt-link
    :to="{ name: link }"
    active-class="font-bold border-red-800"
    class="text-gray-300 hover:border-gray-400 focus:outline-none transition duration-150 ease-in-out"
  >
    <slot></slot>
  </nuxt-link>
</template>
<script>
export default {
  props: {
    link: {
      type: String,
      required: true
    },
  }
}
</script>
<style></style>
