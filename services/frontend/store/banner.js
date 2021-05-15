export const state = () => ({
  visible: showBanner()
})

export const mutations = {
	hide(state) {
	  localStorage.setItem("banner", "hide")
		state.visible = false
	}
}

export const getters = {
	get (state) {
		return state.visible
	}
}

// showBanner checks the localStorage for the banner item, if it exists, the banner is not shown
// if it doesn't it is
// doesn't work with SSR
function showBanner() {
  const banner = localStorage.getItem("banner")
  if (banner === null) {
    return true;
  }
  return false
}
