<script setup lang="ts">
import type { ContentType, FieldDefinition } from '@vibezz/types'

const route = useRoute()
const { apiFetch } = useApi()

const { data: contentType, refresh } = await useAsyncData(`content-type-${route.params.id}`, () =>
  apiFetch<ContentType>(`/api/admin/content-types/${route.params.id}`),
)

const fieldTypes = [
  { value: 'text', label: 'Text' },
  { value: 'textarea', label: 'Textarea' },
  { value: 'richtext', label: 'Rich Text' },
  { value: 'number', label: 'Number' },
  { value: 'boolean', label: 'Boolean' },
  { value: 'date', label: 'Date' },
  { value: 'select', label: 'Select' },
  { value: 'media', label: 'Media' },
  { value: 'url', label: 'URL' },
  { value: 'email', label: 'Email' },
  { value: 'json', label: 'JSON' },
]

const showAddField = ref(false)
const newFieldName = ref('')
const newFieldSlug = ref('')
const newFieldType = ref('text')
const newFieldRequired = ref(false)

watch(newFieldName, (val) => {
  newFieldSlug.value = val.toLowerCase().replace(/[^a-z0-9]+/g, '_').replace(/(^_|_$)/g, '')
})

async function addField() {
  if (!contentType.value) return

  const fields = [
    ...contentType.value.fields.map(f => ({
      name: f.name,
      slug: f.slug,
      type: f.type,
      required: f.required,
    })),
    {
      name: newFieldName.value,
      slug: newFieldSlug.value,
      type: newFieldType.value,
      required: newFieldRequired.value,
    },
  ]

  await apiFetch(`/api/admin/content-types/${route.params.id}`, {
    method: 'PATCH',
    body: JSON.stringify({ fields }),
  })

  newFieldName.value = ''
  newFieldSlug.value = ''
  newFieldType.value = 'text'
  newFieldRequired.value = false
  showAddField.value = false
  refresh()
}

async function removeField(slug: string) {
  if (!contentType.value) return

  const fields = contentType.value.fields
    .filter(f => f.slug !== slug)
    .map(f => ({
      name: f.name,
      slug: f.slug,
      type: f.type,
      required: f.required,
    }))

  await apiFetch(`/api/admin/content-types/${route.params.id}`, {
    method: 'PATCH',
    body: JSON.stringify({ fields }),
  })

  refresh()
}
</script>

<template>
  <div v-if="contentType">
    <div class="flex items-center gap-4 mb-6">
      <UButton to="/content-types" variant="ghost" icon="i-heroicons-arrow-left" />
      <div>
        <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">{{ contentType.name }}</h2>
        <p class="text-sm text-gray-500">{{ contentType.slug }}</p>
      </div>
    </div>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold">Fields</h3>
          <UButton size="sm" @click="showAddField = true" icon="i-heroicons-plus">
            Add Field
          </UButton>
        </div>
      </template>

      <UTable
        :rows="contentType.fields"
        :columns="[
          { key: 'name', label: 'Name' },
          { key: 'slug', label: 'Slug' },
          { key: 'type', label: 'Type' },
          { key: 'required', label: 'Required' },
          { key: 'actions', label: '' },
        ]"
      >
        <template #required-data="{ row }">
          <UBadge :color="row.required ? 'red' : 'gray'" variant="soft">
            {{ row.required ? 'Required' : 'Optional' }}
          </UBadge>
        </template>
        <template #actions-data="{ row }">
          <UButton
            icon="i-heroicons-trash"
            variant="ghost"
            color="red"
            size="xs"
            @click="removeField(row.slug)"
          />
        </template>
      </UTable>

      <div v-if="!contentType.fields.length" class="text-center py-8 text-gray-500">
        No fields defined yet.
      </div>
    </UCard>

    <UModal v-model="showAddField">
      <UCard>
        <template #header>
          <h3 class="text-lg font-semibold">Add Field</h3>
        </template>
        <form @submit.prevent="addField" class="space-y-4">
          <UFormGroup label="Name">
            <UInput v-model="newFieldName" placeholder="Author Name" required />
          </UFormGroup>
          <UFormGroup label="Slug">
            <UInput v-model="newFieldSlug" placeholder="author_name" required />
          </UFormGroup>
          <UFormGroup label="Type">
            <USelect v-model="newFieldType" :options="fieldTypes" option-attribute="label" value-attribute="value" />
          </UFormGroup>
          <UFormGroup>
            <UCheckbox v-model="newFieldRequired" label="Required" />
          </UFormGroup>
          <div class="flex justify-end gap-2">
            <UButton variant="ghost" @click="showAddField = false">Cancel</UButton>
            <UButton type="submit">Add Field</UButton>
          </div>
        </form>
      </UCard>
    </UModal>
  </div>
</template>
