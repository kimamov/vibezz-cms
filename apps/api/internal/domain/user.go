package domain

import (
	"time"

	"github.com/google/uuid"
)

// Ensure user implements UserInterface
var _ UserInterface = (*user)(nil)

// user is the concrete implementation of UserInterface
type user struct {
	*DataObject
}

// NewUser creates a new User instance with default values
func NewUser() UserInterface {
	u := &user{
		DataObject: NewDataObject(),
	}
	u.SetID(uuid.New())
	u.SetRole(RoleAuthor)
	u.SetCreatedAt(time.Now().UTC())
	u.SetUpdatedAt(time.Now().UTC())
	u.SetSoftDeletedAt(SoftDeleteTime)
	return u
}

// NewUserFromData creates a User from existing data
func NewUserFromData(data map[string]string) UserInterface {
	u := &user{
		DataObject: NewDataObject(),
	}
	u.Hydrate(data)
	return u
}

// ID returns the user ID
func (u *user) ID() uuid.UUID {
	return uuid.MustParse(u.Get(ColumnID))
}

// SetID sets the user ID
func (u *user) SetID(id uuid.UUID) {
	u.Set(ColumnID, id.String())
}

// CreatedAt returns the creation time
func (u *user) CreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, u.Get(ColumnCreatedAt))
	return t
}

// SetCreatedAt sets the creation time
func (u *user) SetCreatedAt(createdAt time.Time) {
	u.Set(ColumnCreatedAt, createdAt.Format(time.RFC3339))
}

// UpdatedAt returns the last update time
func (u *user) UpdatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, u.Get(ColumnUpdatedAt))
	return t
}

// SetUpdatedAt sets the last update time
func (u *user) SetUpdatedAt(updatedAt time.Time) {
	u.Set(ColumnUpdatedAt, updatedAt.Format(time.RFC3339))
}

// SoftDeletedAt returns the soft delete timestamp
func (u *user) SoftDeletedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, u.Get(ColumnSoftDeletedAt))
	return t
}

// SetSoftDeletedAt sets the soft delete timestamp
func (u *user) SetSoftDeletedAt(softDeletedAt time.Time) {
	u.Set(ColumnSoftDeletedAt, softDeletedAt.Format(time.RFC3339))
}

// IsSoftDeleted returns true if the user is soft deleted
func (u *user) IsSoftDeleted() bool {
	return IsSoftDeleted(u.SoftDeletedAt())
}

// Email returns the user's email
func (u *user) Email() string {
	return u.Get("email")
}

// SetEmail sets the user's email
func (u *user) SetEmail(email string) {
	u.Set("email", email)
}

// Name returns the user's name
func (u *user) Name() string {
	return u.Get(ColumnName)
}

// SetName sets the user's name
func (u *user) SetName(name string) {
	u.Set(ColumnName, name)
}

// PasswordHash returns the password hash
func (u *user) PasswordHash() string {
	return u.Get("password_hash")
}

// SetPasswordHash sets the password hash
func (u *user) SetPasswordHash(hash string) {
	u.Set("password_hash", hash)
}

// Role returns the user's role
func (u *user) Role() UserRole {
	return UserRole(u.Get("role"))
}

// SetRole sets the user's role
func (u *user) SetRole(role UserRole) {
	u.Set("role", string(role))
}
