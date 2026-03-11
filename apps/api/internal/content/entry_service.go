package content

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EntryService struct {
	pool *pgxpool.Pool
}

func NewEntryService(pool *pgxpool.Pool) *EntryService {
	return &EntryService{pool: pool}
}

func sanitizeSlug(slug string) string {
	slug = strings.Trim(slug, "/")
	slug = strings.TrimSpace(slug)
	return slug
}

func (s *EntryService) buildPath(ctx context.Context, parentID *uuid.UUID, slug string) string {
	slug = sanitizeSlug(slug)

	if parentID == nil {
		if slug == "" {
			return "/"
		}
		return "/" + slug
	}

	var parentPath string
	err := s.pool.QueryRow(ctx, `SELECT path FROM entries WHERE id = $1`, *parentID).Scan(&parentPath)
	if err != nil {
		if slug == "" {
			return "/"
		}
		return "/" + slug
	}

	parentPath = strings.TrimRight(parentPath, "/")
	if slug == "" {
		return parentPath + "/"
	}
	return parentPath + "/" + slug
}

func (s *EntryService) Create(ctx context.Context, input CreateEntryInput) (*Entry, error) {
	id := uuid.New()
	now := time.Now()
	input.Slug = sanitizeSlug(input.Slug)
	var path string
	if input.PathPrefix != nil {
		prefix := strings.TrimRight(*input.PathPrefix, "/")
		if input.Slug == "" {
			path = prefix
		} else {
			path = prefix + "/" + input.Slug
		}
	} else {
		path = s.buildPath(ctx, input.ParentID, input.Slug)
	}

	fieldsJSON, err := json.Marshal(input.Fields)
	if err != nil {
		return nil, err
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO entries (id, content_type_id, title, slug, path, parent_id, author_id, status, fields, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, 'draft', $8, $9, $10)`,
		id, input.ContentTypeID, input.Title, input.Slug, path,
		input.ParentID, input.AuthorID, fieldsJSON, now, now)
	if err != nil {
		return nil, err
	}

	return &Entry{
		ID: id, ContentTypeID: input.ContentTypeID,
		Title: input.Title, Slug: input.Slug, Path: path,
		ParentID: input.ParentID, AuthorID: input.AuthorID,
		Status: "draft", Fields: input.Fields,
		CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (s *EntryService) scanEntry(scan func(dest ...interface{}) error) (*Entry, error) {
	var e Entry
	var fieldsJSON []byte
	err := scan(&e.ID, &e.ContentTypeID, &e.Title, &e.Slug, &e.Path,
		&e.ParentID, &e.AuthorID, &e.Status, &fieldsJSON,
		&e.PublishedAt, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fieldsJSON, &e.Fields); err != nil {
		e.Fields = map[string]interface{}{}
	}
	return &e, nil
}

const entryColumns = `id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at`

func (s *EntryService) GetByID(ctx context.Context, id uuid.UUID) (*Entry, error) {
	row := s.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM entries WHERE id = $1`, entryColumns), id)
	return s.scanEntry(row.Scan)
}

func (s *EntryService) GetByPath(ctx context.Context, path string) (*Entry, error) {
	row := s.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM entries WHERE path = $1 AND status = 'published'`, entryColumns), path)
	return s.scanEntry(row.Scan)
}

func (s *EntryService) GetByPathAndContentType(ctx context.Context, path string, contentTypeID uuid.UUID) (*Entry, error) {
	row := s.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM entries WHERE path = $1 AND content_type_id = $2 AND status = 'published'`, entryColumns), path, contentTypeID)
	return s.scanEntry(row.Scan)
}

func (s *EntryService) GetByContentTypeAndSlug(ctx context.Context, contentTypeID uuid.UUID, slug string) (*Entry, error) {
	row := s.pool.QueryRow(ctx,
		fmt.Sprintf(`SELECT %s FROM entries WHERE content_type_id = $1 AND slug = $2 AND status = 'published'`, entryColumns), contentTypeID, slug)
	return s.scanEntry(row.Scan)
}

