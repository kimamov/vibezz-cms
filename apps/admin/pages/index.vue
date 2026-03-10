<script setup lang="ts">
import type { Entry, ContentType } from '@vibezz/types'

const { apiFetch } = useApi()

const { data: contentTypes } = await useAsyncData('dashboard-types', () =>
  apiFetch<ContentType[]>('/api/admin/content-types'),
)

const { data: recentEntries } = await useAsyncData('dashboard-entries', () =>
  apiFetch<Entry[]>('/api/admin/entries'),
)
</script>

<template>
  <div>
    <h2 class="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">Dashboard</h2>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <UCard>
        <div class="text-center">
          <p class="text-3xl font-bold text-primary">{{ contentTypes?.length || 0 }}</p>
          <p class="text-sm text-gray-500 mt-1">Content Types</p>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <p class="text-3xl font-bold text-primary">{{ recentEntries?.length || 0 }}</p>
          <p class="text-sm text-gray-500 mt-1">Total Entries</p>
        </div>
      </UCard>
      <UCard>
        <div class="text-center">
          <p class="text-3xl font-bold text-primary">
            {{ recentEntries?.filter(e => e.status === 'published').length || 0 }}
          </p>
          <p class="text-sm text-gray-500 mt-1">Published</p>
        </div>
      </UCard>
    </div>

    <UCard>
      <template #header>
        <h3 class="text-lg font-semibold">Recent Entries</h3>
      </template>

      <UTable
        :rows="recentEntries?.slice(0, 10) || []"
        :columns="[
          { key: 'title', label: 'Title' },
          { key: 'status', label: 'Status' },
          { key: 'path', label: 'Path' },
          { key: 'updated_at', label: 'Updated' },
        ]"
      >
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
      </UTable>
    </UCard>
  </div>
</template>
