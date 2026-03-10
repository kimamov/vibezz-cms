import type { Entry, NavigationItem } from '@vibezz/types'

export function useCms() {
  const config = useRuntimeConfig()

  async function fetchPage(path: string): Promise<Entry | null> {
    try {
      const response = await fetch(`${config.public.apiBase}/api/public/routes${path}`)
      if (!response.ok) return null
      return response.json()
    } catch {
      return null
    }
  }

  async function fetchNavigation(): Promise<NavigationItem[]> {
    try {
      const response = await fetch(`${config.public.apiBase}/api/public/navigation`)
      if (!response.ok) return []
      return response.json()
    } catch {
      return []
    }
  }

  function mediaUrl(id: string): string {
    return `${config.public.apiBase}/api/public/media/${id}`
  }

  return { fetchPage, fetchNavigation, mediaUrl }
}
