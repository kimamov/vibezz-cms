<script setup lang="ts">
import type { ContentBlock, BlockTypeDefinition } from '@vibezz/types'

const props = defineProps<{
  modelValue?: ContentBlock[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: ContentBlock[]]
}>()

const { apiFetch } = useApi()

const { data: blockTypesFromApi } = await useAsyncData('block-types', () =>
  apiFetch<BlockTypeDefinition[]>('/api/admin/block-types'),
)

const blockTypes = computed(() => blockTypesFromApi.value ?? [])

const blocks = computed({
  get: () => props.modelValue || [],
  set: (val) => emit('update:modelValue', val),
})

function getBlockType(slug: string): BlockTypeDefinition | undefined {
  return blockTypes.value.find((b) => b.slug === slug)
}

function generateId() {
  return crypto.randomUUID()
}

function addBlock(slug: string) {
  const def = getBlockType(slug)
  const defaultData = def ? { ...def.default_data } : {}
  const updated = [
    ...blocks.value,
    {
      id: generateId(),
      type: slug,
      data: defaultData,
    },
  ]
  emit('update:modelValue', updated)
}

function updateBlock(index: number, data: Record<string, unknown>) {
  const updated = [...blocks.value]
  updated[index] = { ...updated[index], data }
  emit('update:modelValue', updated)
}

function removeBlock(index: number) {
  const updated = blocks.value.filter((_, i) => i !== index)
  emit('update:modelValue', updated)
}

function moveBlock(index: number, direction: -1 | 1) {
  const target = index + direction
  if (target < 0 || target >= blocks.value.length) return
  const updated = [...blocks.value]
  const temp = updated[index]
  updated[index] = updated[target]
  updated[target] = temp
  emit('update:modelValue', updated)
}

function blockLabel(type: string) {
  return getBlockType(type)?.label ?? type
}
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="(block, index) in blocks"
      :key="block.id"
      class="border border-gray-200 dark:border-gray-700 rounded-lg"
    >
      <div class="flex items-center justify-between px-3 py-2 bg-gray-50 dark:bg-gray-800/50 rounded-t-lg border-b border-gray-200 dark:border-gray-700">
        <span class="text-xs font-medium text-gray-500 uppercase tracking-wide">
          {{ blockLabel(block.type) }}
        </span>
        <div class="flex gap-1">
          <UButton
            icon="i-heroicons-chevron-up"
            variant="ghost"
            size="2xs"
            :disabled="index === 0"
            @click="moveBlock(index, -1)"
          />
          <UButton
            icon="i-heroicons-chevron-down"
            variant="ghost"
            size="2xs"
            :disabled="index === blocks.length - 1"
            @click="moveBlock(index, 1)"
          />
          <UButton
            icon="i-heroicons-trash"
            variant="ghost"
            color="red"
            size="2xs"
            @click="removeBlock(index)"
          />
        </div>
      </div>

      <div class="p-3">
        <!-- Heading -->
        <div v-if="block.type === 'heading'" class="space-y-2">
          <UInput
            :model-value="(block.data.text as string)"
            placeholder="Heading text"
            @update:model-value="updateBlock(index, { ...block.data, text: $event })"
          />
          <USelect
            :model-value="String(block.data.level)"
            :options="[
              { value: '1', label: 'H1' },
              { value: '2', label: 'H2' },
              { value: '3', label: 'H3' },
              { value: '4', label: 'H4' },
              { value: '5', label: 'H5' },
              { value: '6', label: 'H6' },
            ]"
            option-attribute="label"
            value-attribute="value"
            @update:model-value="updateBlock(index, { ...block.data, level: Number($event) })"
          />
        </div>

        <!-- Text -->
        <div v-else-if="block.type === 'text'">
          <UTextarea
            :model-value="(block.data.content as string)"
            placeholder="Write your content..."
            :rows="4"
            @update:model-value="updateBlock(index, { ...block.data, content: $event })"
          />
        </div>

        <!-- Image -->
        <div v-else-if="block.type === 'image'" class="space-y-2">
          <MediaPicker
            :model-value="(block.data.media_id as string)"
            @update:model-value="updateBlock(index, { ...block.data, media_id: $event })"
          />
          <UInput
            :model-value="(block.data.caption as string)"
            placeholder="Caption (optional)"
            @update:model-value="updateBlock(index, { ...block.data, caption: $event })"
          />
          <UInput
            :model-value="(block.data.alt as string)"
            placeholder="Alt text (optional)"
            @update:model-value="updateBlock(index, { ...block.data, alt: $event })"
          />
        </div>

        <!-- Quote -->
        <div v-else-if="block.type === 'quote'" class="space-y-2">
          <UTextarea
            :model-value="(block.data.text as string)"
            placeholder="Quote text"
            :rows="3"
            @update:model-value="updateBlock(index, { ...block.data, text: $event })"
          />
          <UInput
            :model-value="(block.data.attribution as string)"
            placeholder="Attribution (optional)"
            @update:model-value="updateBlock(index, { ...block.data, attribution: $event })"
          />
        </div>

        <!-- Code -->
        <div v-else-if="block.type === 'code'" class="space-y-2">
          <UInput
            :model-value="(block.data.language as string)"
            placeholder="Language (e.g. javascript)"
            @update:model-value="updateBlock(index, { ...block.data, language: $event })"
          />
          <UTextarea
            :model-value="(block.data.code as string)"
            placeholder="Code..."
            :rows="6"
            class="font-mono"
            @update:model-value="updateBlock(index, { ...block.data, code: $event })"
          />
        </div>

        <!-- Divider -->
        <div v-else-if="block.type === 'divider'" class="py-2">
          <hr class="border-gray-300 dark:border-gray-600" />
        </div>

        <!-- Plugin blocks: config form from config_fields -->
        <div v-else-if="getBlockType(block.type)?.config_fields?.length" class="space-y-2">
          <UFormGroup
            v-for="field in getBlockType(block.type)!.config_fields"
            :key="field.slug"
            :label="field.name"
            :required="field.required"
          >
            <UInput
              v-if="field.type === 'text' || field.type === 'number'"
              :model-value="(block.data[field.slug] as string | number)"
              :type="field.type"
              @update:model-value="updateBlock(index, { ...block.data, [field.slug]: field.type === 'number' ? Number($event) : $event })"
            />
            <UCheckbox
              v-else-if="field.type === 'boolean'"
              :model-value="(block.data[field.slug] as boolean)"
              @update:model-value="updateBlock(index, { ...block.data, [field.slug]: $event })"
            />
          </UFormGroup>
        </div>

        <!-- Unknown block type -->
        <div v-else class="py-2 text-sm text-gray-500">
          Block type "{{ block.type }}"
        </div>
      </div>
    </div>

    <!-- Add block menu -->
    <div class="flex flex-wrap gap-2 pt-2">
      <UButton
        v-for="bt in blockTypes"
        :key="bt.slug"
        size="sm"
        variant="outline"
        :icon="bt.icon"
        @click="addBlock(bt.slug)"
      >
        {{ bt.label }}
      </UButton>
    </div>
  </div>
</template>
