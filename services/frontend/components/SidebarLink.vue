<template>
  <nuxt-link
    :to="{ name: pathName, params: { id: param } }"
    exact-active-class="bg-gray-700 text-gray-200 font-bold"
    class="mt-1 group flex items-center px-2 py-2 text-sm leading-5 font-medium text-gray-300 hover:bg-gray-700 rounded-md focus:outline-none transition ease-in-out duration-150"
    data-cy='dungeon'
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
