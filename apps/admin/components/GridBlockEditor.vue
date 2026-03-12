<script setup lang="ts">
import type { GridBlockData } from '@vibezz/types'

const props = defineProps<{
  data: Record<string, unknown>
}>()

const emit = defineEmits<{
  update: [data: Record<string, unknown>]
}>()

function update(key: keyof GridBlockData, value: unknown) {
  emit('update', { ...props.data, [key]: value })
}

function updateResponsive(key: 'mobileColumns' | 'tabletColumns' | 'desktopColumns', value: number) {
  emit('update', { ...props.data, [key]: Math.max(1, Math.min(12, value)) })
}
</script>

<template>
  <div class="space-y-4">
    <!-- Grid Configuration -->
    <div class="grid grid-cols-3 gap-3">
      <UFormGroup label="Columns">
        <UInput
          :model-value="(data.columns as number) || 3"
          type="number"
          :min="1"
          :max="12"
          @update:model-value="update('columns', Number($event))"
        />
      </UFormGroup>
      
      <UFormGroup label="Column Gap">
        <UInput
          :model-value="(data.columnGap as string) || '16px'"
          placeholder="16px"
          @update:model-value="update('columnGap', $event)"
        />
      </UFormGroup>
      
      <UFormGroup label="Row Gap">
        <UInput
          :model-value="(data.rowGap as string) || '16px'"
          placeholder="16px"
          @update:model-value="update('rowGap', $event)"
        />
      </UFormGroup>
    </div>

    <!-- Responsive Columns -->
    <div class="border border-gray-200 dark:border-gray-700 rounded-lg p-3 space-y-3">
      <div class="text-xs font-medium text-gray-500 uppercase">Responsive Columns</div>
      
      <div class="grid grid-cols-3 gap-3">
        <UFormGroup label="Mobile (&lt; 640px)">
          <UInput
            :model-value="(data.mobileColumns as number) || 1"
            type="number"
            :min="1"
            :max="12"
            @update:model-value="updateResponsive('mobileColumns', Number($event))"
          />
        </UFormGroup>
        
        <UFormGroup label="Tablet (640px - 1024px)">
          <UInput
            :model-value="(data.tabletColumns as number) || 2"
            type="number"
            :min="1"
            :max="12"
            @update:model-value="updateResponsive('tabletColumns', Number($event))"
          />
        </UFormGroup>
        
        <UFormGroup label="Desktop (> 1024px)">
          <UInput
            :model-value="(data.desktopColumns as number) || 3"
            type="number"
            :min="1"
            :max="12"
            @update:model-value="updateResponsive('desktopColumns', Number($event))"
          />
        </UFormGroup>
      </div>
    </div>

    <!-- Visual Preview -->
    <div class="mt-4 p-3 bg-gray-50 dark:bg-gray-800 rounded">
      <div class="text-xs font-medium text-gray-500 mb-2">Grid Preview</div>
      <div class="grid gap-1" :style="{ gridTemplateColumns: `repeat(${(data.columns as number) || 3}, 1fr)` }">
        <div 
          v-for="i in Math.min(6, (data.columns as number) || 3)" 
          :key="i"
          class="h-8 bg-primary-200 dark:bg-primary-800 rounded flex items-center justify-center text-xs"
        >
          {{ i }}
        </div>
      </div>
      <div class="mt-2 text-xs text-gray-400">
        {{ (data.columns as number) || 3 }} columns × {{ (data.columnGap as string) || '16px' }} gap
      </div>
    </div>
  </div>
</template>