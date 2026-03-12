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
const isOpen = ref(false)
const searchQuery = ref('')

const blocks = computed({
  get: () => props.modelValue || [],
  set: (val) => emit('update:modelValue', val),
})

function generateId() {
  return crypto.randomUUID()
}

function addBlock(slug: string) {
  const def = blockTypes.value.find((b) => b.slug === slug)
  const defaultData = def ? { ...def.default_data } : {}
  const newBlock: ContentBlock = {
    id: generateId(),
    type: slug,
    data: defaultData,
    children: def?.isContainer ? [] : undefined,
  }
  emit('update:modelValue', [...blocks.value, newBlock])
  isOpen.value = false
  searchQuery.value = ''
}

function updateBlock(index: number, block: ContentBlock) {
  const updated = [...blocks.value]
  updated[index] = block
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

// Group block types by category
const blockTypesByCategory = computed(() => {
  const grouped: Record<string, BlockTypeDefinition[]> = {}
  
  for (const bt of blockTypes.value) {
    const category = bt.category || 'other'
    if (!grouped[category]) {
      grouped[category] = []
    }
    grouped[category].push(bt)
  }
  
  // Sort categories
  const categoryOrder = ['basic', 'media', 'layout', 'content', 'other']
  const sorted: Record<string, BlockTypeDefinition[]> = {}
  for (const cat of categoryOrder) {
    if (grouped[cat]) {
      sorted[cat] = grouped[cat]
    }
  }
  
  return sorted
})

const categoryLabels: Record<string, string> = {
  basic: 'Basic Blocks',
  media: 'Media',
  layout: 'Layout',
  content: 'Content',
  other: 'Other',
}

const categoryIcons: Record<string, string> = {
  basic: 'i-heroicons-document-text',
  media: 'i-heroicons-photo',
  layout: 'i-heroicons-view-columns',
  content: 'i-heroicons-newspaper',
  other: 'i-heroicons-square-3-stack-3d',
}

// Filter block types based on search
const filteredBlockTypes = computed(() => {
  if (!searchQuery.value) return blockTypesByCategory.value
  
  const query = searchQuery.value.toLowerCase()
  const filtered: Record<string, BlockTypeDefinition[]> = {}
  
  for (const [category, types] of Object.entries(blockTypesByCategory.value)) {
    const matchingTypes = types.filter(bt => 
      bt.label.toLowerCase().includes(query) ||
      bt.description?.toLowerCase().includes(query) ||
      bt.slug.toLowerCase().includes(query)
    )
    if (matchingTypes.length > 0) {
      filtered[category] = matchingTypes
    }
  }
  
  return filtered
})

function getBlockIcon(slug: string): string {
  const iconMap: Record<string, string> = {
    heading: 'i-heroicons-h1',
    text: 'i-heroicons-document-text',
    image: 'i-heroicons-photo',
    quote: 'i-heroicons-chat-bubble-bottom-center-text',
    code: 'i-heroicons-code-bracket',
    divider: 'i-heroicons-minus',
    container: 'i-heroicons-square-3-stack-3d',
    grid: 'i-heroicons-view-columns',
    section: 'i-heroicons-rectangle-stack',
    column: 'i-heroicons-table-cells',
  }
  return iconMap[slug] || 'i-heroicons-cube'
}
</script>

<template>
  <div class="space-y-3">
    <!-- Blocks List -->
    <div v-if="blocks.length === 0" class="text-center py-8 text-gray-500">
      <UIcon name="i-heroicons-squares-plus" class="w-12 h-12 mx-auto mb-2 text-gray-300" />
      <p>No blocks yet. Add your first block below.</p>
    </div>

    <div v-else class="space-y-3">
      <BlockItem
        v-for="(block, index) in blocks"
        :key="block.id"
        :block="block"
        :block-types="blockTypes"
        :level="0"
        @update="updateBlock(index, $event)"
        @remove="removeBlock(index)"
        @move="moveBlock(index, $event)"
      />
    </div>

    <!-- Add Block Button -->
    <div class="pt-4 border-t border-gray-200 dark:border-gray-700">
      <UButton
        color="primary"
        variant="soft"
        icon="i-heroicons-plus"
        size="sm"
        @click="isOpen = true"
      >
        Add Block
      </UButton>
    </div>

    <!-- Block Selector Dialog -->
    <UModal v-model="isOpen" :ui="{ width: 'md:max-w-3xl' }">
      <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
              Add Block
            </h3>
            <UButton
              color="gray"
              variant="ghost"
              icon="i-heroicons-x-mark"
              @click="isOpen = false"
            />
          </div>
        </template>

        <!-- Search -->
        <div class="p-4 pb-2">
          <UInput
            v-model="searchQuery"
            placeholder="Search blocks..."
            icon="i-heroicons-magnifying-glass"
            size="lg"
          />
        </div>

        <!-- Block Categories -->
        <div class="p-4 max-h-[60vh] overflow-y-auto">
          <div v-if="Object.keys(filteredBlockTypes).length === 0" class="text-center py-8 text-gray-500">
            <UIcon name="i-heroicons-magnifying-glass" class="w-12 h-12 mx-auto mb-2 text-gray-300" />
            <p>No blocks found matching "{{ searchQuery }}"</p>
          </div>

          <div
            v-for="(types, category) in filteredBlockTypes"
            :key="category"
            class="mb-6 last:mb-0"
          >
            <div class="flex items-center gap-2 mb-3">
              <UIcon :name="categoryIcons[category] || 'i-heroicons-folder'" class="text-gray-400" />
              <h4 class="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
                {{ categoryLabels[category] || category }}
              </h4>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <button
                v-for="blockType in types"
                :key="blockType.slug"
                class="flex items-start gap-3 p-3 text-left rounded-lg border border-gray-200 dark:border-gray-700 hover:border-primary-500 dark:hover:border-primary-500 hover:bg-primary-50 dark:hover:bg-primary-900/20 transition-all group"
                @click="addBlock(blockType.slug)"
              >
                <div class="flex-shrink-0 w-10 h-10 rounded-lg bg-gray-100 dark:bg-gray-800 flex items-center justify-center group-hover:bg-primary-100 dark:group-hover:bg-primary-900/30">
                  <UIcon
                    :name="getBlockIcon(blockType.slug)"
                    class="w-5 h-5 text-gray-500 group-hover:text-primary-600 dark:group-hover:text-primary-400"
                  />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="font-medium text-gray-900 dark:text-white">
                      {{ blockType.label }}
                    </span>
                    <UBadge
                      v-if="blockType.isContainer"
                      size="xs"
                      color="primary"
                      variant="soft"
                    >
                      Container
                    </UBadge>
                  </div>
                  <p v-if="blockType.description" class="text-sm text-gray-500 dark:text-gray-400 mt-0.5 line-clamp-2">
                    {{ blockType.description }}
                  </p>
                  <p v-else class="text-sm text-gray-400 dark:text-gray-500 mt-0.5">
                    {{ blockType.slug }}
                  </p>
                </div>
              </button>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex justify-end">
            <UButton
              color="gray"
              variant="ghost"
              @click="isOpen = false"
            >
              Cancel
            </UButton>
          </div>
        </template>
      </UCard>
    </UModal>
  </div>
</template>