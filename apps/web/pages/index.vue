<script setup lang="ts">
const { fetchPage } = useCms()
const { data: page } = await useAsyncData('home', () => fetchPage('/'))
</script>

<template>
  <div>
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
    <div v-else class="text-center py-20">
      <h1 class="text-4xl font-bold mb-4">Welcome to Vibezz</h1>
      <p class="text-gray-500 text-lg">
        Your content will appear here once you publish a page with the path "/".
      </p>
    </div>
  </div>
</template>
