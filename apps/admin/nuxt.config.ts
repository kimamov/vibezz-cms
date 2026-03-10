export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },
  modules: ['@nuxt/ui'],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
    },
  },

  app: {
    head: {
      title: 'Vibezz CMS',
      meta: [
        { name: 'description', content: 'Vibezz CMS Admin Panel' },
      ],
    },
  },

  typescript: {
    strict: true,
  },
})
