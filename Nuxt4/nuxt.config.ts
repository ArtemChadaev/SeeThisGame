// https://nuxt.com/docs/api/configuration/nuxt-config
const apiBaseUrl = process.env.API_BASE_URL || 'http://localhost:8080'
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    '@pinia/nuxt',
    'pinia-plugin-persistedstate/nuxt',
  ],
  css: ['~/assets/css/main.css'],
  ui: {
    colorMode: false, //на время потом темную тему буду делать
  },
  pinia: {
    storesDirs: ['./app/stores/**'],
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8000/',
    },
  },
  routeRules: {
    '/api/_nuxt_icon/**': {}, // Для работы иконок
    '/api/settings/**': { proxy: `${apiBaseUrl}/api/settings/**` },
    '/auth/**': { proxy: `${apiBaseUrl}/auth/**` },
  },
})
