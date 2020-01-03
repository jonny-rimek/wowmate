<template>
	<div id="app" class="min-h-screen bg-orange-200">
		<nav class="">
			<div class="font-bold text-4xl text-orange-900">
				<a href="/">wowmate</a>
			</div>
		</nav>
		<div class="flex">
			<div class="block">
				<div class="w-48">
					<div>
						<a href="/raid/tep">The Eternal palace</a>
					</div>
					<div>
						<a href="/encounter/2298">Commander Shivra</a>
					</div>
					<router-link :to="{name: 'EncounterTopDamage', params: { encounterId: '2298'}}">Commander Shivra</router-link>
					<div>
						<a href="/encounter/2289">Blackwater Behemoth</a>
					</div>
					<div>
						<a href="/encounter/2305">Radiance of Azashara</a>
					</div>
					<div>
						<a href="/encounter/2304">Lady Ashvane</a>
					</div>
					<div>
						<a href="/encounter/2303">Orgozoa</a>
					</div>
					<div>
						<a href="/encounter/2311">The Queen's Court</a>
					</div>
					<div>
						<a href="/encounter/2293">Za'qul</a>
					</div>
					<div>
						<a href="/encounter/aszhara">Queen Azshara</a>
					</div>
				</div>
			</div>
			<router-view></router-view>
		</div>
	</div>
</template>

<script>
export default {
	name: 'app',
	components: {
	}
}
</script>

<style src="./assets/tailwind.css">
