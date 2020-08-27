<template>
  <div>
    <div class="pt-8 pb-20 px-4 sm:px-6 lg:pt-12 lg:pb-28 lg:px-8">
      <div class="relative max-w-lg mx-auto lg:max-w-7xl">
        <div>
          <h2
            class="text-3xl leading-9 tracking-tight font-extrabold text-gray-900 dark:text-gray-200 sm:text-4xl sm:leading-10"
          >
            Press
          </h2>
          <div
            class="mt-3 sm:mt-4 lg:grid lg:grid-cols-2 lg:gap-5 lg:items-center"
          >
            <p class="text-xl leading-7 text-gray-500 dark:text-gray-400">
              News about wowmate directly in your inbox.
            </p>
            <form class="mt-6 flex lg:mt-0 lg:justify-end">
              <input
                aria-label="Email address"
                type="email"
                required
                class="appearance-none w-full px-4 py-2 border border-gray-700 text-base leading-6 rounded-md dark:text-gray-100 dark:bg-gray-900 placeholder-gray-500 focus:outline-none focus:shadow-outline-blue focus:border-blue-300 transition duration-150 ease-in-out lg:max-w-xs"
                placeholder="Enter your email"
              />
              <span class="ml-3 flex-shrink-0 inline-flex rounded-md shadow-sm">
                <button
                  type="button"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-base leading-6 font-medium rounded-md text-white bg-red-600 hover:text-gray-800 focus:outline-none focus:border-indigo-700 focus:shadow-outline-indigo active:bg-indigo-700 transition ease-in-out duration-150"
                >
                  Notify me
                </button>
              </span>
            </form>
          </div>
        </div>
        <div
          class="mt-6 grid gap-16 border-t-2 border-gray-100 dark:border-gray-600 pt-10 lg:grid-cols-2 lg:col-gap-5 lg:row-gap-12"
        >
          <div v-for="article in articles" :key="article.title">
            <p class="text-sm leading-5 text-gray-500 dark:text-gray-400">
              <time datetime="2020-03-16">Mar 16, 2020</time>
            </p>
            <NuxtLink
              :to="{ name: 'news-slug', params: { slug: article.slug } }"
              class="block"
            >
              <h3
                class="mt-2 text-xl leading-7 font-semibold text-gray-900 dark:text-gray-200"
              >
                {{ article.title1 }} {{ article.title2 }}
              </h3>
              <p
                class="mt-3 text-base leading-6 text-gray-500 dark:text-gray-400"
              >
                {{ article.description }}
              </p>
            </NuxtLink>
            <div class="mt-3">
              <NuxtLink
                :to="{ name: 'news-slug', params: { slug: article.slug } }"
                class="text-base leading-6 font-semibold text-red-600 hover:text-gray-800 dark-hover:text-gray-200 transition ease-in-out duration-150"
              >
                Read full story
              </NuxtLink>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  async asyncData({ $content, params }) {
    const articles = await $content('articles', params.slug)
      .only(['title1', 'title2', 'description', 'slug'])
      .sortBy('createdAt', 'asc')
      .fetch()

    return {
      articles
    }
  }
}
</script>

<style></style>
