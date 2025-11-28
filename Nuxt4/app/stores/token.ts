import { FetchError } from 'ofetch'
export const useTokenStore = defineStore(
  'token',
  () => {
    const refreshToken = ref<string | null>(null)
    const accessToken = ref<string | null>(null)
    const expiresAt = ref<number | null>(null)

    const setTokens = (refresh: string, access: string) => {
      refreshToken.value = refresh
      accessToken.value = access
      expiresAt.value = Date.now() + 10 * 60 * 1000
    }

    const logout = () => {
      refreshToken.value = null
      accessToken.value = null
      expiresAt.value = null
    }

    const isTokenExpired = computed(() => {
      if (!expiresAt.value) return true
      return Date.now() > expiresAt.value
    })

    const updateToken = async () => {
      const route = useRoute()
      try {
        const response = await $fetch<{
          refreshToken: string
          accessToken: string
        }>('/auth/refresh', {
          method: 'POST',
          body: {
            refreshToken: refreshToken.value,
          },
        })

        setTokens(response.refreshToken, response.accessToken)
      } catch (error: unknown) {
        if (error instanceof FetchError) {
          if (
            error.response &&
            error.response.status === 400 &&
            error.data.error === 'invalid_token'
          ) {
            logout()
            await navigateTo({
              path: '/login',
              query: {
                previous: route.fullPath,
                error: 'invalid_token',
                code: '400',
              },
            })
          }
        }
      }
    }

    return {
      refreshToken,
      accessToken,
      setTokens,
      logout,
      isTokenExpired,
      updateToken
    }
  },

  {
    persist: true,
  },
)
