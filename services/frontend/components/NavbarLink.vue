<template>
  <nuxt-link
    :to="{ name: link }"
    :prefetch="false"
    active-class="border-red-800 text-gray-700 hover:border-red-800"
    class="text-gray-600 hover:text-gray-700 hover:border-gray-500 focus:outline-none transition duration-150 ease-in-out"
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
    }
  }
}
</script>
<style></style>
