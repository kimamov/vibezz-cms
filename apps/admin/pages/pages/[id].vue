<script setup lang="ts">
import type { Entry, ContentType } from '@vibezz/types'

const route = useRoute()
const { apiFetch } = useApi()
const toast = useToast()

interface PageResponse {
  page: Entry
  content_type: ContentType
}

const { data: response, refresh } = await useAsyncData(`page-${route.params.id}`, () =>
  apiFetch<PageResponse>(`/api/admin/pages/${route.params.id}`),
)

const entry = computed(() => response.value?.page ?? null)
const contentType = computed(() => response.value?.content_type ?? null)

const title = ref('')
const slug = ref('')
const fields = ref<Record<string, unknown>>({})
const saving = ref(false)

watch(entry, (val) => {
  if (val) {
    title.value = val.title
    slug.value = val.slug
    fields.value = { ...val.fields }
  }
}, { immediate: true })

async function save() {
  saving.value = true
  try {
    await apiFetch(`/api/admin/pages/${route.params.id}`, {
      method: 'PATCH',
      body: JSON.stringify({
        title: title.value,
        slug: slug.value,
        fields: fields.value,
      }),
    })
    refresh()
    toast.add({ title: 'Changes saved', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to save', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  } finally {
    saving.value = false
  }
}

async function publish() {
  try {
    await apiFetch(`/api/admin/pages/${route.params.id}/publish`, { method: 'POST' })
    refresh()
    toast.add({ title: 'Page published', icon: 'i-heroicons-check-circle', color: 'green' })
  } catch (e: any) {
    toast.add({ title: 'Failed to publish', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}

async function unpublish() {
  try {
    await apiFetch(`/api/admin/pages/${route.params.id}/unpublish`, { method: 'POST' })
    refresh()
    toast.add({ title: 'Page unpublished', icon: 'i-heroicons-check-circle', color: 'yellow' })
  } catch (e: any) {
    toast.add({ title: 'Failed to unpublish', description: e.message, icon: 'i-heroicons-x-circle', color: 'red' })
  }
}
</script>

<template>
  <div v-if="entry">
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-4">
        <UButton to="/pages" variant="ghost" icon="i-heroicons-arrow-left" />
        <div>
          <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ entry.title }}</h2>
          <p class="text-sm text-gray-500">{{ entry.path }}</p>
        </div>
        <UBadge
          :color="entry.status === 'published' ? 'green' : 'yellow'"
          variant="soft"
        >
          {{ entry.status }}
        </UBadge>
      </div>
      <div class="flex gap-2">
        <UButton
          v-if="entry.status === 'draft'"
          color="green"
          @click="publish"
        >
          Publish
        </UButton>
        <UButton
          v-else
          color="yellow"
          variant="outline"
          @click="unpublish"
        >
          Unpublish
        </UButton>
      </div>
    </div>

    <form @submit.prevent="save" class="space-y-6">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">General</h3>
        </template>
        <div class="space-y-4">
          <UFormGroup label="Title">
            <UInput v-model="title" required />
          </UFormGroup>
          <UFormGroup label="Slug" hint="Leave empty for root page (/)">
            <UInput v-model="slug" />
          </UFormGroup>
        </div>
      </UCard>

      <UCard v-if="contentType?.fields?.length">
        <template #header>
          <h3 class="text-lg font-semibold">Content</h3>
        </template>
        <div class="space-y-4">
          <UFormGroup
            v-for="field in contentType.fields"
            :key="field.slug"
            :label="field.name"
            :required="field.required"
          >
            <MediaPicker
              v-if="field.type === 'media'"
              v-model="(fields[field.slug] as string)"
            />
            <BlockEditor
              v-else-if="field.type === 'blocks'"
              v-model="(fields[field.slug] as any[])"
            />
            <UInput
              v-else-if="field.type === 'text' || field.type === 'url' || field.type === 'email'"
              v-model="(fields[field.slug] as string)"
              :type="field.type === 'email' ? 'email' : field.type === 'url' ? 'url' : 'text'"
              :required="field.required"
            />
            <UTextarea
              v-else-if="field.type === 'textarea' || field.type === 'richtext'"
              v-model="(fields[field.slug] as string)"
              :required="field.required"
              :rows="field.type === 'richtext' ? 8 : 4"
            />
            <UInput
              v-else-if="field.type === 'number'"
              v-model="(fields[field.slug] as string)"
              type="number"
              :required="field.required"
            />
            <UInput
              v-else-if="field.type === 'date'"
              v-model="(fields[field.slug] as string)"
              type="date"
              :required="field.required"
            />
            <UCheckbox
              v-else-if="field.type === 'boolean'"
              v-model="(fields[field.slug] as boolean)"
              :label="field.name"
            />
            <UTextarea
              v-else-if="field.type === 'json'"
              v-model="(fields[field.slug] as string)"
              :rows="6"
              :required="field.required"
              placeholder="{}"
            />
            <UInput
              v-else
              v-model="(fields[field.slug] as string)"
              :required="field.required"
            />
          </UFormGroup>
        </div>
      </UCard>

      <div class="flex justify-end">
        <UButton type="submit" :loading="saving">
          Save Changes
        </UButton>
      </div>
    </form>
  </div>
</template>
