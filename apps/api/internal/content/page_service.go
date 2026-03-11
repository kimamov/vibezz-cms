package content

import (
	"context"
	"errors"
	"sort"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pageContentTypeSlug = "page"

var errNotAPage = errors.New("entry is not a page")

type PageService struct {
	pool       *pgxpool.Pool
	entries    *EntryService
	contentTypes *ContentTypeService
}

func NewPageService(pool *pgxpool.Pool, entries *EntryService, contentTypes *ContentTypeService) *PageService {
	return &PageService{pool: pool, entries: entries, contentTypes: contentTypes}
}

func (s *PageService) pageTypeID(ctx context.Context) (uuid.UUID, error) {
	ct, err := s.contentTypes.GetBySlug(ctx, pageContentTypeSlug)
	if err != nil {
		return uuid.Nil, err
	}
	return ct.ID, nil
}

func (s *PageService) ListTree(ctx context.Context) ([]PageNode, error) {
	pageID, err := s.pageTypeID(ctx)
	if err != nil {
		return nil, err
	}

	list, err := s.entries.List(ctx, EntryFilters{ContentTypeID: &pageID})
	if err != nil {
		return nil, err
	}

	sort.Slice(list, func(i, j int) bool { return list[i].Path < list[j].Path })

	byParent := make(map[uuid.UUID][]Entry)
	for _, e := range list {
		var key uuid.UUID
		if e.ParentID != nil {
			key = *e.ParentID
		}
		byParent[key] = append(byParent[key], e)
	}

	var build func(parentID uuid.UUID) []PageNode
	build = func(parentID uuid.UUID) []PageNode {
		entries := byParent[parentID]
		if len(entries) == 0 {
			return nil
		}
		nodes := make([]PageNode, 0, len(entries))
		for _, e := range entries {
			nodes = append(nodes, PageNode{
				Entry:    e,
				Children: build(e.ID),
			})
		}
		return nodes
	}

	return build(uuid.Nil), nil
}

func (s *PageService) Create(ctx context.Context, authorID uuid.UUID, title, slug string, parentID *uuid.UUID, fields map[string]interface{}) (*Entry, error) {
	pageID, err := s.pageTypeID(ctx)
	if err != nil {
		return nil, err
	}
	return s.entries.Create(ctx, CreateEntryInput{
		ContentTypeID: pageID,
		Title:         title,
		Slug:          slug,
		ParentID:      parentID,
		Fields:        fields,
		AuthorID:      authorID,
	})
}

func (s *PageService) Get(ctx context.Context, id uuid.UUID) (*Entry, error) {
	pageID, err := s.pageTypeID(ctx)
	if err != nil {
		return nil, err
	}
	entry, err := s.entries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry.ContentTypeID != pageID {
		return nil, errNotAPage
	}
	return entry, nil
}

func (s *PageService) Update(ctx context.Context, id uuid.UUID, input UpdateEntryInput) (*Entry, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	return s.entries.Update(ctx, id, input)
}

func (s *PageService) Publish(ctx context.Context, id uuid.UUID) (*Entry, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	return s.entries.Publish(ctx, id)
}

func (s *PageService) Unpublish(ctx context.Context, id uuid.UUID) (*Entry, error) {
	if _, err := s.Get(ctx, id); err != nil {
		return nil, err
	}
	return s.entries.Unpublish(ctx, id)
}

func (s *PageService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return s.entries.Delete(ctx, id)
}

func (s *PageService) GetByPath(ctx context.Context, path string) (*Entry, error) {
	pageID, err := s.pageTypeID(ctx)
	if err != nil {
		return nil, err
	}
	return s.entries.GetByPathAndContentType(ctx, path, pageID)
}
