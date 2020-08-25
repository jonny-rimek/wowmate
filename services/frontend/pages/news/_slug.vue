<template>
  <article class="max-w-5xl mx-auto">
    <!-- This component requires Tailwind CSS >= 1.5.1 and @tailwindcss/ui >= 0.4.0 -->
    <div class="relative py-16 overflow-hidden">
      <div class="hidden lg:block lg:absolute lg:inset-y-0 lg:h-full lg:w-full">
        <!-- <div class="relative h-full text-lg max-w-prose mx-auto"></div> -->
      </div>
      <div class="relative px-4 sm:px-6 lg:px-8">
        <div class="text-lg max-w-prose mx-auto mb-6">
          <p
            class="text-base text-center leading-6 text-red-600 font-semibold tracking-wide uppercase"
          >
            WOW mate
          </p>
          <h1
            class="mt-2 mb-8 text-3xl text-center leading-8 font-extrabold tracking-tight text-gray-900 dark:text-gray-200 sm:text-4xl sm:leading-10"
          >
            that's amazing
          </h1>
          <p class="text-xl text-gray-500 leading-8"></p>
          <p>
            {{ article.description }}
          </p>
        </div>
        <div class="prose prose-lg text-gray-500 mx-auto">
          <p>
            <nuxt-content :document="article" />
          </p>
        </div>
      </div>
    </div>
  </article>
</template>
<script>
export default {
  async asyncData({ $content, params }) {
    const article = await $content('articles', params.slug).fetch()

    return { article }
  }
}
</script>
