export default {
  ssr: false,
  /*
  ** Headers of the page
  */
  server: {
    port: 3003
  },
  env: {
    baseUrl:
      process.env.NUXT_ENV === 'local'
        ? 'http://localhost:3000'
        : process.env.NUXT_ENV === 'dev'
        ? 'https://api.dev.wowmate.io'
        : 'https://api.wowmate.io'
  },
  components: true,
  head: {
    title: process.env.npm_package_name || '',
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { hid: 'description', name: 'description', content: process.env.npm_package_description || '' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
    ]
  },
  /*
  ** Customize the progress-bar color
  */
  loading: { color: '#E64E58' },
  /*
  ** Global CSS
  */
  css: [
  ],
  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
    "~/plugins/wowmateApi"
  ],
  /*
  ** Nuxt.js dev-modules
  */
  buildModules: [
	'@nuxtjs/tailwindcss',
  ],
  tailwindcss: {
    jit: true
  },
  /*
  ** Nuxt.js modules
  */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
	'@nuxtjs/axios',
	'@nuxt/content',
  ],
  /*
  ** Axios module configuration
  ** See https://axios.nuxtjs.org/options
  */
  axios: {
  },
  /*
  ** Build configuration
  */
  build: {
    /*
    ** You can extend webpack config here
    */
    extend (config, ctx) {
    }
  },
  //purgeCSS: {
  //  whitelist: ["dark-mode"]
  //}
}
