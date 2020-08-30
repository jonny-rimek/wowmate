<template>
  <div class="p-5 lg:hidden">
    <SidebarLink
      v-for="item in dungeons"
      :param="item.id"
      :pathName="item.pathName"
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
        { name: 'All dungeons', pathName: 'mythicplus-dungeon' },
        { name: "Atal'Dazar", pathName: 'mythicplus-dungeon-id', id: 2144 },
        { name: 'Freehold', pathName: 'mythicplus-dungeon-id', id: 2145 },
        { name: "Kings'Rest", pathName: 'mythicplus-dungeon-id', id: 2146 },
        {
          name: 'Shrine of the Storm',
          pathName: 'mythicplus-dungeon-id',
          id: 2147
        },
        {
          name: 'Siege of Boralus',
          pathName: 'mythicplus-dungeon-id',
          id: 2148
        },
        {
          name: 'Temple of Sethraliss',
          pathName: 'mythicplus-dungeon-id',
          id: 2149
        },
        {
          name: 'The MOTHERLODE!!',
          pathName: 'mythicplus-dungeon-id',
          id: 2154
        },
        { name: 'The Underrot', pathName: 'mythicplus-dungeon-id', id: 2164 },
        { name: 'Tol Dagor', pathName: 'mythicplus-dungeon-id', id: 2174 },
        { name: 'Waycrest Manor', pathName: 'mythicplus-dungeon-id', id: 2184 },
        {
          name: 'Operation: Mechagon - Junkyard',
          pathName: 'mythicplus-dungeon-id',
          id: 2194
        },
        {
          name: 'Operation: Mechagon - Workshop',
          pathName: 'mythicplus-dungeon-id',
          id: 2104
        }
      ]
    }
  }
}
</script>

<style></style>
