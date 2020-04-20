<template>
  <nuxt-link
    :to="{ name: 'mythicplus-id', params: { id: link } }"
    :prefetch="false"
    class="mt-1 group flex items-center px-3 py-2 text-sm leading-5 font-medium text-gray-500 rounded-md hover:text-gray-700 hover:bg-gray-50 focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150"
  >
    <!-- active-class="text-gray-700 bg-gray-100 hover:bg-gray-200" -->
    <span class="truncate">
      <slot></slot>
    </span>
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
