/*
** TailwindCSS Configuration File
**
** Docs: https://tailwindcss.com/docs/configuration
** Default: https://github.com/tailwindcss/tailwindcss/blob/master/stubs/defaultConfig.stub.js
*/
module.exports = {
  theme: {
   extend: {
    colors: {
      red: {
		'800': '#E64E58',
		'600': '#EB7179',
		'400': '#F0949B',
		'200': '#F5B8BC',
        }
      }
    }
  },
  variants: {},
  plugins: [
	  //the import definately works, warning is ignorable
    require('@tailwindcss/ui')({
	  layout: 'sidebar',
	})
  ]
}
