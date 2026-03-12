package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vibezz/cms/internal/db"
	"github.com/vibezz/cms/internal/domain"
)

// Ensure store implements StoreInterface
var _ domain.StoreInterface = (*store)(nil)

// store is the implementation of domain.StoreInterface
type store struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewStore creates a new store instance
func NewStore(pool *pgxpool.Pool) domain.StoreInterface {
	return &store{
		pool:    pool,
		queries: db.New(pool),
	}
}

// UserCreate creates a new user
func (s *store) UserCreate(ctx context.Context, user domain.UserInterface) error {
	_, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID(),
		Email:        user.Email(),
		Name:         user.Name(),
		PasswordHash: user.PasswordHash(),
		Role:         string(user.Role()),
		CreatedAt:    user.CreatedAt(),
		UpdatedAt:    user.UpdatedAt(),
	})
	return err
}

// UserUpdate - NOT IMPLEMENTED: requires SQL update
func (s *store) UserUpdate(ctx context.Context, user domain.UserInterface) error {
	return errors.New("user update not implemented")
}

// UserDelete - NOT IMPLEMENTED: requires SQL update
func (s *store) UserDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("user delete not implemented")
}

// UserSoftDelete - NOT IMPLEMENTED: requires soft delete column
func (s *store) UserSoftDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("user soft delete not implemented")
}

// UserFindByID finds a user by ID
func (s *store) UserFindByID(ctx context.Context, id uuid.UUID) (domain.UserInterface, error) {
	u, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbUserToDomain(u), nil
}

// UserFindByEmail finds a user by email
func (s *store) UserFindByEmail(ctx context.Context, email string) (domain.UserInterface, error) {
	u, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return dbUserToDomain(u), nil
}

// UserList - NOT IMPLEMENTED: requires SQL update
func (s *store) UserList(ctx context.Context, query domain.QueryInterface) ([]domain.UserInterface, error) {
	return nil, errors.New("user list not implemented")
}

// UserCount - NOT IMPLEMENTED: requires SQL update
func (s *store) UserCount(ctx context.Context, query domain.QueryInterface) (int64, error) {
	return s.queries.CountAdminUsers(ctx)
}

// ContentTypeCreate creates a new content type
func (s *store) ContentTypeCreate(ctx context.Context, ct domain.ContentTypeInterface) error {
	fields, err := json.Marshal(ct.Fields())
	if err != nil {
		return err
	}

	_, err = s.queries.CreateContentType(ctx, db.CreateContentTypeParams{
		ID:        ct.ID(),
		Name:      ct.Name(),
		Slug:      ct.Slug(),
		Fields:    fields,
		CreatedAt: ct.CreatedAt(),
		UpdatedAt: ct.UpdatedAt(),
	})
	return err
}

// ContentTypeUpdate updates an existing content type
func (s *store) ContentTypeUpdate(ctx context.Context, ct domain.ContentTypeInterface) error {
	fields, err := json.Marshal(ct.Fields())
	if err != nil {
		return err
	}

	return s.queries.UpdateContentType(ctx, db.UpdateContentTypeParams{
		ID:        ct.ID(),
		Name:      ct.Name(),
		Fields:    fields,
		UpdatedAt: time.Now().UTC(),
	})
}

// ContentTypeDelete deletes a content type
func (s *store) ContentTypeDelete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteContentType(ctx, id)
}

// ContentTypeFindByID finds a content type by ID
func (s *store) ContentTypeFindByID(ctx context.Context, id uuid.UUID) (domain.ContentTypeInterface, error) {
	ct, err := s.queries.GetContentTypeByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbContentTypeToDomain(ct), nil
}

// ContentTypeFindBySlug finds a content type by slug
func (s *store) ContentTypeFindBySlug(ctx context.Context, slug string) (domain.ContentTypeInterface, error) {
	ct, err := s.queries.GetContentTypeBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	return dbContentTypeToDomain(ct), nil
}

