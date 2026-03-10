import type { ApiError } from '@vibezz/types'

export function useApi() {
  const config = useRuntimeConfig()
  const auth = useAuth()

  async function apiFetch<T>(
    path: string,
    options: RequestInit = {},
  ): Promise<T> {
    const headers: Record<string, string> = {
      ...(options.headers as Record<string, string>),
    }

    if (auth.token.value) {
      headers['Authorization'] = `Bearer ${auth.token.value}`
    }

    if (!(options.body instanceof FormData)) {
      headers['Content-Type'] = 'application/json'
    }

    const response = await fetch(`${config.public.apiBase}${path}`, {
      ...options,
      headers,
    })

    if (!response.ok) {
      const error: ApiError = await response.json().catch(() => ({
        error: 'Request failed',
      }))
      throw new Error(error.error)
    }

    if (response.status === 204) {
      return undefined as T
    }

    return response.json()
  }

  return { apiFetch }
}
