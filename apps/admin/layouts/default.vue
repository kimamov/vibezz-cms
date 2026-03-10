<script setup lang="ts">
const auth = useAuth()
const route = useRoute()

const navigation = [
  { label: 'Dashboard', to: '/', icon: 'i-heroicons-home' },
  { label: 'Content Types', to: '/content-types', icon: 'i-heroicons-rectangle-stack' },
  { label: 'Entries', to: '/entries', icon: 'i-heroicons-document-text' },
  { label: 'Media', to: '/media', icon: 'i-heroicons-photo' },
]

const isActive = (path: string) => {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <div class="flex h-screen bg-gray-50 dark:bg-gray-950">
    <aside class="w-64 bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-800 flex flex-col">
      <div class="p-4 border-b border-gray-200 dark:border-gray-800">
        <h1 class="text-xl font-bold text-primary">Vibezz CMS</h1>
      </div>

      <nav class="flex-1 p-3 space-y-1">
        <NuxtLink
          v-for="item in navigation"
          :key="item.to"
          :to="item.to"
          class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors"
          :class="isActive(item.to)
            ? 'bg-primary-50 dark:bg-primary-950 text-primary font-medium'
            : 'text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800'"
        >
          <span :class="item.icon" class="w-5 h-5" />
          {{ item.label }}
        </NuxtLink>
      </nav>

      <div class="p-3 border-t border-gray-200 dark:border-gray-800">
        <div class="flex items-center justify-between px-3 py-2">
          <div class="text-sm">
            <p class="font-medium text-gray-900 dark:text-gray-100">{{ auth.user.value?.name }}</p>
            <p class="text-gray-500 text-xs">{{ auth.user.value?.role }}</p>
          </div>
          <UButton
            icon="i-heroicons-arrow-right-on-rectangle"
            variant="ghost"
            size="sm"
            @click="auth.logout()"
          />
        </div>
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto">
      <div class="p-8">
        <slot />
      </div>
    </main>
  </div>
</template>
