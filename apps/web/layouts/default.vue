<script setup lang="ts">
import type { NavigationItem } from '@vibezz/types'

const { fetchNavigation } = useCms()
const { data: navigation } = await useAsyncData('navigation', fetchNavigation)
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <header class="border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900">
      <div class="max-w-6xl mx-auto px-4 py-4 flex items-center justify-between">
        <NuxtLink to="/" class="text-xl font-bold">Vibezz</NuxtLink>
        <nav class="flex gap-6">
          <NuxtLink
            v-for="item in navigation"
            :key="item.id"
            :to="item.path"
            class="text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors"
          >
            {{ item.title }}
          </NuxtLink>
        </nav>
      </div>
    </header>

    <main class="flex-1">
      <div class="max-w-6xl mx-auto px-4 py-8">
        <slot />
      </div>
    </main>

    <footer class="border-t border-gray-200 dark:border-gray-800 py-6">
      <div class="max-w-6xl mx-auto px-4 text-center text-sm text-gray-500">
        Powered by Vibezz CMS
      </div>
    </footer>
  </div>
</template>
