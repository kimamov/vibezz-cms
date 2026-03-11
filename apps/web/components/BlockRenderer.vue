<script setup lang="ts">
import type { ContentBlock } from '@vibezz/types'
import { isMediaField } from '@vibezz/types'

defineProps<{
  blocks: ContentBlock[]
}>()

function headingTag(level: number): string {
  const tags: Record<number, string> = { 1: 'h1', 2: 'h2', 3: 'h3', 4: 'h4', 5: 'h5', 6: 'h6' }
  return tags[level] || 'h2'
}
</script>

<template>
  <div class="space-y-6">
    <template v-for="block in blocks" :key="block.id">
      <!-- Heading -->
      <component
        v-if="block.type === 'heading'"
        :is="headingTag(block.data.level as number)"
        class="font-bold"
        :class="{
          'text-4xl': block.data.level === 1,
          'text-3xl': block.data.level === 2,
          'text-2xl': block.data.level === 3,
          'text-xl': block.data.level === 4,
          'text-lg': block.data.level === 5,
          'text-base': block.data.level === 6,
        }"
      >
        {{ block.data.text }}
      </component>

      <!-- Text -->
      <div
        v-else-if="block.type === 'text'"
        class="prose dark:prose-invert max-w-none"
        v-html="block.data.content"
      />

      <!-- Image -->
      <figure v-else-if="block.type === 'image' && block.data.media" class="w-full">
        <img
          v-if="isMediaField(block.data.media) && block.data.media.mime_type.startsWith('image/')"
          :src="block.data.media.url"
          :alt="(block.data.alt as string) || block.data.media.filename"
          class="w-full rounded-lg"
        />
        <figcaption
          v-if="block.data.caption"
          class="text-sm text-gray-500 mt-2 text-center"
        >
          {{ block.data.caption }}
        </figcaption>
      </figure>

      <!-- Quote -->
      <blockquote
        v-else-if="block.type === 'quote'"
        class="border-l-4 border-gray-300 dark:border-gray-600 pl-4 italic text-gray-700 dark:text-gray-300"
      >
        <p>{{ block.data.text }}</p>
        <footer v-if="block.data.attribution" class="text-sm text-gray-500 mt-2 not-italic">
          &mdash; {{ block.data.attribution }}
        </footer>
      </blockquote>

      <!-- Code -->
      <pre
        v-else-if="block.type === 'code'"
        class="bg-gray-900 text-gray-100 rounded-lg p-4 overflow-x-auto text-sm"
      ><code>{{ block.data.code }}</code></pre>

      <!-- Divider -->
      <hr
        v-else-if="block.type === 'divider'"
        class="border-gray-200 dark:border-gray-700"
      />

      <!-- News List (plugin block; data.items from enricher) -->
      <section v-else-if="block.type === 'news_list' && block.data.items" class="space-y-4">
        <h2 class="text-2xl font-bold">News</h2>
        <ul class="space-y-3">
          <li
            v-for="item in block.data.items"
            :key="item.id"
            class="border-b border-gray-200 dark:border-gray-700 pb-3 last:border-0"
          >
            <NuxtLink
              :to="item.path || `/news/${item.slug}`"
              class="font-medium text-primary hover:underline"
            >
              {{ item.headline ?? item.title }}
            </NuxtLink>
            <p v-if="item.excerpt" class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              {{ item.excerpt }}
            </p>
          </li>
        </ul>
      </section>

      <!-- Generic plugin block fallback -->
      <div v-else class="text-sm text-gray-500">
        Block type "{{ block.type }}"
      </div>
    </template>
  </div>
</template>
