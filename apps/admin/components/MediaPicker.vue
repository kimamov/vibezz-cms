<script setup lang="ts">
import type { MediaFile } from '@vibezz/types'

const props = defineProps<{
  modelValue?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { apiFetch } = useApi()
const config = useRuntimeConfig()

const showPicker = ref(false)
const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

const { data: files, refresh } = await useAsyncData('media-picker', () =>
  apiFetch<MediaFile[]>('/api/admin/media'),
  { lazy: true },
)

const selectedFile = computed(() => {
  if (!props.modelValue || !files.value) return null
  return files.value.find(f => f.id === props.modelValue) ?? null
})

function mediaUrl(id: string) {
  return `${config.public.apiBase}/api/public/media/${id}`
}

function isImage(mimeType: string) {
  return mimeType.startsWith('image/')
}

function select(file: MediaFile) {
  emit('update:modelValue', file.id)
  showPicker.value = false
}

function clear() {
  emit('update:modelValue', '')
}

function openPicker() {
  refresh()
  showPicker.value = true
}

function triggerUpload() {
  fileInput.value?.click()
}

async function handleUpload(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return

  uploading.value = true
  try {
    let lastUploaded: MediaFile | null = null
    for (const file of input.files) {
      const formData = new FormData()
      formData.append('file', file)
      lastUploaded = await apiFetch<MediaFile>('/api/admin/media', {
        method: 'POST',
        body: formData,
      })
    }
    await refresh()
    if (lastUploaded) {
      select(lastUploaded)
    }
  } finally {
    uploading.value = false
    input.value = ''
  }
}
</script>

<template>
  <div>
    <div v-if="selectedFile" class="flex items-center gap-3 p-3 border border-gray-200 dark:border-gray-700 rounded-lg">
      <div class="w-16 h-16 flex-shrink-0 rounded overflow-hidden bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
        <img
          v-if="isImage(selectedFile.mime_type)"
          :src="mediaUrl(selectedFile.id)"
          :alt="selectedFile.filename"
          class="object-cover w-full h-full"
        />
        <span v-else class="i-heroicons-document text-2xl text-gray-400" />
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{{ selectedFile.filename }}</p>
        <p class="text-xs text-gray-500">{{ selectedFile.mime_type }}</p>
      </div>
      <div class="flex gap-1">
        <UButton size="xs" variant="ghost" @click="openPicker">Change</UButton>
        <UButton size="xs" variant="ghost" color="red" @click="clear">Remove</UButton>
      </div>
    </div>

    <div v-else>
      <UButton variant="outline" @click="openPicker" icon="i-heroicons-photo">
        Select Media
      </UButton>
    </div>

    <UModal v-model="showPicker">
      <UCard>
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Select Media</h3>
            <div class="flex gap-2">
              <UButton size="sm" variant="outline" @click="triggerUpload" :loading="uploading" icon="i-heroicons-arrow-up-tray">
                Upload
              </UButton>
            </div>
            <input
              ref="fileInput"
              type="file"
              class="hidden"
              @change="handleUpload"
            />
          </div>
        </template>

        <div v-if="files?.length" class="grid grid-cols-4 gap-3 max-h-96 overflow-y-auto">
          <button
            v-for="file in files"
            :key="file.id"
            type="button"
            class="rounded-lg border-2 p-1 transition-colors cursor-pointer"
            :class="file.id === modelValue
              ? 'border-primary bg-primary-50 dark:bg-primary-950'
              : 'border-transparent hover:border-gray-300 dark:hover:border-gray-600'"
            @click="select(file)"
          >
            <div class="aspect-square flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded overflow-hidden">
              <img
                v-if="isImage(file.mime_type)"
                :src="mediaUrl(file.id)"
                :alt="file.filename"
                class="object-cover w-full h-full"
              />
              <span v-else class="i-heroicons-document text-3xl text-gray-400" />
            </div>
            <p class="text-xs text-center mt-1 truncate text-gray-700 dark:text-gray-300">{{ file.filename }}</p>
          </button>
        </div>

        <div v-else class="text-center py-8 text-gray-500">
          No media files yet. Upload one to get started.
        </div>
      </UCard>
    </UModal>
  </div>
</template>
