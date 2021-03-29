export default function(context, inject){

  inject('wowmateApi', {
    getSummaries: getKeys,
    getPerDungeonSummaries: getKeysPerDungeon,
    getMythicplusPlayerOverallDamage: getPlayerDamageDone,
  })

  async function getPlayerDamageDone() {
    const log = await context.app.$axios.$request({
      url: process.env.baseUrl + '/api/combatlogs/keys/' + context.route.params.id + '/player-damage-done',
      method: "get",
    })

    return log
  }

    async function getKeysPerDungeon(){
    let logs
    if (context.route.query.next != null && context.route.prev == null) {
      logs = await context.app.$axios.$request({
        url: process.env.baseUrl + '/api/combatlogs/keys/' + context.route.params.id,
        method: "get",
        params: {
          next: encodeURIComponent(context.route.query.next)
        }
      })
    } else if (context.route.query.prev != null && context.route.next == null) {
      logs = await context.app.$axios.$request({
        url: process.env.baseUrl + '/api/combatlogs/keys/' + context.route.params.id,
        method: "get",
        params: {
          prev: encodeURIComponent(context.route.query.prev)
        }
      })
    } else {
      logs = await context.app.$axios.$request({
        url: process.env.baseUrl + '/api/combatlogs/keys/' + context.route.params.id,
        method: "get",
      })
    }
    let lastSk
    let firstSk

    if (logs.last_sk == null) {
      lastSk = ""
    } else {
      lastSk = logs.last_sk
    }
    if (logs.first_sk == null) {
      firstSk = ""
    } else {
      firstSk = logs.first_sk
    }
    //console.log(firstSk)
    //console.log(lastSk)
    context.store.commit('pagination/set', {
      last: lastSk,
      first: firstSk
      //last: logs.last_sk.Value,
      //first: logs.first_sk.Value
    })
    return logs
  }
  async function getKeys() {
    let logs
    if (context.route.query.next != null && context.route.prev == null) {
      logs = await context.app.$axios.$request({
        url: process.env.baseUrl + '/api/combatlogs/keys',
        method: "get",
        params: {
          next: encodeURIComponent(context.route.query.next)
        }
      })
    } else if (context.route.query.prev != null && context.route.next == null) {
      logs = await context.app.$axios.$request({
        url: process.env.baseUrl + '/api/combatlogs/keys',
        method: "get",
        params: {
          prev: encodeURIComponent(context.route.query.prev)
        }
      })
    } else {
      logs = await context.app.$axios.$request({
        url: process.env.baseUrl + '/api/combatlogs/keys',
        method: "get",
      })
    }
    let lastSk
    let firstSk

    if (logs.last_sk == null) {
      lastSk = ""
    } else {
      lastSk = logs.last_sk
    }
    if (logs.first_sk == null) {
      firstSk = ""
    } else {
      firstSk = logs.first_sk
    }
    //console.log(firstSk)
    //console.log(lastSk)
    context.store.commit('pagination/set', {
      last: lastSk,
      first: firstSk
      //last: logs.last_sk.Value,
      //first: logs.first_sk.Value
    })
    return logs
  }
  /*


   */
}
