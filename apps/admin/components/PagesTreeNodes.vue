<script setup lang="ts">
interface PageNode {
  id: string
  title: string
  slug: string
  path: string
  status: string
  updated_at: string
  children?: PageNode[]
}

defineProps<{
  nodes: PageNode[]
}>()

const emit = defineEmits<{
  publish: [id: string]
  unpublish: [id: string]
  delete: [id: string]
}>()
</script>

<template>
  <div class="divide-y divide-gray-100 dark:divide-gray-800">
    <div
      v-for="node in nodes"
      :key="node.id"
      class="flex flex-col"
    >
      <div class="flex items-center justify-between gap-4 py-2 px-2 hover:bg-gray-50 dark:hover:bg-gray-800/50 rounded">
        <div class="flex items-center gap-3 min-w-0 flex-1">
          <NuxtLink
            :to="`/pages/${node.id}`"
            class="text-primary hover:underline font-medium truncate"
          >
            {{ node.title || node.path || 'Untitled' }}
          </NuxtLink>
          <span class="text-sm text-gray-500 truncate">{{ node.path }}</span>
          <UBadge
            :color="node.status === 'published' ? 'green' : 'yellow'"
            variant="soft"
            class="flex-shrink-0"
          >
            {{ node.status }}
          </UBadge>
        </div>
        <div class="flex gap-1 flex-shrink-0">
          <UButton
            v-if="node.status === 'draft'"
            size="xs"
            variant="ghost"
            color="green"
            @click="emit('publish', node.id)"
          >
            Publish
          </UButton>
          <UButton
            v-else
            size="xs"
            variant="ghost"
            color="yellow"
            @click="emit('unpublish', node.id)"
          >
            Unpublish
          </UButton>
          <UButton
            icon="i-heroicons-trash"
            variant="ghost"
            color="red"
            size="xs"
            @click="emit('delete', node.id)"
          />
        </div>
      </div>
      <div v-if="node.children?.length" class="pl-6 border-l-2 border-gray-200 dark:border-gray-700 ml-2">
        <PagesTreeNodes
          :nodes="node.children"
          @publish="(id) => emit('publish', id)"
          @unpublish="(id) => emit('unpublish', id)"
          @delete="(id) => emit('delete', id)"
        />
      </div>
    </div>
  </div>
</template>
