<template>
  <div class="p-5 lg:hidden">
    <SidebarLink
      v-for="item in dungeons"
      :param="item.id"
      :pathName="item.pathName"
      :key="item.id"
      data-cy='dungeon-mobile'
    >
      {{ item.name }}
    </SidebarLink>
  </div>
</template>

<script>

export default {
  layout: 'mythicplus',
  data() {
    return {
      dungeons: [
        { name: 'All dungeons', pathName: 'mythicplus' },
        { name: "Atal'Dazar", pathName: 'mythicplus-id', id: 2144 },
        { name: 'Freehold', pathName: 'mythicplus-id', id: 2145 },
        { name: "Kings'Rest", pathName: 'mythicplus-id', id: 2146 },
        {
          name: 'Shrine of the Storm',
          pathName: 'mythicplus-id',
          id: 2147
        },
        {
          name: 'Siege of Boralus',
          pathName: 'mythicplus-id',
          id: 2148
        },
        {
          name: 'Temple of Sethraliss',
          pathName: 'mythicplus-id',
          id: 2149
        },
        {
          name: 'The MOTHERLODE!!',
          pathName: 'mythicplus-id',
          id: 2154
        },
        { name: 'The Underrot', pathName: 'mythicplus-id', id: 2164 },
        { name: 'Tol Dagor', pathName: 'mythicplus-id', id: 2174 },
        { name: 'Waycrest Manor', pathName: 'mythicplus-id', id: 2184 },
        {
          name: 'Operation: Mechagon - Junkyard',
          pathName: 'mythicplus-id',
          id: 2194
        },
        {
          name: 'Operation: Mechagon - Workshop',
          pathName: 'mythicplus-id',
          id: 2104
        }
      ]
    }
  }
}
</script>

<style></style>
