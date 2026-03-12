<script setup lang="ts">
import type { ContentBlock, BlockStyle } from '@vibezz/types'
import { isMediaField } from '@vibezz/types'

const props = defineProps<{
  blocks: ContentBlock[]
  level?: number
}>()

const level = computed(() => props.level ?? 0)

function headingTag(blockLevel: number): string {
  const tags: Record<number, string> = { 1: 'h1', 2: 'h2', 3: 'h3', 4: 'h4', 5: 'h5', 6: 'h6' }
  return tags[blockLevel] || 'h2'
}

function getStyleClasses(style?: BlockStyle): string {
  if (!style) return ''
  return style.className || ''
}

function getStyleObject(style?: BlockStyle): Record<string, string> {
  if (!style?.inlineStyle) return {}
  return style.inlineStyle
}

function getResponsiveClasses(style?: BlockStyle): string {
  if (!style?.responsive) return ''
  const classes: string[] = []
  
  if (style.responsive.mobile?.className) {
    classes.push(style.responsive.mobile.className)
  }
  if (style.responsive.tablet?.className) {
    classes.push(`md:${style.responsive.tablet.className}`)
  }
  if (style.responsive.desktop?.className) {
    classes.push(`lg:${style.responsive.desktop.className}`)
  }
  
  return classes.join(' ')
}

// Generate Tailwind grid column classes for responsive grids
function getGridColumnClasses(data: Record<string, unknown>): string {
  const mobileCols = (data.mobileColumns as number) || 1
  const tabletCols = (data.tabletColumns as number) || 2
  const desktopCols = (data.desktopColumns as number) || (data.columns as number) || 3
  
  // Map column counts to Tailwind classes
  const colClasses: Record<number, string> = {
    1: 'grid-cols-1',
    2: 'grid-cols-2',
    3: 'grid-cols-3',
    4: 'grid-cols-4',
    5: 'grid-cols-5',
    6: 'grid-cols-6',
    7: 'grid-cols-7',
    8: 'grid-cols-8',
    9: 'grid-cols-9',
    10: 'grid-cols-10',
    11: 'grid-cols-11',
    12: 'grid-cols-12',
  }
  
  const classes: string[] = []
  
  // Mobile (default)
  classes.push(colClasses[mobileCols] || 'grid-cols-1')
  
  // Tablet (md breakpoint)
  classes.push(`md:${colClasses[tabletCols] || 'grid-cols-2'}`)
  
  // Desktop (lg breakpoint)
  classes.push(`lg:${colClasses[desktopCols] || 'grid-cols-3'}`)
  
  return classes.join(' ')
}
</script>

<template>
    <template v-for="block in blocks" :key="block.id">
      <!-- Heading -->
      <component
        v-if="block.type === 'heading'"
        :is="headingTag(block.data.level as number)"
        class="font-bold"
        :class="[
          {
            'text-4xl': block.data.level === 1,
            'text-3xl': block.data.level === 2,
            'text-2xl': block.data.level === 3,
            'text-xl': block.data.level === 4,
            'text-lg': block.data.level === 5,
            'text-base': block.data.level === 6,
          },
          getStyleClasses(block.style),
          getResponsiveClasses(block.style),
        ]"
        :style="getStyleObject(block.style)"
      >
        {{ block.data.text }}
      </component>

      <!-- Text -->
      <div
        v-else-if="block.type === 'text'"
        class="prose dark:prose-invert max-w-none"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
        v-html="block.data.content"
      />

      <!-- Image -->
      <figure
        v-else-if="block.type === 'image' && block.data.media"
        class="w-full"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
      >
        <img
          v-if="isMediaField(block.data.media) && block.data.media.mime_type.startsWith('image/')"
          :src="block.data.media.url"
          :alt="(block.data.alt as string) || block.data.media.filename"
          class="w-full rounded-lg"
          :class="{
            'aspect-video object-cover': block.data.aspectRatio === '16:9',
            'aspect-[4/3] object-cover': block.data.aspectRatio === '4:3',
            'aspect-square object-cover': block.data.aspectRatio === '1:1',
            'aspect-[3/2] object-cover': block.data.aspectRatio === '3:2',
          }"
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
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
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
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
      ><code :class="`language-${block.data.language || 'plaintext'}`">{{ block.data.code }}</code></pre>

      <!-- Divider -->
      <hr
        v-else-if="block.type === 'divider'"
        class="border-gray-200 dark:border-gray-700"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
      />

      <!-- Container Block -->
      <div
        v-else-if="block.type === 'container'"
        class="container-block"
        :class="[
          {
            'flex': block.data.layout === 'flex',
            'grid': block.data.layout === 'grid',
            'flex flex-row': block.data.layout === 'row',
            'flex flex-col': block.data.layout === 'column' || block.data.layout === 'stack',
          },
          getStyleClasses(block.style),
          getResponsiveClasses(block.style),
        ]"
        :style="{
          gap: block.data.gap || '16px',
          padding: block.data.padding || '0px',
          background: block.data.background || 'transparent',
          maxWidth: block.data.maxWidth || 'none',
          ...getStyleObject(block.style),
        }"
      >
        <BlockRenderer
          v-if="block.children?.length"
          :blocks="block.children"
          :level="level + 1"
        />
      </div>

      <!-- Grid Block -->
      <div
        v-else-if="block.type === 'grid'"
        class="grid"
        :class="[
          getGridColumnClasses(block.data),
          getStyleClasses(block.style),
          getResponsiveClasses(block.style),
        ]"
        :style="{
          gap: `${block.data.rowGap || '16px'} ${block.data.columnGap || '16px'}`,
          ...getStyleObject(block.style),
        }"
      >
        <BlockRenderer
          v-if="block.children?.length"
          :blocks="block.children"
          :level="level + 1"
        />
      </div>

      <!-- Section Block -->
      <section
        v-else-if="block.type === 'section'"
        class="section-block"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="{
          backgroundColor: block.data.background || 'transparent',
          backgroundImage: block.data.backgroundImage ? `url(${block.data.backgroundImage})` : 'none',
          backgroundSize: block.data.backgroundSize || 'cover',
          padding: `${block.data.paddingY || '64px'} ${block.data.paddingX || '0px'}`,
          minHeight: block.data.minHeight || 'auto',
          ...getStyleObject(block.style),
        }"
      >
        <BlockRenderer
          v-if="block.children?.length"
          :blocks="block.children"
          :level="level + 1"
        />
      </section>

      <!-- Column Block -->
      <div
        v-else-if="block.type === 'column'"
        class="column-block"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="{
          flex: block.data.width ? `0 0 ${block.data.width}` : '1',
          padding: block.data.padding || '0px',
          ...getStyleObject(block.style),
        }"
      >
        <BlockRenderer
          v-if="block.children?.length"
          :blocks="block.children"
          :level="level + 1"
        />
      </div>

      <!-- News List (plugin block; data.items from enricher) -->
      <section
        v-else-if="block.type === 'news_list' && block.data.items"
        class="space-y-4"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
      >
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
      <div
        v-else
        class="text-sm text-gray-500"
        :class="[getStyleClasses(block.style), getResponsiveClasses(block.style)]"
        :style="getStyleObject(block.style)"
      >
        Block type "{{ block.type }}"
      </div>
    </template>
</template>