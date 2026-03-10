<script setup lang="ts">
const route = useRoute()
const { fetchPage } = useCms()

const path = computed(() => '/' + (Array.isArray(route.params.slug) ? route.params.slug.join('/') : route.params.slug))

const { data: page, error } = await useAsyncData(
  `page-${path.value}`,
  () => fetchPage(path.value),
)

if (error.value || !page.value) {
  throw createError({
    statusCode: 404,
    statusMessage: 'Page Not Found',
  })
}

useHead({
  title: page.value?.title || 'Page',
})
</script>

<template>
  <div v-if="page">
    <h1 class="text-4xl font-bold mb-4">{{ page.title }}</h1>
    <div v-if="page.fields">
      <div
        v-for="(value, key) in page.fields"
        :key="String(key)"
        class="mb-4"
      >
        <div v-if="typeof value === 'string'" class="prose dark:prose-invert max-w-none" v-html="value" />
        <p v-else>{{ value }}</p>
      </div>
    </div>
  </div>
</template>
