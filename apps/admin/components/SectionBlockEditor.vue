<script setup lang="ts">
import type { SectionBlockData } from '@vibezz/types'

const props = defineProps<{
  data: Record<string, unknown>
}>()

const emit = defineEmits<{
  update: [data: Record<string, unknown>]
}>()

const backgroundSizeOptions = [
  { value: 'cover', label: 'Cover' },
  { value: 'contain', label: 'Contain' },
  { value: 'auto', label: 'Auto' },
]

function update(key: keyof SectionBlockData, value: unknown) {
  emit('update', { ...props.data, [key]: value })
}
</script>

<template>
  <div class="space-y-3">
    <!-- Background -->
    <UFormGroup label="Background Color">
      <div class="flex gap-2">
        <input
          type="color"
          :value="(data.background as string)?.startsWith('#') ? data.background : '#ffffff'"
          class="w-10 h-10 rounded border border-gray-300 cursor-pointer"
          @input="update('background', ($event.target as HTMLInputElement).value)"
        />
        <UInput
          :model-value="(data.background as string)"
          placeholder="Color or transparent"
          class="flex-1"
          @update:model-value="update('background', $event)"
        />
      </div>
    </UFormGroup>

    <!-- Background Image -->
    <UFormGroup label="Background Image">
      <div class="space-y-2">
        <MediaPicker
          :model-value="(data.backgroundImage as string)"
          @update:model-value="update('backgroundImage', $event)"
        />
        <USelect
          v-if="data.backgroundImage"
          :model-value="(data.backgroundSize as string) || 'cover'"
          :options="backgroundSizeOptions"
          placeholder="Background size"
          @update:model-value="update('backgroundSize', $event)"
        />
      </div>
    </UFormGroup>

    <!-- Padding -->
    <div class="grid grid-cols-2 gap-3">
      <UFormGroup label="Vertical Padding">
        <UInput
          :model-value="(data.paddingY as string) || '64px'"
          placeholder="64px"
          @update:model-value="update('paddingY', $event)"
        />
      </UFormGroup>
      
      <UFormGroup label="Horizontal Padding">
        <UInput
          :model-value="(data.paddingX as string) || '0px'"
          placeholder="0px"
          @update:model-value="update('paddingX', $event)"
        />
      </UFormGroup>
    </div>

    <!-- Min Height -->
    <UFormGroup label="Min Height">
      <UInput
        :model-value="(data.minHeight as string) || ''"
        placeholder="e.g., 400px or 50vh"
        @update:model-value="update('minHeight', $event)"
      />
    </UFormGroup>

    <!-- Preview -->
    <div 
      class="mt-4 p-4 rounded border-2 border-dashed border-gray-300 dark:border-gray-600 text-center"
      :style="{
        backgroundColor: (data.background as string) || 'transparent',
        padding: `${(data.paddingY as string) || '20px'} ${(data.paddingX as string) || '20px'}`,
        minHeight: (data.minHeight as string) || '100px'
      }"
    >
      <span class="text-xs text-gray-400">Section Preview</span>
    </div>
  </div>
</template>