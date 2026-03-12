package domain

import (
	"time"

	"github.com/google/uuid"
)

// entryVersion is the concrete implementation of EntryVersionInterface
type entryVersion struct {
	id        uuid.UUID
	entryID   uuid.UUID
	title     string
	fields    map[string]interface{}
	authorID  uuid.UUID
	createdAt time.Time
}

// NewEntryVersion creates a new EntryVersion instance
func NewEntryVersion() EntryVersionInterface {
	return &entryVersion{
		id:        uuid.New(),
		createdAt: time.Now().UTC(),
		fields:    make(map[string]interface{}),
	}
}

// ID returns the version ID
func (v *entryVersion) ID() uuid.UUID {
	return v.id
}

// SetID sets the version ID
func (v *entryVersion) SetID(id uuid.UUID) {
	v.id = id
}

// EntryID returns the entry ID
func (v *entryVersion) EntryID() uuid.UUID {
	return v.entryID
}

// SetEntryID sets the entry ID
func (v *entryVersion) SetEntryID(entryID uuid.UUID) {
	v.entryID = entryID
}

// Title returns the title
func (v *entryVersion) Title() string {
	return v.title
}

// SetTitle sets the title
func (v *entryVersion) SetTitle(title string) {
	v.title = title
}

// Fields returns the fields
func (v *entryVersion) Fields() map[string]interface{} {
	return v.fields
}

// SetFields sets the fields
func (v *entryVersion) SetFields(fields map[string]interface{}) error {
	v.fields = fields
	return nil
}

// AuthorID returns the author ID
func (v *entryVersion) AuthorID() uuid.UUID {
	return v.authorID
}

// SetAuthorID sets the author ID
func (v *entryVersion) SetAuthorID(authorID uuid.UUID) {
	v.authorID = authorID
}

// CreatedAt returns the creation time
func (v *entryVersion) CreatedAt() time.Time {
	return v.createdAt
}

// SetCreatedAt sets the creation time
func (v *entryVersion) SetCreatedAt(createdAt time.Time) {
	v.createdAt = createdAt
}
