export const state = () => ({
	next: "",
  prev: "",
})

export const mutations = {
	set(state, sk) {
		state.next = sk.last
    state.prev = sk.first
	}
}

export const getters = {
	get (state) {
		return state
	}
}
