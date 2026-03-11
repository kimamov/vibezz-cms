<script setup lang="ts">
import type { Entry } from '@vibezz/types'

interface PageNode {
  id: string
  title: string
  slug: string
  path: string
  status: string
  updated_at: string
  children?: PageNode[]
}

const { apiFetch } = useApi()
const toast = useToast()

const { data: tree, refresh } = await useAsyncData('pages-tree', () =>
  apiFetch<PageNode[]>('/api/admin/pages'),
)

const showCreate = ref(false)
const newTitle = ref('')
const newSlug = ref('')
const newParentId = ref('')

watch(newTitle, (val) => {
  newSlug.value = val.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
})

const parentOptions = computed(() => {
  const options: { value: string; label: string }[] = [{ value: '', label: '— Root —' }]
  function collect(nodes: PageNode[] | undefined, prefix: string) {
    if (!nodes) return
    for (const n of nodes) {
      options.push({ value: n.id, label: `${prefix}${n.title}` })
      collect(n.children, `${prefix}${n.title} / `)
    }
  }
  collect(tree.value, '')
  return options
})

async function create() {
  try {
    await apiFetch('/api/admin/pages', {
      method: 'POST',
      body: JSON.stringify({
        title: newTitle.value,
        slug: newSlug.value,
        parent_id: newParentId.value ? newParentId.value : null,
        fields: { blocks: [] },
      }),
    })
    newTitle.value = ''
    newSlug.value = ''
    newParentId.value = ''
    showCreate.value = false
    refresh()
    toast.add({ title: 'Page created', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to create page', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function publish(id: string) {
  try {
    await apiFetch(`/api/admin/pages/${id}/publish`, { method: 'POST' })
    refresh()
    toast.add({ title: 'Page published', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to publish', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function unpublish(id: string) {
  try {
    await apiFetch(`/api/admin/pages/${id}/unpublish`, { method: 'POST' })
    refresh()
    toast.add({ title: 'Page unpublished', icon: 'i-heroicons-check-circle', color: 'yellow' })
  } catch (e: any) {
    toast.add({ title: 'Failed to unpublish', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function remove(id: string) {
  if (!confirm('Delete this page? Child pages will need a new parent or be deleted.')) return
  try {
    await apiFetch(`/api/admin/pages/${id}`, { method: 'DELETE' })
    refresh()
    toast.add({ title: 'Page deleted', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to delete', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Pages</h2>
      <UButton @click="showCreate = true" icon="i-heroicons-plus">
        New Page
      </UButton>
    </div>

    <UModal v-model="showCreate">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">Create Page</h3>
        </template>
        <form @submit.prevent="create" class="space-y-4">
          <UFormGroup label="Parent Page">
            <USelect
              v-model="newParentId"
              :options="parentOptions"
              option-attribute="label"
              value-attribute="value"
              placeholder="Root (no parent)"
            />
          </UFormGroup>
          <UFormGroup label="Title">
            <UInput v-model="newTitle" placeholder="About" required />
          </UFormGroup>
          <UFormGroup label="Slug" hint="Leave empty for root page (/)">
            <UInput v-model="newSlug" placeholder="about" />
          </UFormGroup>
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" @click="showCreate = false">Cancel</UButton>
            <UButton type="submit">Create</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <UCard>
      <div class="space-y-0">
        <template v-if="tree?.length">
          <PagesTreeNodes
            :nodes="tree"
            @publish="publish"
            @unpublish="unpublish"
            @delete="remove"
          />
        </template>
        <div v-else class="text-center py-12 text-gray-500">
          No pages yet. Create one to get started.
        </div>
      </div>
    </UCard>
  </div>
</template>
