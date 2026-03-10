export interface User {
  id: string
  email: string
  name: string
  role: 'admin' | 'editor' | 'author'
  created_at: string
  updated_at: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
}

export interface ContentType {
  id: string
  name: string
  slug: string
  fields: FieldDefinition[]
  created_at: string
  updated_at: string
}

export interface FieldDefinition {
  name: string
  slug: string
  type: FieldType
  required: boolean
  options?: Record<string, unknown>
  sort_order: number
}

export type FieldType =
  | 'text'
  | 'textarea'
  | 'richtext'
  | 'number'
  | 'boolean'
  | 'date'
  | 'select'
  | 'media'
  | 'url'
  | 'email'
  | 'json'

export interface Entry {
  id: string
  content_type_id: string
  title: string
  slug: string
  path: string
  parent_id?: string
  author_id: string
  status: 'draft' | 'published' | 'archived'
  fields: Record<string, unknown>
  published_at?: string
  created_at: string
  updated_at: string
}

export interface CreateEntryInput {
  content_type_id: string
  title: string
  slug: string
  parent_id?: string
  fields: Record<string, unknown>
}

export interface UpdateEntryInput {
  title?: string
  slug?: string
  fields?: Record<string, unknown>
}

export interface NavigationItem {
  id: string
  title: string
  slug: string
  path: string
  children?: NavigationItem[]
}

export interface MediaFile {
  id: string
  filename: string
  mime_type: string
  size: number
  storage_path: string
  uploader_id: string
  created_at: string
}

export interface ApiError {
  error: string
  details?: string
}
