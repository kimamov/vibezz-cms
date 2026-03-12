package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DataObjectInterface provides dirty tracking capabilities for entities
type DataObjectInterface interface {
	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()
	Hydrate(data map[string]string)
	Get(key string) string
	Set(key string, value string)
	IsDirty() bool
}

// EntityInterface is the base interface for all entities
type EntityInterface interface {
	DataObjectInterface

	ID() uuid.UUID
	SetID(id uuid.UUID)

	CreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)

	UpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time)

	SoftDeletedAt() time.Time
	SetSoftDeletedAt(softDeletedAt time.Time)

	IsSoftDeleted() bool
}

// UserInterface represents a user entity
type UserInterface interface {
	EntityInterface

	Email() string
	SetEmail(email string)

	Name() string
	SetName(name string)

	PasswordHash() string
	SetPasswordHash(hash string)

	Role() UserRole
	SetRole(role UserRole)
}

// ContentTypeInterface represents a content type definition
type ContentTypeInterface interface {
	EntityInterface

	Name() string
	SetName(name string)

	Slug() string
	SetSlug(slug string)

	Fields() []FieldDefinition
	SetFields(fields []FieldDefinition) error
}

// FieldDefinition represents a single field in a content type
type FieldDefinition struct {
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	Type      string      `json:"type"`
	Required  bool        `json:"required"`
	Options   interface{} `json:"options,omitempty"`
	SortOrder int         `json:"sort_order"`
}

// EntryInterface represents a content entry
type EntryInterface interface {
	EntityInterface

	ContentTypeID() uuid.UUID
	SetContentTypeID(contentTypeID uuid.UUID)

	Title() string
	SetTitle(title string)

	Slug() string
	SetSlug(slug string)

	Path() string
	SetPath(path string)

	ParentID() *uuid.UUID
	SetParentID(parentID *uuid.UUID)

	AuthorID() uuid.UUID
	SetAuthorID(authorID uuid.UUID)

	Status() EntityStatus
	SetStatus(status EntityStatus)

	Fields() map[string]interface{}
	SetFields(fields map[string]interface{}) error

	PublishedAt() *time.Time
	SetPublishedAt(publishedAt *time.Time)

	IsPublished() bool
}

// EntryVersionInterface represents a version of an entry
type EntryVersionInterface interface {
	ID() uuid.UUID
	SetID(id uuid.UUID)

	EntryID() uuid.UUID
	SetEntryID(entryID uuid.UUID)

	Title() string
	SetTitle(title string)

	Fields() map[string]interface{}
	SetFields(fields map[string]interface{}) error

	AuthorID() uuid.UUID
	SetAuthorID(authorID uuid.UUID)

	CreatedAt() time.Time
	SetCreatedAt(createdAt time.Time)
}

// MediaFileInterface represents a media file entity
type MediaFileInterface interface {
	EntityInterface

	Filename() string
	SetFilename(filename string)

	MimeType() string
	SetMimeType(mimeType string)

	Size() int64
	SetSize(size int64)

	StoragePath() string
	SetStoragePath(path string)

	UploaderID() uuid.UUID
	SetUploaderID(uploaderID uuid.UUID)
}

// StoreInterface defines all operations available in the store
type StoreInterface interface {
	// User operations
	UserCreate(ctx context.Context, user UserInterface) error
	UserUpdate(ctx context.Context, user UserInterface) error
	UserDelete(ctx context.Context, id uuid.UUID) error
	UserSoftDelete(ctx context.Context, id uuid.UUID) error
	UserFindByID(ctx context.Context, id uuid.UUID) (UserInterface, error)
	UserFindByEmail(ctx context.Context, email string) (UserInterface, error)
	UserList(ctx context.Context, query QueryInterface) ([]UserInterface, error)
	UserCount(ctx context.Context, query QueryInterface) (int64, error)

	// ContentType operations
	ContentTypeCreate(ctx context.Context, ct ContentTypeInterface) error
	ContentTypeUpdate(ctx context.Context, ct ContentTypeInterface) error
	ContentTypeDelete(ctx context.Context, id uuid.UUID) error
	ContentTypeFindByID(ctx context.Context, id uuid.UUID) (ContentTypeInterface, error)
	ContentTypeFindBySlug(ctx context.Context, slug string) (ContentTypeInterface, error)
	ContentTypeList(ctx context.Context, query QueryInterface) ([]ContentTypeInterface, error)

	// Entry operations
	EntryCreate(ctx context.Context, entry EntryInterface) error
	EntryUpdate(ctx context.Context, entry EntryInterface) error
	EntryDelete(ctx context.Context, id uuid.UUID) error
	EntrySoftDelete(ctx context.Context, id uuid.UUID) error
	EntryFindByID(ctx context.Context, id uuid.UUID) (EntryInterface, error)
	EntryFindBySlug(ctx context.Context, slug string, contentTypeID uuid.UUID) (EntryInterface, error)
	EntryFindByPath(ctx context.Context, path string) (EntryInterface, error)
	EntryList(ctx context.Context, query QueryInterface) ([]EntryInterface, error)
	EntryCount(ctx context.Context, query QueryInterface) (int64, error)
	EntryListTree(ctx context.Context, contentTypeID uuid.UUID) ([]EntryInterface, error)

	// EntryVersion operations
	EntryVersionCreate(ctx context.Context, version EntryVersionInterface) error
	EntryVersionList(ctx context.Context, entryID uuid.UUID) ([]EntryVersionInterface, error)

	// MediaFile operations
	MediaFileCreate(ctx context.Context, media MediaFileInterface) error
	MediaFileDelete(ctx context.Context, id uuid.UUID) error
	MediaFileFindByID(ctx context.Context, id uuid.UUID) (MediaFileInterface, error)
	MediaFileList(ctx context.Context, query QueryInterface) ([]MediaFileInterface, error)
}

