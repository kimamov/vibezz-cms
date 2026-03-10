import type { TokenPair, User } from '@vibezz/types'

const tokenState = ref<string | null>(null)
const userState = ref<User | null>(null)

export function useAuth() {
  const config = useRuntimeConfig()

  const token = computed(() => tokenState.value)
  const user = computed(() => userState.value)
  const isAuthenticated = computed(() => !!tokenState.value)

  function init() {
    if (import.meta.client) {
      tokenState.value = localStorage.getItem('vibezz_token')
      const refreshToken = localStorage.getItem('vibezz_refresh_token')
      if (tokenState.value) {
        fetchMe().catch(() => {
          if (refreshToken) {
            refresh(refreshToken).catch(() => logout())
          } else {
            logout()
          }
        })
      }
    }
  }

  async function login(email: string, password: string) {
    const response = await fetch(`${config.public.apiBase}/api/admin/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    })

    if (!response.ok) {
      const data = await response.json().catch(() => ({ error: 'Login failed' }))
      throw new Error(data.error)
    }

    const tokens: TokenPair = await response.json()
    tokenState.value = tokens.access_token
    if (import.meta.client) {
      localStorage.setItem('vibezz_token', tokens.access_token)
      localStorage.setItem('vibezz_refresh_token', tokens.refresh_token)
    }

    await fetchMe()
  }

  async function refresh(refreshToken: string) {
    const response = await fetch(`${config.public.apiBase}/api/admin/auth/refresh`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    })

    if (!response.ok) throw new Error('Refresh failed')

    const tokens: TokenPair = await response.json()
    tokenState.value = tokens.access_token
    if (import.meta.client) {
      localStorage.setItem('vibezz_token', tokens.access_token)
      localStorage.setItem('vibezz_refresh_token', tokens.refresh_token)
    }
  }

  async function fetchMe() {
    const response = await fetch(`${config.public.apiBase}/api/admin/me`, {
      headers: { Authorization: `Bearer ${tokenState.value}` },
    })
    if (!response.ok) throw new Error('Not authenticated')
    userState.value = await response.json()
  }

  function logout() {
    tokenState.value = null
    userState.value = null
    if (import.meta.client) {
      localStorage.removeItem('vibezz_token')
      localStorage.removeItem('vibezz_refresh_token')
    }
    navigateTo('/login')
  }

  return { token, user, isAuthenticated, init, login, logout }
}
