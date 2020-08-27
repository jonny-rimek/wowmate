<template>
  <div class="p-5 lg:hidden">
    <SidebarLink pathName="mythicplus">
      All Dungeons
    </SidebarLink>
    <SidebarLink
      v-for="item in dungeons"
      :param="item.id"
      pathName="mythicplus-id"
      :key="item.id"
    >
      {{ item.name }}
    </SidebarLink>
  </div>
</template>

<script>
import SidebarLink from '~/components/SidebarLink.vue'

export default {
  components: {
    SidebarLink
  },
  layout: 'mythicplus',
  data() {
    return {
      dungeons: [
        { name: "Atal'Dazar", id: 2144 },
        { name: 'Freehold', id: 2145 },
        { name: "Kings'Rest", id: 2146 },
        { name: 'Shrine of the Storm', id: 2147 },
        { name: 'Siege of Boralus', id: 2148 },
        { name: 'Temple of Sethraliss', id: 2149 },
        { name: 'The MOTHERLODE!!', id: 2154 },
        { name: 'The Underrot', id: 2164 },
        { name: 'Tol Dagor', id: 2174 },
        { name: 'Waycrest Manor', id: 2184 },
        { name: 'Operation: Mechagon - Junkyard', id: 2194 },
        { name: 'Operation: Mechagon - Workshop', id: 2104 }
      ]
    }
  }
}
</script>

<style></style>
