<template>
  <div>
    <!-- TODO: add nice patterns -->
    <main
      class="my-8 mx-auto max-w-screen-xl px-4 sm:mt-12 sm:px-6 md:mt-20 xl:mt-24"
    >
      <div class="lg:grid lg:grid-cols-12 lg:gap-8">
        <div
          class="sm:text-center md:max-w-2xl md:mx-auto lg:col-span-6 lg:text-left"
        >
          <div
            class="text-sm font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400 sm:text-base lg:text-sm xl:text-base"
          >
            Coming soon
          </div>
          <h2
            class="mt-1 text-4xl tracking-tight leading-10 font-extrabold text-gray-900 dark:text-gray-200 sm:leading-none sm:text-6xl lg:text-5xl xl:text-6xl"
          >
            Combatlogs for World of Warcraft
            <br class="hidden md:inline" />
            <span class="text-red-600">made easy</span>
          </h2>
          <p
            class="mt-3 text-base text-gray-500 dark:text-gray-400 sm:mt-5 sm:text-xl lg:text-lg xl:text-xl"
          >
            Nothing on this site actually works yet. I'm aiming to have
            something usable for shadowlands release
          </p>
          <div
            class="mt-5 sm:max-w-lg sm:mx-auto sm:text-center lg:text-left lg:mx-0"
          >
            <p class="text-base font-medium text-gray-900 dark:text-gray-200">
              Sign up to get notified when it’s ready.
            </p>
            <form action="#" method="POST" class="mt-3 sm:flex">
              <input
                aria-label="Email"
                class="appearance-none block w-full px-3 py-3 border border-gray-300 dark:border-gray-700 dark:bg-gray-900 text-base leading-6 rounded-md placeholder-gray-500 shadow-sm focus:outline-none focus:placeholder-gray-400 focus:shadow-outline-blue focus:border-blue-300 transition duration-150 ease-in-out sm:flex-1"
                placeholder="doesn't work either"
              />
              <button
                type="submit"
                class="mt-3 w-full px-6 py-3 border border-transparent text-base leading-6 font-medium rounded-md text-white bg-red-600 shadow-sm hover:text-gray-800 focus:outline-none focus:border-gray-900 focus:shadow-outline-gray active:bg-gray-900 transition duration-150 ease-in-out sm:mt-0 sm:ml-3 sm:flex-shrink-0 sm:inline-flex sm:items-center sm:w-auto"
              >
                don't Notify me
              </button>
            </form>
            <p class="mt-3 text-sm leading-5 text-gray-500 dark:text-gray-400">
              We care about the protection of your data. Read our
              <a
                href="#"
                class="font-medium text-gray-900 dark:text-gray-200 underline"
                >Privacy Policy</a
              >.
            </p>
          </div>
        </div>
        <div
          class="mt-12 relative sm:max-w-lg sm:mx-auto lg:mt-0 lg:max-w-none lg:mx-0 lg:col-span-6 lg:flex lg:items-center"
        >
          <div class="relative mx-auto w-full rounded-lg lg:max-w-md">
            <button
              type="button"
              class="relative block w-full rounded-lg overflow-hidden focus:outline-none focus:shadow-outline"
            >
              <img
                class="w-full"
                src="undraw_dark_analytics_7axy.svg"
                alt="graph"
              />
              <div
                class="absolute inset-0 w-full h-full flex items-center justify-center"
              ></div>
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
// import Logo from '~/components/Logo.vue'

export default {}
</script>

<style>
/* Sample `apply` at-rules with Tailwind CSS
.container {
  @apply min-h-screen flex justify-center items-center text-center mx-auto;
}
*/
</style>
