<script setup lang="ts">
import type { Entry, ContentType } from '@vibezz/types'

const { apiFetch } = useApi()
const toast = useToast()

const { data: entries, refresh } = await useAsyncData('entries', () =>
  apiFetch<Entry[]>('/api/admin/entries'),
)

const { data: contentTypes } = await useAsyncData('entry-content-types', () =>
  apiFetch<ContentType[]>('/api/admin/content-types?exclude_slug=page'),
)

const showCreate = ref(false)
const newTitle = ref('')
const newSlug = ref('')
const newContentTypeId = ref('')

watch(newTitle, (val) => {
  newSlug.value = val.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
})

const contentTypeOptions = computed(() =>
  (contentTypes.value || []).map(ct => ({ value: ct.id, label: ct.name })),
)

async function create() {
  try {
    await apiFetch('/api/admin/entries', {
      method: 'POST',
      body: JSON.stringify({
        content_type_id: newContentTypeId.value,
        title: newTitle.value,
        slug: newSlug.value,
        fields: {},
      }),
    })
    newTitle.value = ''
    newSlug.value = ''
    newContentTypeId.value = ''
    showCreate.value = false
    refresh()
    toast.add({ title: 'Entry created', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to create entry', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function publish(id: string) {
  try {
    await apiFetch(`/api/admin/entries/${id}/publish`, { method: 'POST' })
    refresh()
    toast.add({ title: 'Entry published', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to publish', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function unpublish(id: string) {
  try {
    await apiFetch(`/api/admin/entries/${id}/unpublish`, { method: 'POST' })
    refresh()
    toast.add({ title: 'Entry unpublished', icon: 'i-heroicons-check-circle', color: 'yellow' })
  } catch (e: any) {
    toast.add({ title: 'Failed to unpublish', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function remove(id: string) {
  if (!confirm('Delete this entry?')) return
  try {
    await apiFetch(`/api/admin/entries/${id}`, { method: 'DELETE' })
    refresh()
    toast.add({ title: 'Entry deleted', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to delete', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Entries</h2>
      <UButton @click="showCreate = true" icon="i-heroicons-plus">
        New Entry
      </UButton>
    </div>

    <UModal v-model="showCreate">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">Create Entry</h3>
        </template>
        <form @submit.prevent="create" class="space-y-4">
          <UFormGroup label="Content Type">
            <USelect
              v-model="newContentTypeId"
              :options="contentTypeOptions"
              option-attribute="label"
              value-attribute="value"
              placeholder="Select a content type"
              required
            />
          </UFormGroup>
          <UFormGroup label="Title">
            <UInput v-model="newTitle" placeholder="My First Post" required />
          </UFormGroup>
          <UFormGroup label="Slug" hint="Leave empty for root page (/)">
            <UInput v-model="newSlug" placeholder="my-first-post" />
          </UFormGroup>
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" @click="showCreate = false">Cancel</UButton>
            <UButton type="submit">Create</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <UCard>
      <UTable
        :rows="entries || []"
        :columns="[
          { key: 'title', label: 'Title' },
          { key: 'path', label: 'Path' },
          { key: 'status', label: 'Status' },
          { key: 'updated_at', label: 'Updated' },
          { key: 'actions', label: '' },
        ]"
      >
        <template #title-data="{ row }">
          <NuxtLink
            :to="`/entries/${row.id}`"
            class="text-primary hover:underline font-medium"
          >
            {{ row.title }}
          </NuxtLink>
        </template>
        <template #status-data="{ row }">
          <UBadge
            :color="row.status === 'published' ? 'green' : 'yellow'"
            variant="soft"
          >
            {{ row.status }}
          </UBadge>
        </template>
        <template #updated_at-data="{ row }">
          {{ new Date(row.updated_at).toLocaleDateString() }}
        </template>
        <template #actions-data="{ row }">
          <div class="flex gap-1">
            <UButton
              v-if="row.status === 'draft'"
              size="xs"
              variant="ghost"
              color="green"
              @click="publish(row.id)"
            >
              Publish
            </UButton>
            <UButton
              v-else
              size="xs"
              variant="ghost"
              color="yellow"
              @click="unpublish(row.id)"
            >
              Unpublish
            </UButton>
            <UButton
              icon="i-heroicons-trash"
              variant="ghost"
              color="red"
              size="xs"
              @click="remove(row.id)"
            />
          </div>
        </template>
      </UTable>
    </UCard>

    <div v-if="!entries?.length" class="text-center py-12 text-gray-500">
      No entries yet. Create one to get started.
    </div>
  </div>
</template>
