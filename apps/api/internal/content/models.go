package content

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	Name         string     `json:"name"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type ContentType struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	Slug      string            `json:"slug"`
	Fields    []FieldDefinition `json:"fields"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type FieldDefinition struct {
	Name       string      `json:"name"`
	Slug       string      `json:"slug"`
	Type       string      `json:"type"`
	Required   bool        `json:"required"`
	Options    interface{} `json:"options,omitempty"`
	SortOrder  int         `json:"sort_order"`
}

type FieldDefinitionInput struct {
	Name     string      `json:"name" binding:"required"`
	Slug     string      `json:"slug" binding:"required"`
	Type     string      `json:"type" binding:"required"`
	Required bool        `json:"required"`
	Options  interface{} `json:"options,omitempty"`
}

type Entry struct {
	ID            uuid.UUID              `json:"id"`
	ContentTypeID uuid.UUID              `json:"content_type_id"`
	Title         string                 `json:"title"`
	Slug          string                 `json:"slug"`
	Path          string                 `json:"path"`
	ParentID      *uuid.UUID             `json:"parent_id,omitempty"`
	AuthorID      uuid.UUID              `json:"author_id"`
	Status        string                 `json:"status"`
	Fields        map[string]interface{} `json:"fields"`
	PublishedAt   *time.Time             `json:"published_at,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type CreateEntryInput struct {
	ContentTypeID uuid.UUID              `json:"content_type_id"`
	Title         string                 `json:"title"`
	Slug          string                 `json:"slug"`
	ParentID      *uuid.UUID             `json:"parent_id,omitempty"`
	Fields        map[string]interface{} `json:"fields"`
	AuthorID      uuid.UUID              `json:"author_id"`
}

type UpdateEntryInput struct {
	Title  *string                 `json:"title,omitempty"`
	Slug   *string                 `json:"slug,omitempty"`
	Fields *map[string]interface{} `json:"fields,omitempty"`
}

type EntryFilters struct {
	ContentTypeID *uuid.UUID
	ParentID      *uuid.UUID
	Status        string
}

type NavigationItem struct {
	ID       uuid.UUID        `json:"id"`
	Title    string           `json:"title"`
	Slug     string           `json:"slug"`
	Path     string           `json:"path"`
	Children []NavigationItem `json:"children,omitempty"`
}

type MediaFile struct {
	ID          uuid.UUID  `json:"id"`
	Filename    string     `json:"filename"`
	MimeType    string     `json:"mime_type"`
	Size        int64      `json:"size"`
	StoragePath string     `json:"storage_path"`
	UploaderID  uuid.UUID  `json:"uploader_id"`
	CreatedAt   time.Time  `json:"created_at"`
}
