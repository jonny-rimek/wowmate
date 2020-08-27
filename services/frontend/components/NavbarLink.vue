<template>
  <nuxt-link
    :to="{ name: link }"
    :prefetch="false"
    active-class="border-red-800 hover:border-red-800 text-gray-700"
    class="text-gray-500 dark:text-gray-300 dark-hover:text-gray-300 hover:text-gray-700 hover:border-gray-400 focus:outline-none transition duration-150 ease-in-out"
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
