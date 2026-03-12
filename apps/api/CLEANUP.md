# Vibezz CMS API Cleanup Summary

## Overview
This cleanup refactors the `apps/api` codebase to follow the cleaner architectural patterns from the `cmsstore` reference implementation. The goal is to establish a more maintainable, testable, and scalable foundation.

## Key Patterns Applied from cmsstore

### 1. Interface-First Design
All entities are now defined by interfaces first, located in `internal/domain/interfaces.go`:
- `UserInterface`
- `ContentTypeInterface`
- `EntryInterface`
- `EntryVersionInterface`
- `MediaFileInterface`
- `StoreInterface`

This enables:
- Dependency inversion
- Easy mocking for tests
- Clean separation of concerns
- Multiple implementations if needed

### 2. DataObject for Dirty Tracking
Created `internal/domain/dataobject.go` with a base `DataObject` struct that provides:
- Automatic dirty tracking (tracks which fields changed)
- `Data()` - returns all data
- `DataChanged()` - returns only modified fields
- `IsDirty()` - check if any field was modified
- `MarkAsNotDirty()` - reset dirty state
- `Hydrate()` - populate from existing data

All entities embed this struct for automatic change tracking.

### 3. Consistent Entity Structure
Each entity follows a consistent pattern:

```
internal/domain/
├── interfaces.go         # Interface definitions
├── consts.go             # Constants (column names, statuses)
├── dataobject.go         # Base dirty tracking struct
├── user.go              # User entity implementation
├── content_type.go      # ContentType entity implementation
├── entry.go             # Entry entity implementation
├── entry_version.go     # EntryVersion entity implementation
├── media_file.go        # MediaFile entity implementation
├── query.go             # Base query builder
├── entry_query.go       # Entry-specific query builder
└── user_query.go        # User-specific query builder
```

### 4. Fluent Setters (Modified for Go)
Unlike cmsstore which returns interface types (fluent chaining), we use void return types to avoid Go's covariance issues. This is still clean and readable:

```go
user := domain.NewUser()
user.SetName("John")
user.SetEmail("john@example.com")
user.SetRole(domain.RoleAdmin)
```

### 5. Query Builder Pattern
Created flexible query builders in:
- `internal/domain/query.go` - Base query with limit, offset, order by
- `internal/domain/entry_query.go` - Entry-specific filters (content_type_id, status, parent_id)
- `internal/domain/user_query.go` - User-specific filters (role, email)

Usage:
```go
query := domain.NewEntryQuery().
    SetContentTypeID(contentTypeID).
    SetStatus(domain.StatusActive).
    SetLimit(10)

entries, err := store.EntryList(ctx, query)
```

### 6. Store Pattern
Created `internal/store/store.go` as a bridge between domain interfaces and the existing sqlc-generated database layer:

```go
store := store.NewStore(pool)

// Create
user := domain.NewUser()
user.SetName("John")
user.SetEmail("john@example.com")
err := store.UserCreate(ctx, user)

// Read
user, err := store.UserFindByID(ctx, userID)

// Update
user.SetName("Jane")
err := store.UserUpdate(ctx, user)

// List
query := domain.NewUserQuery().SetLimit(10)
users, err := store.UserList(ctx, query)
```

## Files Created

### New Domain Package (`internal/domain/`)
1. **consts.go** - Constants for statuses, column names, soft delete handling
2. **interfaces.go** - All entity and store interfaces
3. **dataobject.go** - Base DataObject struct with dirty tracking
4. **user.go** - User entity implementation
5. **content_type.go** - ContentType entity implementation
6. **entry.go** - Entry entity implementation
7. **entry_version.go** - EntryVersion entity implementation
8. **media_file.go** - MediaFile entity implementation
9. **query.go** - Base query builder
10. **entry_query.go** - Entry-specific query builder
11. **user_query.go** - User-specific query builder

### New Store Package (`internal/store/`)
1. **store.go** - Store implementation bridging domain to sqlc

## Migration Path

### Phase 1: Add Missing SQL Queries (Next Step)
The following database operations need SQL query definitions:

**Users:**
- `UpdateUser` - Update user details
- `DeleteUser` - Delete user
- `ListUsers` - List users with pagination
- `CountUsers` - Count users

**Entries:**
- `GetEntryBySlug` - Find entry by slug within content type
- `CountEntries` - Count entries
- `ListEntryTree` - Get hierarchical entry tree

**Entry Versions:**
- `CreateEntryVersion` - Create version snapshot
- `ListEntryVersions` - List versions for an entry

**Media Files:**
- Already implemented

**Soft Deletes:**
- Add `soft_deleted_at` column to all tables
- Create soft delete queries for each entity

### Phase 2: Update Existing Services
Gradually migrate existing services from `internal/content/` to use the new domain types:

**Before:**
```go
// internal/content/user_service.go
type UserService struct {
    pool *pgxpool.Pool
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (*User, error) {
    // Direct SQL
}
```

**After:**
```go
func (s *UserService) Create(ctx context.Context, input CreateUserInput) (domain.UserInterface, error) {
    user := domain.NewUser()
    user.SetName(input.Name)
    user.SetEmail(input.Email)
    // ... set other fields
    
    err := s.store.UserCreate(ctx, user)
    return user, err
}
```

### Phase 3: Update Handlers
Update HTTP handlers to use domain types:

```go
// Handler creates a user
func (h *UserHandler) Create(c *gin.Context) {
    var req createUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    user, err := h.service.Create(c.Request.Context(), req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(201, user) // domain.UserInterface automatically JSON-serializable
}
```

### Phase 4: Remove Duplicate Models
Once migration is complete:
1. Remove `internal/content/models.go`
2. Remove `internal/db/models.go` (sqlc-generated, keep the queries)
3. Update sqlc config to generate only query methods, not models

## Benefits of This Cleanup

1. **Better Testability**: Interfaces enable easy mocking
2. **Clear Contracts**: Interfaces define exactly what each entity can do
3. **Dirty Tracking**: Automatic change tracking for partial updates
4. **Type Safety**: Less `any` types, more structured data
5. **Consistency**: All entities follow the same pattern
6. **Separation of Concerns**: Domain logic separate from database logic
7. **Easier Maintenance**: Changes to one layer don't cascade

## Quick Start for New Features

To add a new entity (e.g., `Tag`):

1. **Add to interfaces.go:**
```go
type TagInterface interface {
    EntityInterface
    Name() string
    SetName(name string)
    Slug() string
    SetSlug(slug string)
}
```

2. **Create tag.go:**
```go
type tag struct {
    *DataObject
}

func NewTag() TagInterface {
    t := &tag{DataObject: NewDataObject()}
    t.SetID(uuid.New())
    // ... set defaults
    return t
}
// ... implement interface methods
```

3. **Add to StoreInterface:**
```go
TagCreate(ctx context.Context, tag TagInterface) error
TagFindByID(ctx context.Context, id uuid.UUID) (TagInterface, error)
// ... etc
```

4. **Implement in store.go:**
```go
func (s *store) TagCreate(ctx context.Context, tag domain.TagInterface) error {
    // Convert domain to db params and call queries.CreateTag
}
```

## Notes

- The old `internal/content/models.go` and `internal/db/models.go` still exist for backward compatibility
- Services can gradually migrate to the new domain types
- The store bridges old and new during migration
- Once fully migrated, the duplicate models can be removed
