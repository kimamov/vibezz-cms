package domain

import (
	"github.com/google/uuid"
)

// entryQuery is the implementation of EntryQueryInterface
type entryQuery struct {
	*query
}

// NewEntryQuery creates a new EntryQuery instance
func NewEntryQuery() EntryQueryInterface {
	return &entryQuery{
		query: NewQuery().(*query),
	}
}

// HasContentTypeID returns true if content type ID is set
func (q *entryQuery) HasContentTypeID() bool {
	_, ok := q.properties["content_type_id"]
	return ok
}

// ContentTypeID returns the content type ID
func (q *entryQuery) ContentTypeID() uuid.UUID {
	if val, ok := q.properties["content_type_id"].(uuid.UUID); ok {
		return val
	}
	return uuid.Nil
}

// SetContentTypeID sets the content type ID
func (q *entryQuery) SetContentTypeID(contentTypeID uuid.UUID) EntryQueryInterface {
	q.properties["content_type_id"] = contentTypeID
	return q
}

// HasParentID returns true if parent ID is set
func (q *entryQuery) HasParentID() bool {
	_, ok := q.properties["parent_id"]
	return ok
}

// ParentID returns the parent ID
func (q *entryQuery) ParentID() *uuid.UUID {
	if val, ok := q.properties["parent_id"].(*uuid.UUID); ok {
		return val
	}
	return nil
}

// SetParentID sets the parent ID
func (q *entryQuery) SetParentID(parentID *uuid.UUID) EntryQueryInterface {
	q.properties["parent_id"] = parentID
	return q
}

// HasStatus returns true if status is set
func (q *entryQuery) HasStatus() bool {
	_, ok := q.properties["status"]
	return ok
}

// Status returns the status
func (q *entryQuery) Status() EntityStatus {
	if val, ok := q.properties["status"].(EntityStatus); ok {
		return val
	}
	return ""
}

// SetStatus sets the status
func (q *entryQuery) SetStatus(status EntityStatus) EntryQueryInterface {
	q.properties["status"] = status
	return q
}

// HasSoftDeleteIncluded returns true if soft delete inclusion is set
func (q *entryQuery) HasSoftDeleteIncluded() bool {
	_, ok := q.properties["soft_delete_included"]
	return ok
}

// SoftDeleteIncluded returns whether to include soft deleted records
func (q *entryQuery) SoftDeleteIncluded() bool {
	if val, ok := q.properties["soft_delete_included"].(bool); ok {
		return val
	}
	return false
}

// SetSoftDeleteIncluded sets whether to include soft deleted records
func (q *entryQuery) SetSoftDeleteIncluded(included bool) EntryQueryInterface {
	q.properties["soft_delete_included"] = included
	return q
}
