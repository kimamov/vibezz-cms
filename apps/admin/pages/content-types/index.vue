<script setup lang="ts">
import type { ContentType } from '@vibezz/types'

const { apiFetch } = useApi()
const toast = useToast()

const { data: contentTypes, refresh } = await useAsyncData('content-types', () =>
  apiFetch<ContentType[]>('/api/admin/content-types?exclude_slug=page'),
)

const showCreate = ref(false)
const newName = ref('')
const newSlug = ref('')

watch(newName, (val) => {
  newSlug.value = val.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/(^-|-$)/g, '')
})

async function create() {
  try {
    await apiFetch('/api/admin/content-types', {
      method: 'POST',
      body: JSON.stringify({
        name: newName.value,
        slug: newSlug.value,
        fields: [],
      }),
    })
    newName.value = ''
    newSlug.value = ''
    showCreate.value = false
    refresh()
    toast.add({ title: 'Content type created', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to create', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function remove(id: string) {
  if (!confirm('Delete this content type? All its entries will also be removed.')) return
  try {
    await apiFetch(`/api/admin/content-types/${id}`, { method: 'DELETE' })
    refresh()
    toast.add({ title: 'Content type deleted', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to delete', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Content Types</h2>
      <UButton @click="showCreate = true" icon="i-heroicons-plus">
        New Content Type
      </UButton>
    </div>

    <UModal v-model="showCreate">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">Create Content Type</h3>
        </template>
        <form @submit.prevent="create" class="space-y-4">
          <UFormGroup label="Name">
            <UInput v-model="newName" placeholder="Blog Post" required />
          </UFormGroup>
          <UFormGroup label="Slug">
            <UInput v-model="newSlug" placeholder="blog-post" required />
          </UFormGroup>
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" @click="showCreate = false">Cancel</UButton>
            <UButton type="submit">Create</UButton>
          </div>
        </form>
      </UCard>
    </UModal>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <UCard
        v-for="ct in contentTypes"
        :key="ct.id"
      >
        <div class="flex items-start justify-between">
          <div>
            <NuxtLink
              :to="`/content-types/${ct.id}`"
              class="text-lg font-semibold text-gray-900 dark:text-gray-100 hover:text-primary"
            >
              {{ ct.name }}
            </NuxtLink>
            <p class="text-sm text-gray-500 mt-1">{{ ct.slug }}</p>
            <p class="text-xs text-gray-400 mt-2">{{ ct.fields.length }} fields</p>
          </div>
          <UDropdown :items="[[
            { label: 'Edit', icon: 'i-heroicons-pencil', to: `/content-types/${ct.id}` },
            { label: 'Delete', icon: 'i-heroicons-trash', click: () => remove(ct.id) },
          ]]">
            <UButton icon="i-heroicons-ellipsis-vertical" variant="ghost" size="sm" />
          </UDropdown>
        </div>
      </UCard>
    </div>

    <div v-if="!contentTypes?.length" class="text-center py-12 text-gray-500">
      No content types yet. Create one to get started.
    </div>
  </div>
</template>
