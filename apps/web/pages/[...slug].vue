<script setup lang="ts">
import { isMediaField, isBlockArray } from '@vibezz/types'

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
        <BlockRenderer v-if="isBlockArray(value)" :blocks="value" />
        <img
          v-else-if="isMediaField(value) && value.mime_type.startsWith('image/')"
          :src="value.url"
          :alt="value.filename"
          class="max-w-full rounded-lg"
        />
        <a
          v-else-if="isMediaField(value)"
          :href="value.url"
          target="_blank"
          class="text-primary hover:underline"
        >
          {{ value.filename }}
        </a>
        <div v-else-if="typeof value === 'string'" class="prose dark:prose-invert max-w-none" v-html="value" />
        <p v-else>{{ value }}</p>
      </div>
    </div>
  </div>
</template>
