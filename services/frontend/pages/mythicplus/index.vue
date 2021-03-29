<template>
  <TopMPlus :logs="logs" :title="title"></TopMPlus>
</template>

<script>
export default {
  watchQuery: true,
  async asyncData(context) {
    const logs = await context.app.$wowmateApi.getSummaries()
    return { logs }
  },
  layout: 'mythicplus',
  data() {
    return {
      logs: {},
      title: 'Not All dungeons',
    }
  },
}
</script>

<style></style>