// ContentTypeList lists content types
func (s *store) ContentTypeList(ctx context.Context, query domain.QueryInterface) ([]domain.ContentTypeInterface, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	cts, err := s.queries.ListContentTypes(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.ContentTypeInterface, len(cts))
	for i, ct := range cts {
		result[i] = dbContentTypeToDomain(ct)
	}
	return result, nil
}

// EntryCreate creates a new entry
func (s *store) EntryCreate(ctx context.Context, entry domain.EntryInterface) error {
	fields, err := json.Marshal(entry.Fields())
	if err != nil {
		return err
	}

	var parentID pgtype.UUID
	if entry.ParentID() != nil {
		parentID = pgtype.UUID{Bytes: *entry.ParentID(), Valid: true}
	}

	_, err = s.queries.CreateEntry(ctx, db.CreateEntryParams{
		ID:            entry.ID(),
		ContentTypeID: entry.ContentTypeID(),
		Title:         entry.Title(),
		Slug:          entry.Slug(),
		Path:          entry.Path(),
		ParentID:      parentID,
		AuthorID:      entry.AuthorID(),
		Status:        string(entry.Status()),
		Fields:        fields,
		CreatedAt:     entry.CreatedAt(),
		UpdatedAt:     entry.UpdatedAt(),
	})
	return err
}

// EntryUpdate updates an existing entry
func (s *store) EntryUpdate(ctx context.Context, entry domain.EntryInterface) error {
	fields, err := json.Marshal(entry.Fields())
	if err != nil {
		return err
	}

	return s.queries.UpdateEntry(ctx, db.UpdateEntryParams{
		ID:        entry.ID(),
		Title:     entry.Title(),
		Slug:      entry.Slug(),
		Path:      entry.Path(),
		Fields:    fields,
		UpdatedAt: time.Now().UTC(),
	})
}

// EntryDelete permanently deletes an entry
func (s *store) EntryDelete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteEntry(ctx, id)
}

// EntrySoftDelete - NOT IMPLEMENTED: requires soft delete column
func (s *store) EntrySoftDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("entry soft delete not implemented")
}

// EntryFindByID finds an entry by ID
func (s *store) EntryFindByID(ctx context.Context, id uuid.UUID) (domain.EntryInterface, error) {
	e, err := s.queries.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbEntryToDomain(e), nil
}

// EntryFindBySlug - NOT IMPLEMENTED: requires SQL update
func (s *store) EntryFindBySlug(ctx context.Context, slug string, contentTypeID uuid.UUID) (domain.EntryInterface, error) {
	return nil, errors.New("entry find by slug not implemented")
}

// EntryFindByPath finds an entry by path
func (s *store) EntryFindByPath(ctx context.Context, path string) (domain.EntryInterface, error) {
	e, err := s.queries.GetEntryByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	return dbEntryToDomain(e), nil
}

// EntryList lists entries
func (s *store) EntryList(ctx context.Context, query domain.QueryInterface) ([]domain.EntryInterface, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	entries, err := s.queries.ListEntries(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.EntryInterface, len(entries))
	for i, e := range entries {
		result[i] = dbEntryToDomain(e)
	}
	return result, nil
}

// EntryCount - NOT IMPLEMENTED: requires SQL update
func (s *store) EntryCount(ctx context.Context, query domain.QueryInterface) (int64, error) {
	return 0, errors.New("entry count not implemented")
}

// EntryListTree - NOT IMPLEMENTED: requires SQL update
func (s *store) EntryListTree(ctx context.Context, contentTypeID uuid.UUID) ([]domain.EntryInterface, error) {
	return nil, errors.New("entry list tree not implemented")
}

// EntryVersionCreate - NOT IMPLEMENTED: requires SQL update
func (s *store) EntryVersionCreate(ctx context.Context, version domain.EntryVersionInterface) error {
	return errors.New("entry version create not implemented")
}

// EntryVersionList - NOT IMPLEMENTED: requires SQL update
func (s *store) EntryVersionList(ctx context.Context, entryID uuid.UUID) ([]domain.EntryVersionInterface, error) {
	return nil, errors.New("entry version list not implemented")
}

