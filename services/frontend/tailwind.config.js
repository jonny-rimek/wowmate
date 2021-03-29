/*
** TailwindCSS Configuration File
**
** Docs: https://tailwindcss.com/docs/configuration
** Default: https://github.com/tailwindcss/tailwindcss/blob/master/stubs/defaultConfig.stub.js
*/
const colors = require('tailwindcss/colors') //ignore waring

module.exports = {
  theme: {
    //colors, //this would add all new colors, but it changes the old gray I use as bg for dark mode
    /* to import a single new color
    colors: {
      cyan: colors.cyan
    },
     */
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
  plugins: [
	  //the import definitely works, warning is ignorable
    require('@tailwindcss/ui')({
      layout: 'sidebar',
    }),
    require('@tailwindcss/forms'),
  ]
}
