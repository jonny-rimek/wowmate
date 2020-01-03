<template>
	<div class="">
	  	Player: {{ items[0].caster_name }}
		<table>
			<tr>
				<th>Damage</th>
				<th>Encounter ID</th>
			<tr v-for="item in items">
				<td><router-link :to="{name: 'EncounterTopDamage', params: {encounter_id: item.gsi1pk }}">{{ item.gsi1pk}}</router-link></td>
				<td><router-link :to="{name: 'BossFightDamage', params: {boss_fight_uuid: item.gsi2pk }}">{{ item.sk }}</router-link></td>
			</tr>
		</table>
	</div>
</template>

<script>
import axios from 'axios';

export default {
	name: 'PlayerDamage',
	data() {
		return {
			items: null
		}
	},
	mounted() {
		axios.get('https://api.wmate.net/api/damage/caster/' + this.$route.params.caster_id).then(response => (this.items = response.data ))
	}
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
