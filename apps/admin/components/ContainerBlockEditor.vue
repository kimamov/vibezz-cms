<script setup lang="ts">
import type { ContainerBlockData } from '@vibezz/types'

const props = defineProps<{
  data: Record<string, unknown>
}>()

const emit = defineEmits<{
  update: [data: Record<string, unknown>]
}>()

const layoutOptions = [
  { value: 'flex', label: 'Flex' },
  { value: 'grid', label: 'Grid' },
  { value: 'row', label: 'Row' },
  { value: 'column', label: 'Column' },
  { value: 'stack', label: 'Stack' },
]

function update(key: keyof ContainerBlockData, value: unknown) {
  emit('update', { ...props.data, [key]: value })
}
</script>

<template>
  <div class="space-y-3">
    <div class="grid grid-cols-2 gap-3">
      <UFormGroup label="Layout">
        <USelect
          :model-value="(data.layout as string) || 'flex'"
          :options="layoutOptions"
          @update:model-value="update('layout', $event)"
        />
      </UFormGroup>
      
      <UFormGroup label="Gap">
        <UInput
          :model-value="(data.gap as string) || '16px'"
          placeholder="16px"
          @update:model-value="update('gap', $event)"
        />
      </UFormGroup>
    </div>

    <div class="grid grid-cols-2 gap-3">
      <UFormGroup label="Padding">
        <UInput
          :model-value="(data.padding as string) || '0px'"
          placeholder="0px"
          @update:model-value="update('padding', $event)"
        />
      </UFormGroup>
      
      <UFormGroup label="Max Width">
        <UInput
          :model-value="(data.maxWidth as string) || ''"
          placeholder="1200px or 100%"
          @update:model-value="update('maxWidth', $event)"
        />
      </UFormGroup>
    </div>

    <UFormGroup label="Background">
      <div class="flex gap-2">
        <input
          type="color"
          :value="(data.background as string)?.startsWith('#') ? data.background : '#ffffff'"
          class="w-10 h-10 rounded border border-gray-300 cursor-pointer"
          @input="update('background', ($event.target as HTMLInputElement).value)"
        />
        <UInput
          :model-value="(data.background as string)"
          placeholder="Color, gradient, or transparent"
          class="flex-1"
          @update:model-value="update('background', $event)"
        />
      </div>
    </UFormGroup>

    <!-- Preview -->
    <div class="mt-4 p-3 bg-gray-50 dark:bg-gray-800 rounded text-xs text-gray-500">
      <div class="font-medium mb-1">Preview:</div>
      <div class="font-mono">
        Layout: {{ data.layout || 'flex' }} | Gap: {{ data.gap || '16px' }} | Padding: {{ data.padding || '0px' }}
      </div>
    </div>
  </div>
</template>