func (s *EntryService) List(ctx context.Context, filters EntryFilters) ([]Entry, error) {
	query := fmt.Sprintf(`SELECT %s FROM entries WHERE 1=1`, entryColumns)
	args := []interface{}{}
	argIdx := 1

	if filters.ContentTypeID != nil {
		query += fmt.Sprintf(` AND content_type_id = $%d`, argIdx)
		args = append(args, *filters.ContentTypeID)
		argIdx++
	}
	if filters.ParentID != nil {
		query += fmt.Sprintf(` AND parent_id = $%d`, argIdx)
		args = append(args, *filters.ParentID)
		argIdx++
	}
	if filters.Status != "" {
		query += fmt.Sprintf(` AND status = $%d`, argIdx)
		args = append(args, filters.Status)
		argIdx++
	}
	query += ` ORDER BY created_at DESC`

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		e, err := s.scanEntry(rows.Scan)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *e)
	}
	if entries == nil {
		entries = []Entry{}
	}
	return entries, nil
}

func (s *EntryService) Update(ctx context.Context, id uuid.UUID, input UpdateEntryInput) (*Entry, error) {
	entry, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		entry.Title = *input.Title
	}
	if input.Slug != nil {
		entry.Slug = sanitizeSlug(*input.Slug)
		entry.Path = s.buildPath(ctx, entry.ParentID, entry.Slug)
	}
	if input.Fields != nil {
		entry.Fields = *input.Fields
	}

	fieldsJSON, err := json.Marshal(entry.Fields)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	_, err = s.pool.Exec(ctx,
		`UPDATE entries SET title=$1, slug=$2, path=$3, fields=$4, updated_at=$5 WHERE id=$6`,
		entry.Title, entry.Slug, entry.Path, fieldsJSON, now, id)
	if err != nil {
		return nil, err
	}

	entry.UpdatedAt = now
	return entry, nil
}

func (s *EntryService) Publish(ctx context.Context, id uuid.UUID) (*Entry, error) {
	now := time.Now()
	_, err := s.pool.Exec(ctx,
		`UPDATE entries SET status='published', published_at=$1, updated_at=$2 WHERE id=$3`,
		now, now, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *EntryService) Unpublish(ctx context.Context, id uuid.UUID) (*Entry, error) {
	now := time.Now()
	_, err := s.pool.Exec(ctx,
		`UPDATE entries SET status='draft', published_at=NULL, updated_at=$1 WHERE id=$2`,
		now, id)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *EntryService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM entries WHERE id = $1`, id)
	return err
}

func (s *EntryService) GetNavigation(ctx context.Context) ([]NavigationItem, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, title, slug, path, parent_id FROM entries
		 WHERE status = 'published' ORDER BY path`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type flatItem struct {
		ID       uuid.UUID
		Title    string
		Slug     string
		Path     string
		ParentID *uuid.UUID
	}

	var items []flatItem
	for rows.Next() {
		var item flatItem
		if err := rows.Scan(&item.ID, &item.Title, &item.Slug, &item.Path, &item.ParentID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	childMap := make(map[uuid.UUID][]NavigationItem)
	var roots []NavigationItem

	for _, item := range items {
		nav := NavigationItem{
			ID: item.ID, Title: item.Title,
			Slug: item.Slug, Path: item.Path,
		}
		if item.ParentID != nil {
			childMap[*item.ParentID] = append(childMap[*item.ParentID], nav)
		} else {
			roots = append(roots, nav)
		}
	}

	var attachChildren func(items []NavigationItem) []NavigationItem
	attachChildren = func(items []NavigationItem) []NavigationItem {
		for i := range items {
			if children, ok := childMap[items[i].ID]; ok {
				items[i].Children = attachChildren(children)
			}
		}
		return items
	}

	if roots == nil {
		roots = []NavigationItem{}
	}
	return attachChildren(roots), nil
}
