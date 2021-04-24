<template>
  <TopMPlus :logs="logs" :title="title"></TopMPlus>
</template>

<script>

export default {
  watchQuery: true,
  async asyncData(context) {
    const logs = await context.app.$wowmateApi.getPerDungeonSummaries()
    return { logs }
  },
  layout: 'mythicplus', //trigger!
  data() {
    return {
      logs: {},
      title: 'Not All dungeons',
    }
  },
}
</script>

<style></style>
