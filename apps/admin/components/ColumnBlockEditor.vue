<script setup lang="ts">
import type { ColumnBlockData } from '@vibezz/types'

const props = defineProps<{
  data: Record<string, unknown>
}>()

const emit = defineEmits<{
  update: [data: Record<string, unknown>]
}>()

function update(key: keyof ColumnBlockData, value: unknown) {
  emit('update', { ...props.data, [key]: value })
}
</script>

<template>
  <div class="space-y-3">
    <UFormGroup label="Width">
      <USelect
        :model-value="(data.width as string) || ''"
        :options="[
          { value: '', label: 'Auto' },
          { value: '1/2', label: '1/2 (50%)' },
          { value: '1/3', label: '1/3 (33%)' },
          { value: '2/3', label: '2/3 (66%)' },
          { value: '1/4', label: '1/4 (25%)' },
          { value: '3/4', label: '3/4 (75%)' },
          { value: '1/5', label: '1/5 (20%)' },
          { value: '100%', label: '100%' },
        ]"
        @update:model-value="update('width', $event)"
      />
    </UFormGroup>

    <UFormGroup label="Padding">
      <UInput
        :model-value="(data.padding as string) || '0px'"
        placeholder="0px"
        @update:model-value="update('padding', $event)"
      />
    </UFormGroup>

    <div class="text-xs text-gray-500 mt-2">
      Columns are typically used inside Grid or Container blocks to control layout.
    </div>
  </div>
</template>