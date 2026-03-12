<script setup lang="ts">
import type { ContentBlock, BlockTypeDefinition, BlockFieldDefinition } from '@vibezz/types'

const props = defineProps<{
  block: ContentBlock
  blockTypes: BlockTypeDefinition[]
  level?: number
}>()

const emit = defineEmits<{
  update: [block: ContentBlock]
  remove: []
  move: [direction: -1 | 1]
}>()

const isOpen = ref(true)
const isAddChildOpen = ref(false)
const childSearchQuery = ref('')
const level = computed(() => props.level ?? 0)
const indentClass = computed(() => `ml-${Math.min(level.value * 4, 12)}`)

const blockType = computed(() => 
  props.blockTypes.find((b) => b.slug === props.block.type)
)

const isContainer = computed(() => blockType.value?.isContainer ?? false)

// Group block types by category for child selector
const childBlockTypesByCategory = computed(() => {
  const grouped: Record<string, BlockTypeDefinition[]> = {}
  const availableTypes = props.blockTypes.filter((bt) => !bt.isContainer || level.value < 3)
  
  for (const bt of availableTypes) {
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

// Filter child block types based on search
const filteredChildBlockTypes = computed(() => {
  if (!childSearchQuery.value) return childBlockTypesByCategory.value
  
  const query = childSearchQuery.value.toLowerCase()
  const filtered: Record<string, BlockTypeDefinition[]> = {}
  
  for (const [category, types] of Object.entries(childBlockTypesByCategory.value)) {
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

function addChildBlockAndClose(slug: string) {
  addChildBlock(slug)
  isAddChildOpen.value = false
  childSearchQuery.value = ''
}

function getFieldType(field: BlockFieldDefinition): string {
  const typeMap: Record<string, string> = {
    text: 'text',
    number: 'number',
    boolean: 'checkbox',
    textarea: 'textarea',
    richtext: 'textarea',
    select: 'select',
    media: 'media',
    color: 'text',
    code: 'textarea',
  }
  return typeMap[field.type] || 'text'
}

function updateData(key: string, value: unknown) {
  emit('update', {
    ...props.block,
    data: { ...props.block.data, [key]: value },
  })
}

function updateChildren(children: ContentBlock[]) {
  emit('update', {
    ...props.block,
    children,
  })
}

function addChildBlock(slug: string) {
  const def = props.blockTypes.find((b) => b.slug === slug)
  const defaultData = def ? { ...def.default_data } : {}
  const newChild: ContentBlock = {
    id: crypto.randomUUID(),
    type: slug,
    data: defaultData,
    children: def?.isContainer ? [] : undefined,
  }
  const currentChildren = props.block.children || []
  updateChildren([...currentChildren, newChild])
}

function updateChildBlock(index: number, child: ContentBlock) {
  const children = [...(props.block.children || [])]
  children[index] = child
  updateChildren(children)
}

function removeChildBlock(index: number) {
  const children = (props.block.children || []).filter((_, i) => i !== index)
  updateChildren(children)
}

function moveChildBlock(index: number, direction: -1 | 1) {
  const target = index + direction
  const children = props.block.children || []
  if (target < 0 || target >= children.length) return
  const updated = [...children]
  const temp = updated[index]
  updated[index] = updated[target]
  updated[target] = temp
  updateChildren(updated)
}

function generateId() {
  return crypto.randomUUID()
}
</script>

<template>
  <div 
    class="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden"
    :class="{ 'shadow-sm': level === 0 }"
  >
    <!-- Header -->
    <div 
      class="flex items-center justify-between px-3 py-2 bg-gray-50 dark:bg-gray-800/50 border-b border-gray-200 dark:border-gray-700"
    >
      <div class="flex items-center gap-2">
        <UButton
          v-if="isContainer"
          :icon="isOpen ? 'i-heroicons-chevron-down' : 'i-heroicons-chevron-right'"
          variant="ghost"
          size="2xs"
          color="gray"
          @click="isOpen = !isOpen"
        />
        <UIcon v-if="blockType?.icon" :name="blockType.icon" class="text-gray-400" />
        <span class="text-xs font-medium text-gray-600 dark:text-gray-400 uppercase tracking-wide">
          {{ blockType?.label ?? block.type }}
        </span>
        <UBadge v-if="isContainer" size="xs" color="primary" variant="soft">
          {{ block.children?.length ?? 0 }} items
        </UBadge>
      </div>
      <div class="flex gap-1">
        <UButton
          icon="i-heroicons-chevron-up"
          variant="ghost"
          size="2xs"
          color="gray"
          @click="emit('move', -1)"
        />
        <UButton
          icon="i-heroicons-chevron-down"
          variant="ghost"
          size="2xs"
          color="gray"
          @click="emit('move', 1)"
        />
        <UButton
          icon="i-heroicons-trash"
          variant="ghost"
          color="red"
          size="2xs"
          @click="emit('remove')"
        />
      </div>
    </div>

    <!-- Content -->
    <div v-show="!isContainer || isOpen" class="p-3 space-y-3">
      <!-- Render block-specific editor based on type -->
      
      <!-- Heading -->
      <template v-if="block.type === 'heading'">
        <div class="space-y-2">
          <UInput
            :model-value="(block.data.text as string)"
            placeholder="Heading text"
            @update:model-value="updateData('text', $event)"
          />
          <USelect
            :model-value="String(block.data.level ?? 2)"
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
            @update:model-value="updateData('level', Number($event))"
          />
        </div>
      </template>

      <!-- Text -->
      <template v-else-if="block.type === 'text'">
        <UTextarea
          :model-value="(block.data.content as string)"
          placeholder="Write your content..."
          :rows="4"
          @update:model-value="updateData('content', $event)"
        />
      </template>

      <!-- Image -->
      <template v-else-if="block.type === 'image'">
        <div class="space-y-2">
          <MediaPicker
            :model-value="(block.data.media_id as string)"
            @update:model-value="updateData('media_id', $event)"
          />
          <UInput
            :model-value="(block.data.caption as string)"
            placeholder="Caption (optional)"
            @update:model-value="updateData('caption', $event)"
          />
          <UInput
            :model-value="(block.data.alt as string)"
            placeholder="Alt text (optional)"
            @update:model-value="updateData('alt', $event)"
          />
        </div>
      </template>

      <!-- Quote -->
      <template v-else-if="block.type === 'quote'">
        <div class="space-y-2">
          <UTextarea
            :model-value="(block.data.text as string)"
            placeholder="Quote text"
            :rows="3"
            @update:model-value="updateData('text', $event)"
          />
          <UInput
            :model-value="(block.data.attribution as string)"
            placeholder="Attribution (optional)"
            @update:model-value="updateData('attribution', $event)"
          />
        </div>
      </template>

      <!-- Code -->
      <template v-else-if="block.type === 'code'">
        <div class="space-y-2">
          <USelect
            :model-value="(block.data.language as string) || 'plaintext'"
            :options="[
              { value: 'plaintext', label: 'Plain Text' },
              { value: 'javascript', label: 'JavaScript' },
              { value: 'typescript', label: 'TypeScript' },
              { value: 'go', label: 'Go' },
              { value: 'python', label: 'Python' },
              { value: 'html', label: 'HTML' },
              { value: 'css', label: 'CSS' },
              { value: 'sql', label: 'SQL' },
              { value: 'json', label: 'JSON' },
            ]"
            @update:model-value="updateData('language', $event)"
          />
          <UTextarea
            :model-value="(block.data.code as string)"
            placeholder="Code..."
            :rows="6"
            class="font-mono text-sm"
            @update:model-value="updateData('code', $event)"
          />
        </div>
      </template>

      <!-- Divider -->
      <template v-else-if="block.type === 'divider'">
        <div class="py-2">
          <hr class="border-gray-300 dark:border-gray-600" />
        </div>
      </template>

      <!-- Container Block -->
      <template v-else-if="block.type === 'container'">
        <ContainerBlockEditor
          :data="block.data"
          @update="emit('update', { ...block, data: $event })"
        />
      </template>

      <!-- Grid Block -->
      <template v-else-if="block.type === 'grid'">
        <GridBlockEditor
          :data="block.data"
          @update="emit('update', { ...block, data: $event })"
        />
      </template>

      <!-- Section Block -->
      <template v-else-if="block.type === 'section'">
        <SectionBlockEditor
          :data="block.data"
          @update="emit('update', { ...block, data: $event })"
        />
      </template>

      <!-- Column Block -->
      <template v-else-if="block.type === 'column'">
        <ColumnBlockEditor
          :data="block.data"
          @update="emit('update', { ...block, data: $event })"
        />
      </template>

      <!-- Plugin blocks with fields -->
      <template v-else-if="blockType?.fields?.length">
        <div class="space-y-3">
          <UFormGroup
            v-for="field in blockType.fields"
            :key="field.name"
            :label="field.label"
            :required="field.required"
          >
            <template v-if="field.type === 'text'">
              <UInput
                :model-value="(block.data[field.name] as string)"
                @update:model-value="updateData(field.name, $event)"
              />
            </template>
            <template v-else-if="field.type === 'number'">
              <UInput
                :model-value="(block.data[field.name] as number)"
                type="number"
                @update:model-value="updateData(field.name, Number($event))"
              />
            </template>
            <template v-else-if="field.type === 'boolean'">
              <UCheckbox
                :model-value="(block.data[field.name] as boolean)"
                @update:model-value="updateData(field.name, $event)"
              />
            </template>
            <template v-else-if="field.type === 'textarea' || field.type === 'richtext'">
              <UTextarea
                :model-value="(block.data[field.name] as string)"
                :rows="4"
                @update:model-value="updateData(field.name, $event)"
              />
            </template>
            <template v-else-if="field.type === 'select' && field.options">
              <USelect
                :model-value="(block.data[field.name] as string)"
                :options="(field.options as string[]).map(o => ({ value: o, label: o }))"
                @update:model-value="updateData(field.name, $event)"
              />
            </template>
            <template v-else-if="field.type === 'media'">
              <MediaPicker
                :model-value="(block.data[field.name] as string)"
                @update:model-value="updateData(field.name, $event)"
              />
            </template>
            <template v-else-if="field.type === 'color'">
              <div class="flex gap-2">
                <input
                  type="color"
                  :value="(block.data[field.name] as string) || '#000000'"
                  class="w-10 h-10 rounded border border-gray-300 cursor-pointer"
                  @input="updateData(field.name, ($event.target as HTMLInputElement).value)"
                />
                <UInput
                  :model-value="(block.data[field.name] as string)"
                  placeholder="#000000"
                  @update:model-value="updateData(field.name, $event)"
                />
              </div>
            </template>
          </UFormGroup>
        </div>
      </template>

      <!-- Unknown block type -->
      <template v-else>
        <div class="py-2 text-sm text-gray-500">
          Block type "{{ block.type }}"
        </div>
      </template>

      <!-- Nested Children for Container Blocks -->
      <div v-if="isContainer && block.children" class="mt-4 space-y-2">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-gray-500 uppercase">Children</span>
        </div>
        
        <div class="space-y-2">
          <BlockItem
            v-for="(child, index) in block.children"
            :key="child.id"
            :block="child"
            :block-types="blockTypes"
            :level="level + 1"
            @update="updateChildBlock(index, $event)"
            @remove="removeChildBlock(index)"
            @move="moveChildBlock(index, $event)"
          />
        </div>

        <!-- Add child block -->
        <div class="pt-2">
          <UButton
            size="xs"
            variant="soft"
            color="primary"
            icon="i-heroicons-plus"
            @click="isAddChildOpen = true"
          >
            Add Block
          </UButton>
        </div>

        <!-- Add Child Block Dialog -->
        <UModal v-model="isAddChildOpen" :ui="{ width: 'md:max-w-3xl' }">
          <UCard :ui="{ ring: '', divide: 'divide-y divide-gray-100 dark:divide-gray-800' }">
            <template #header>
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white">
                  Add Block to {{ blockType?.label || 'Container' }}
                </h3>
                <UButton
                  color="gray"
                  variant="ghost"
                  icon="i-heroicons-x-mark"
                  @click="isAddChildOpen = false"
                />
              </div>
            </template>

            <!-- Search -->
            <div class="p-4 pb-2">
              <UInput
                v-model="childSearchQuery"
                placeholder="Search blocks..."
                icon="i-heroicons-magnifying-glass"
                size="lg"
              />
            </div>

            <!-- Block Categories -->
            <div class="p-4 max-h-[50vh] overflow-y-auto">
              <div v-if="Object.keys(filteredChildBlockTypes).length === 0" class="text-center py-8 text-gray-500">
                <UIcon name="i-heroicons-magnifying-glass" class="w-12 h-12 mx-auto mb-2 text-gray-300" />
                <p>No blocks found matching "{{ childSearchQuery }}"</p>
              </div>

              <div
                v-for="(types, category) in filteredChildBlockTypes"
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
                    @click="addChildBlockAndClose(blockType.slug)"
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
                  @click="isAddChildOpen = false"
                >
                  Cancel
                </UButton>
              </div>
            </template>
          </UCard>
        </UModal>
      </div>
    </div>
  </div>
</template>