// MediaFileCreate creates a media file
func (s *store) MediaFileCreate(ctx context.Context, media domain.MediaFileInterface) error {
	_, err := s.queries.CreateMediaFile(ctx, db.CreateMediaFileParams{
		ID:          media.ID(),
		Filename:    media.Filename(),
		MimeType:    media.MimeType(),
		Size:        media.Size(),
		StoragePath: media.StoragePath(),
		UploaderID:  media.UploaderID(),
		CreatedAt:   media.CreatedAt(),
	})
	return err
}

// MediaFileDelete deletes a media file
func (s *store) MediaFileDelete(ctx context.Context, id uuid.UUID) error {
	return s.queries.DeleteMediaFile(ctx, id)
}

// MediaFileFindByID finds a media file by ID
func (s *store) MediaFileFindByID(ctx context.Context, id uuid.UUID) (domain.MediaFileInterface, error) {
	m, err := s.queries.GetMediaFileByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return dbMediaFileToDomain(m), nil
}

// MediaFileList lists media files
func (s *store) MediaFileList(ctx context.Context, query domain.QueryInterface) ([]domain.MediaFileInterface, error) {
	mediaFiles, err := s.queries.ListMediaFiles(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]domain.MediaFileInterface, len(mediaFiles))
	for i, m := range mediaFiles {
		result[i] = dbMediaFileToDomain(m)
	}
	return result, nil
}

// Helper functions to convert db models to domain models

func dbUserToDomain(u db.User) domain.UserInterface {
	user := domain.NewUserFromData(map[string]string{
		"id":            u.ID.String(),
		"email":         u.Email,
		"name":          u.Name,
		"password_hash": u.PasswordHash,
		"role":          u.Role,
		"created_at":    u.CreatedAt.Format(time.RFC3339),
		"updated_at":    u.UpdatedAt.Format(time.RFC3339),
	})
	user.SetSoftDeletedAt(domain.SoftDeleteTime)
	return user
}

func dbContentTypeToDomain(ct db.ContentType) domain.ContentTypeInterface {
	contentType := domain.NewContentTypeFromData(map[string]string{
		"id":         ct.ID.String(),
		"name":       ct.Name,
		"slug":       ct.Slug,
		"created_at": ct.CreatedAt.Format(time.RFC3339),
		"updated_at": ct.UpdatedAt.Format(time.RFC3339),
	})
	contentType.SetSoftDeletedAt(domain.SoftDeleteTime)

	var fields []domain.FieldDefinition
	if err := json.Unmarshal(ct.Fields, &fields); err == nil {
		contentType.SetFields(fields)
	}

	return contentType
}

func dbEntryToDomain(e db.Entry) domain.EntryInterface {
	data := map[string]string{
		"id":              e.ID.String(),
		"content_type_id": e.ContentTypeID.String(),
		"title":           e.Title,
		"slug":            e.Slug,
		"path":            e.Path,
		"author_id":       e.AuthorID.String(),
		"status":          e.Status,
		"created_at":      e.CreatedAt.Format(time.RFC3339),
		"updated_at":      e.UpdatedAt.Format(time.RFC3339),
	}

	if e.ParentID.Valid {
		data["parent_id"] = uuid.UUID(e.ParentID.Bytes).String()
	}

	entry := domain.NewEntryFromData(data)
	entry.SetSoftDeletedAt(domain.SoftDeleteTime)

	if e.PublishedAt.Valid {
		entry.SetPublishedAt(&e.PublishedAt.Time)
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(e.Fields, &fields); err == nil {
		entry.SetFields(fields)
	}

	return entry
}

func dbMediaFileToDomain(m db.MediaFile) domain.MediaFileInterface {
	media := domain.NewMediaFileFromData(map[string]string{
		"id":           m.ID.String(),
		"filename":     m.Filename,
		"mime_type":    m.MimeType,
		"size":         string(rune(m.Size)),
		"storage_path": m.StoragePath,
		"uploader_id":  m.UploaderID.String(),
		"created_at":   m.CreatedAt.Format(time.RFC3339),
	})
	media.SetSoftDeletedAt(domain.SoftDeleteTime)
	return media
}
