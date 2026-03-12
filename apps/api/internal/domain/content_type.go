package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Ensure contentType implements ContentTypeInterface
var _ ContentTypeInterface = (*contentType)(nil)

// contentType is the concrete implementation of ContentTypeInterface
type contentType struct {
	*DataObject
}

// NewContentType creates a new ContentType instance with default values
func NewContentType() ContentTypeInterface {
	ct := &contentType{
		DataObject: NewDataObject(),
	}
	ct.SetID(uuid.New())
	ct.SetCreatedAt(time.Now().UTC())
	ct.SetUpdatedAt(time.Now().UTC())
	ct.SetSoftDeletedAt(SoftDeleteTime)
	ct.SetFields([]FieldDefinition{})
	return ct
}

// NewContentTypeFromData creates a ContentType from existing data
func NewContentTypeFromData(data map[string]string) ContentTypeInterface {
	ct := &contentType{
		DataObject: NewDataObject(),
	}
	ct.Hydrate(data)
	return ct
}

// ID returns the content type ID
func (ct *contentType) ID() uuid.UUID {
	return uuid.MustParse(ct.Get(ColumnID))
}

// SetID sets the content type ID
func (ct *contentType) SetID(id uuid.UUID) {
	ct.Set(ColumnID, id.String())
}

// CreatedAt returns the creation time
func (ct *contentType) CreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, ct.Get(ColumnCreatedAt))
	return t
}

// SetCreatedAt sets the creation time
func (ct *contentType) SetCreatedAt(createdAt time.Time) {
	ct.Set(ColumnCreatedAt, createdAt.Format(time.RFC3339))
}

// UpdatedAt returns the last update time
func (ct *contentType) UpdatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, ct.Get(ColumnUpdatedAt))
	return t
}

// SetUpdatedAt sets the last update time
func (ct *contentType) SetUpdatedAt(updatedAt time.Time) {
	ct.Set(ColumnUpdatedAt, updatedAt.Format(time.RFC3339))
}

// SoftDeletedAt returns the soft delete timestamp
func (ct *contentType) SoftDeletedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, ct.Get(ColumnSoftDeletedAt))
	return t
}

// SetSoftDeletedAt sets the soft delete timestamp
func (ct *contentType) SetSoftDeletedAt(softDeletedAt time.Time) {
	ct.Set(ColumnSoftDeletedAt, softDeletedAt.Format(time.RFC3339))
}

// IsSoftDeleted returns true if the content type is soft deleted
func (ct *contentType) IsSoftDeleted() bool {
	return IsSoftDeleted(ct.SoftDeletedAt())
}

// Name returns the content type name
func (ct *contentType) Name() string {
	return ct.Get(ColumnName)
}

// SetName sets the content type name
func (ct *contentType) SetName(name string) {
	ct.Set(ColumnName, name)
}

// Slug returns the content type slug
func (ct *contentType) Slug() string {
	return ct.Get(ColumnSlug)
}

// SetSlug sets the content type slug
func (ct *contentType) SetSlug(slug string) {
	ct.Set(ColumnSlug, slug)
}

// Fields returns the content type fields
func (ct *contentType) Fields() []FieldDefinition {
	fieldsJSON := ct.Get(ColumnFields)
	if fieldsJSON == "" {
		return []FieldDefinition{}
	}

	var fields []FieldDefinition
	if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
		return []FieldDefinition{}
	}
	return fields
}

// SetFields sets the content type fields
func (ct *contentType) SetFields(fields []FieldDefinition) error {
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	ct.Set(ColumnFields, string(fieldsJSON))
	return nil
}
