<template>
  <NavbarLink
    :link="link"
    class="block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
  >
    <slot></slot>
  </NavbarLink>
</template>
<script>
import NavbarLink from '~/components/NavbarLink.vue'
export default {
  components: {
    NavbarLink
  },
  props: {
    link: {
      type: String,
      required: true
    }
  }
}
</script>
<style></style>
