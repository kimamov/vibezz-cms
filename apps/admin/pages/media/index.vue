<script setup lang="ts">
import type { MediaFile } from '@vibezz/types'

const { apiFetch } = useApi()
const config = useRuntimeConfig()

const { data: files, refresh } = await useAsyncData('media', () =>
  apiFetch<MediaFile[]>('/api/admin/media'),
)

const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

function triggerUpload() {
  fileInput.value?.click()
}

async function handleUpload(event: Event) {
  const input = event.target as HTMLInputElement
  if (!input.files?.length) return

  uploading.value = true
  try {
    for (const file of input.files) {
      const formData = new FormData()
      formData.append('file', file)
      await apiFetch('/api/admin/media', {
        method: 'POST',
        body: formData,
      })
    }
    refresh()
  } finally {
    uploading.value = false
    input.value = ''
  }
}

async function remove(id: string) {
  if (!confirm('Delete this file?')) return
  await apiFetch(`/api/admin/media/${id}`, { method: 'DELETE' })
  refresh()
}

function mediaUrl(id: string) {
  return `${config.public.apiBase}/api/public/media/${id}`
}

function isImage(mimeType: string) {
  return mimeType.startsWith('image/')
}

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100">Media Library</h2>
      <UButton @click="triggerUpload" :loading="uploading" icon="i-heroicons-arrow-up-tray">
        Upload
      </UButton>
      <input
        ref="fileInput"
        type="file"
        multiple
        class="hidden"
        @change="handleUpload"
      />
    </div>

    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
      <UCard
        v-for="file in files"
        :key="file.id"
        class="group relative"
      >
        <div class="aspect-square flex items-center justify-center bg-gray-100 dark:bg-gray-800 rounded overflow-hidden">
          <img
            v-if="isImage(file.mime_type)"
            :src="mediaUrl(file.id)"
            :alt="file.filename"
            class="object-cover w-full h-full"
          />
          <span v-else class="i-heroicons-document text-4xl text-gray-400" />
        </div>
        <div class="mt-2">
          <p class="text-xs font-medium text-gray-900 dark:text-gray-100 truncate">
            {{ file.filename }}
          </p>
          <p class="text-xs text-gray-500">{{ formatSize(file.size) }}</p>
        </div>
        <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
          <UButton
            icon="i-heroicons-trash"
            variant="solid"
            color="red"
            size="xs"
            @click="remove(file.id)"
          />
        </div>
      </UCard>
    </div>

    <div v-if="!files?.length" class="text-center py-12 text-gray-500">
      No files uploaded yet.
    </div>
  </div>
</template>