// QueryInterface is the base interface for all query builders
type QueryInterface interface {
	Validate() error

	HasLimit() bool
	Limit() int
	SetLimit(limit int) QueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) QueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) QueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) QueryInterface
}

// EntryQueryInterface extends QueryInterface for entry-specific queries
type EntryQueryInterface interface {
	QueryInterface

	HasContentTypeID() bool
	ContentTypeID() uuid.UUID
	SetContentTypeID(contentTypeID uuid.UUID) EntryQueryInterface

	HasParentID() bool
	ParentID() *uuid.UUID
	SetParentID(parentID *uuid.UUID) EntryQueryInterface

	HasStatus() bool
	Status() EntityStatus
	SetStatus(status EntityStatus) EntryQueryInterface

	HasSoftDeleteIncluded() bool
	SoftDeleteIncluded() bool
	SetSoftDeleteIncluded(included bool) EntryQueryInterface
}

// UserQueryInterface extends QueryInterface for user-specific queries
type UserQueryInterface interface {
	QueryInterface

	HasRole() bool
	Role() UserRole
	SetRole(role UserRole) UserQueryInterface

	HasEmail() bool
	Email() string
	SetEmail(email string) UserQueryInterface
}

// BlockInterface represents a content block that can be rendered
type BlockInterface interface {
	EntityInterface

	// BlockType returns the type of block (e.g., "heading", "container", "grid")
	BlockType() string
	SetBlockType(blockType string)

	// EntryID returns the entry this block belongs to
	EntryID() uuid.UUID
	SetEntryID(entryID uuid.UUID)

	// ParentID returns the parent block ID (for nested blocks)
	ParentID() *uuid.UUID
	SetParentID(parentID *uuid.UUID)

	// Sequence is the order of this block within its container
	Sequence() int
	SetSequence(sequence int)

	// BlockData contains block-specific data (e.g., text content, image URL, grid config)
	BlockData() map[string]interface{}
	SetBlockData(data map[string]interface{}) error

	// Children returns child blocks (for container blocks)
	Children() []BlockInterface
	SetChildren(children []BlockInterface)

	// Style contains CSS classes and inline styles
	Style() BlockStyle
	SetStyle(style BlockStyle)

	// IsContainer returns true if this block can contain other blocks
	IsContainer() bool

	// Clone creates a deep copy of the block
	Clone() BlockInterface
}

// BlockStyle represents styling options for a block
type BlockStyle struct {
	ClassName   string                 `json:"className,omitempty"`
	InlineStyle map[string]string      `json:"inlineStyle,omitempty"`
	Responsive  *ResponsiveBreakpoints `json:"responsive,omitempty"`
}

// ResponsiveBreakpoints contains responsive breakpoint styles
type ResponsiveBreakpoints struct {
	Mobile  *BlockStyle `json:"mobile,omitempty"`
	Tablet  *BlockStyle `json:"tablet,omitempty"`
	Desktop *BlockStyle `json:"desktop,omitempty"`
}

// ContainerBlockInterface extends BlockInterface for container blocks
type ContainerBlockInterface interface {
	BlockInterface

	// Layout returns the layout type (flex, grid, stack)
	Layout() string
	SetLayout(layout string)

	// Gap returns the spacing between child blocks
	Gap() string
	SetGap(gap string)

	// Padding returns the padding config
	Padding() string
	SetPadding(padding string)

	// Background returns background color/image
	Background() string
	SetBackground(background string)

	// MaxWidth constrains the container width
	MaxWidth() string
	SetMaxWidth(maxWidth string)
}

// GridBlockInterface extends BlockInterface for grid layouts
type GridBlockInterface interface {
	BlockInterface

	// Columns returns the number of columns (1-12)
	Columns() int
	SetColumns(columns int)

	// ColumnGap returns the gap between columns
	ColumnGap() string
	SetColumnGap(gap string)

	// RowGap returns the gap between rows
	RowGap() string
	SetRowGap(gap string)

	// ResponsiveColumns returns columns per breakpoint
	ResponsiveColumns() ResponsiveColumns
	SetResponsiveColumns(cols ResponsiveColumns)
}

// ResponsiveColumns defines columns at different breakpoints
type ResponsiveColumns struct {
	Mobile  int `json:"mobile"`
	Tablet  int `json:"tablet"`
	Desktop int `json:"desktop"`
}

// BlockQueryInterface for querying blocks
type BlockQueryInterface interface {
	QueryInterface

	HasEntryID() bool
	EntryID() uuid.UUID
	SetEntryID(entryID uuid.UUID) BlockQueryInterface

	HasParentID() bool
	ParentID() *uuid.UUID
	SetParentID(parentID *uuid.UUID) BlockQueryInterface

	HasBlockType() bool
	BlockType() string
	SetBlockType(blockType string) BlockQueryInterface
}
