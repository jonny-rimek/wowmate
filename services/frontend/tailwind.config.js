/*
** TailwindCSS Configuration File
**
** Docs: https://tailwindcss.com/docs/configuration
** Default: https://github.com/tailwindcss/tailwindcss/blob/master/stubs/defaultConfig.stub.js
*/
module.exports = {
  theme: {
   darkSelector: ".dark-mode",
   extend: {
    colors: {
      red: {
		'800': '#E64E58',
		'600': '#EB7179',
		'400': '#F0949B',
		'200': '#F5B8BC',
		'50': '#fdf2f2',
        }
      }
    }
  },
  variants: {
	backgroundColor: ['dark', 'dark-hover', 'dark-group-hover', 'dark-even', 'dark-odd', 'hover'],
	borderColor: ['dark', 'dark-disabled', 'dark-focus', 'dark-focus-within', 'hover'],
	textColor: ['dark', 'dark-hover', 'dark-active', 'dark-placeholder', 'hover']
  },
  plugins: [
	  //the import definately works, warning is ignorable
    require('@tailwindcss/ui')({
	  layout: 'sidebar',
	}),
	require("tailwindcss-dark-mode")()
  ]
}
