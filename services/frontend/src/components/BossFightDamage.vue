<template>
  <div class="">
	  Bossfight:
		<table>
			<tr>
				<th>Player</th>
				<th>Damage</th>
			</tr>
			<tr v-for="item in items">
				<td><router-link :to="{name: 'PlayerDamage', params: { caster_id: item.gsi3pk }}">{{ item.caster_name }}</router-link></td>
				<td><router-link :to="{name: 'BossFightDamage', params: {boss_fight_uuid: item.gsi2pk }}">{{ item.sk }}</router-link></td>
			</tr>
		</table>
	</div>
</template>

<script>
import axios from 'axios';

export default {
	name: 'BossFightDamageDamage',
	data() {
		return {
			items: null
		}
	},
	mounted() {
		axios.get('https://api.wmate.net/api/damage/boss-fight/' + this.$route.params.boss_fight_uuid).then(response => (this.items = response.data ))
	}
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
