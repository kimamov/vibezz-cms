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
  | 'blocks'

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

export interface MediaFieldValue {
  id: string
  filename: string
  mime_type: string
  size: number
  url: string
}

export function isMediaField(value: unknown): value is MediaFieldValue {
  return typeof value === 'object' && value !== null && 'url' in value && 'mime_type' in value
}

export type BlockType = 'heading' | 'text' | 'image' | 'quote' | 'divider' | 'code' | (string & {})

export interface BlockConfigField {
  name: string
  slug: string
  type: string
  required: boolean
}

export interface BlockTypeDefinition {
  slug: string
  label: string
  icon: string
  config_fields?: BlockConfigField[]
  default_data: Record<string, unknown>
}

export interface ContentBlock {
  id: string
  type: string
  data: Record<string, unknown>
}

export interface HeadingBlockData {
  text: string
  level: 1 | 2 | 3 | 4 | 5 | 6
}

export interface TextBlockData {
  content: string
}

export interface ImageBlockData {
  media_id?: string
  media?: MediaFieldValue
  caption?: string
  alt?: string
}

export interface QuoteBlockData {
  text: string
  attribution?: string
}

export interface CodeBlockData {
  code: string
  language?: string
}

export function isBlockArray(value: unknown): value is ContentBlock[] {
  return Array.isArray(value) && value.every(
    item => typeof item === 'object' && item !== null && 'type' in item && 'data' in item,
  )
}

export interface ApiError {
  error: string
  details?: string
}
