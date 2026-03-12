package domain

import (
	"time"
)

// EntityStatus represents the status of any entity
type EntityStatus string

const (
	StatusDraft    EntityStatus = "draft"
	StatusActive   EntityStatus = "active"
	StatusInactive EntityStatus = "inactive"
)

// ColumnNames - standard column names for database queries
const (
	ColumnID            = "id"
	ColumnCreatedAt     = "created_at"
	ColumnUpdatedAt     = "updated_at"
	ColumnSoftDeletedAt = "soft_deleted_at"
	ColumnStatus        = "status"
	ColumnName          = "name"
	ColumnSlug          = "slug"
	ColumnTitle         = "title"
	ColumnContent       = "content"
	ColumnFields        = "fields"
	ColumnMetas         = "metas"
	ColumnAuthorID      = "author_id"
	ColumnParentID      = "parent_id"
	ColumnPath          = "path"
	ColumnPublishedAt   = "published_at"
)

// SoftDeleteTime is the maximum datetime for non-deleted records
var SoftDeleteTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

// IsSoftDeleted checks if a timestamp represents a soft-deleted record
func IsSoftDeleted(t time.Time) bool {
	return t.Before(SoftDeleteTime)
}

// UserRole represents user roles in the system
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleEditor UserRole = "editor"
	RoleAuthor UserRole = "author"
)

// ContentTypeFieldType represents the type of a content type field
type ContentTypeFieldType string

const (
	FieldTypeText     ContentTypeFieldType = "text"
	FieldTypeNumber   ContentTypeFieldType = "number"
	FieldTypeBoolean  ContentTypeFieldType = "boolean"
	FieldTypeDate     ContentTypeFieldType = "date"
	FieldTypeMedia    ContentTypeFieldType = "media"
	FieldTypeRelation ContentTypeFieldType = "relation"
	FieldTypeBlocks   ContentTypeFieldType = "blocks"
)
