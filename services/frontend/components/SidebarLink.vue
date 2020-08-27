<template>
  <nuxt-link
    :to="{ name: pathName, params: { id: param } }"
    active-class="text-gray-700 bg-gray-200 dark:bg-gray-700 dark:text-gray-200 font-bold"
    class="mt-1 group flex items-center px-2 py-2 text-sm leading-5 font-medium dark:text-gray-300 text-gray-500 hover:bg-gray-200 dark-hover:bg-gray-700 dark-hover:text-gray-300 hover:text-gray-700 rounded-md focus:outline-none focus:bg-gray-100 transition ease-in-out duration-150"
  >
    <span class="truncate">
      <slot></slot>
    </span>
  </nuxt-link>
</template>
<script>
export default {
  props: {
    param: {
      type: Number
    },
    pathName: {
      type: String,
      required: true
    }
  }
}
</script>
<style></style>
