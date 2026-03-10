<script setup lang="ts">
import type { Entry, ContentType } from '@vibezz/types'

const route = useRoute()
const { apiFetch } = useApi()

const { data: entry, refresh } = await useAsyncData(`entry-${route.params.id}`, () =>
  apiFetch<Entry>(`/api/admin/entries/${route.params.id}`),
)

const { data: contentType } = await useAsyncData(`entry-ct-${route.params.id}`, async () => {
  if (!entry.value) return null
  return apiFetch<ContentType>(`/api/admin/content-types/${entry.value.content_type_id}`)
})

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
    await apiFetch(`/api/admin/entries/${route.params.id}`, {
      method: 'PATCH',
      body: JSON.stringify({
        title: title.value,
        slug: slug.value,
        fields: fields.value,
      }),
    })
    refresh()
  } finally {
    saving.value = false
  }
}

async function publish() {
  await apiFetch(`/api/admin/entries/${route.params.id}/publish`, { method: 'POST' })
  refresh()
}

async function unpublish() {
  await apiFetch(`/api/admin/entries/${route.params.id}/unpublish`, { method: 'POST' })
  refresh()
}
</script>

<template>
  <div v-if="entry">
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-4">
        <UButton to="/entries" variant="ghost" icon="i-heroicons-arrow-left" />
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
          <UFormGroup label="Slug">
            <UInput v-model="slug" required />
          </UFormGroup>
        </div>
      </UCard>

      <UCard v-if="contentType?.fields.length">
        <template #header>
          <h3 class="text-lg font-semibold">Fields</h3>
        </template>
        <div class="space-y-4">
          <UFormGroup
            v-for="field in contentType.fields"
            :key="field.slug"
            :label="field.name"
            :required="field.required"
          >
            <UInput
              v-if="field.type === 'text' || field.type === 'url' || field.type === 'email'"
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
