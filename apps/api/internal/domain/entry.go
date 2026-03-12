package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Ensure entry implements EntryInterface
var _ EntryInterface = (*entry)(nil)

// entry is the concrete implementation of EntryInterface
type entry struct {
	*DataObject
}

// NewEntry creates a new Entry instance with default values
func NewEntry() EntryInterface {
	e := &entry{
		DataObject: NewDataObject(),
	}
	e.SetID(uuid.New())
	e.SetStatus(StatusDraft)
	e.SetCreatedAt(time.Now().UTC())
	e.SetUpdatedAt(time.Now().UTC())
	e.SetSoftDeletedAt(SoftDeleteTime)
	e.SetFields(map[string]interface{}{})
	return e
}

// NewEntryFromData creates an Entry from existing data
func NewEntryFromData(data map[string]string) EntryInterface {
	e := &entry{
		DataObject: NewDataObject(),
	}
	e.Hydrate(data)
	return e
}

// ID returns the entry ID
func (e *entry) ID() uuid.UUID {
	return uuid.MustParse(e.Get(ColumnID))
}

// SetID sets the entry ID
func (e *entry) SetID(id uuid.UUID) {
	e.Set(ColumnID, id.String())
}

// CreatedAt returns the creation time
func (e *entry) CreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, e.Get(ColumnCreatedAt))
	return t
}

// SetCreatedAt sets the creation time
func (e *entry) SetCreatedAt(createdAt time.Time) {
	e.Set(ColumnCreatedAt, createdAt.Format(time.RFC3339))
}

// UpdatedAt returns the last update time
func (e *entry) UpdatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, e.Get(ColumnUpdatedAt))
	return t
}

// SetUpdatedAt sets the last update time
func (e *entry) SetUpdatedAt(updatedAt time.Time) {
	e.Set(ColumnUpdatedAt, updatedAt.Format(time.RFC3339))
}

// SoftDeletedAt returns the soft delete timestamp
func (e *entry) SoftDeletedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, e.Get(ColumnSoftDeletedAt))
	return t
}

// SetSoftDeletedAt sets the soft delete timestamp
func (e *entry) SetSoftDeletedAt(softDeletedAt time.Time) {
	e.Set(ColumnSoftDeletedAt, softDeletedAt.Format(time.RFC3339))
}

// IsSoftDeleted returns true if the entry is soft deleted
func (e *entry) IsSoftDeleted() bool {
	return IsSoftDeleted(e.SoftDeletedAt())
}

// ContentTypeID returns the content type ID
func (e *entry) ContentTypeID() uuid.UUID {
	return uuid.MustParse(e.Get("content_type_id"))
}

// SetContentTypeID sets the content type ID
func (e *entry) SetContentTypeID(contentTypeID uuid.UUID) {
	e.Set("content_type_id", contentTypeID.String())
}

// Title returns the entry title
func (e *entry) Title() string {
	return e.Get(ColumnTitle)
}

// SetTitle sets the entry title
func (e *entry) SetTitle(title string) {
	e.Set(ColumnTitle, title)
}

// Slug returns the entry slug
func (e *entry) Slug() string {
	return e.Get(ColumnSlug)
}

// SetSlug sets the entry slug
func (e *entry) SetSlug(slug string) {
	e.Set(ColumnSlug, slug)
}

// Path returns the entry path
func (e *entry) Path() string {
	return e.Get(ColumnPath)
}

// SetPath sets the entry path
func (e *entry) SetPath(path string) {
	e.Set(ColumnPath, path)
}

// ParentID returns the parent entry ID
func (e *entry) ParentID() *uuid.UUID {
	parentIDStr := e.Get(ColumnParentID)
	if parentIDStr == "" {
		return nil
	}
	id := uuid.MustParse(parentIDStr)
	return &id
}

// SetParentID sets the parent entry ID
func (e *entry) SetParentID(parentID *uuid.UUID) {
	if parentID == nil {
		e.Set(ColumnParentID, "")
	} else {
		e.Set(ColumnParentID, parentID.String())
	}
}

// AuthorID returns the author ID
func (e *entry) AuthorID() uuid.UUID {
	return uuid.MustParse(e.Get(ColumnAuthorID))
}

// SetAuthorID sets the author ID
func (e *entry) SetAuthorID(authorID uuid.UUID) {
	e.Set(ColumnAuthorID, authorID.String())
}

// Status returns the entry status
func (e *entry) Status() EntityStatus {
	return EntityStatus(e.Get(ColumnStatus))
}

// SetStatus sets the entry status
func (e *entry) SetStatus(status EntityStatus) {
	e.Set(ColumnStatus, string(status))
}

// Fields returns the entry fields
func (e *entry) Fields() map[string]interface{} {
	fieldsJSON := e.Get(ColumnFields)
	if fieldsJSON == "" {
		return map[string]interface{}{}
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
		return map[string]interface{}{}
	}
	return fields
}

// SetFields sets the entry fields
func (e *entry) SetFields(fields map[string]interface{}) error {
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	e.Set(ColumnFields, string(fieldsJSON))
	return nil
}

// PublishedAt returns the published at time
func (e *entry) PublishedAt() *time.Time {
	publishedAtStr := e.Get(ColumnPublishedAt)
	if publishedAtStr == "" {
		return nil
	}
	t, _ := time.Parse(time.RFC3339, publishedAtStr)
	return &t
}

// SetPublishedAt sets the published at time
func (e *entry) SetPublishedAt(publishedAt *time.Time) {
	if publishedAt == nil {
		e.Set(ColumnPublishedAt, "")
	} else {
		e.Set(ColumnPublishedAt, publishedAt.Format(time.RFC3339))
	}
}

// IsPublished returns true if the entry is published
func (e *entry) IsPublished() bool {
	return e.Status() == StatusActive && e.PublishedAt() != nil
}
