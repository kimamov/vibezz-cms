export default defineNuxtConfig({
  compatibilityDate: '2025-01-01',
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8080',
    },
  },

  app: {
    head: {
      title: 'Vibezz',
      meta: [
        { name: 'description', content: 'Powered by Vibezz CMS' },
      ],
    },
  },

  typescript: {
    strict: true,
  },
})
