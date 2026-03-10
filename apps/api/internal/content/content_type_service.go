package content

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContentTypeService struct {
	pool *pgxpool.Pool
}

func NewContentTypeService(pool *pgxpool.Pool) *ContentTypeService {
	return &ContentTypeService{pool: pool}
}

func (s *ContentTypeService) Create(ctx context.Context, name, slug string, fields []FieldDefinitionInput) (*ContentType, error) {
	id := uuid.New()
	now := time.Now()

	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	_, err = s.pool.Exec(ctx,
		`INSERT INTO content_types (id, name, slug, fields, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, name, slug, fieldsJSON, now, now)
	if err != nil {
		return nil, err
	}

	defs := make([]FieldDefinition, len(fields))
	for i, f := range fields {
		defs[i] = FieldDefinition{
			Name: f.Name, Slug: f.Slug, Type: f.Type,
			Required: f.Required, Options: f.Options, SortOrder: i,
		}
	}

	return &ContentType{
		ID: id, Name: name, Slug: slug,
		Fields: defs, CreatedAt: now, UpdatedAt: now,
	}, nil
}

func (s *ContentTypeService) List(ctx context.Context) ([]ContentType, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, name, slug, fields, created_at, updated_at
		 FROM content_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ContentType
	for rows.Next() {
		var ct ContentType
		var fieldsJSON []byte
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Slug, &fieldsJSON, &ct.CreatedAt, &ct.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(fieldsJSON, &ct.Fields); err != nil {
			ct.Fields = []FieldDefinition{}
		}
		types = append(types, ct)
	}
	if types == nil {
		types = []ContentType{}
	}
	return types, nil
}

func (s *ContentTypeService) GetByID(ctx context.Context, id uuid.UUID) (*ContentType, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, name, slug, fields, created_at, updated_at
		 FROM content_types WHERE id = $1`, id)

	var ct ContentType
	var fieldsJSON []byte
	if err := row.Scan(&ct.ID, &ct.Name, &ct.Slug, &fieldsJSON, &ct.CreatedAt, &ct.UpdatedAt); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(fieldsJSON, &ct.Fields); err != nil {
		ct.Fields = []FieldDefinition{}
	}
	return &ct, nil
}

func (s *ContentTypeService) Update(ctx context.Context, id uuid.UUID, name *string, fields *[]FieldDefinitionInput) (*ContentType, error) {
	ct, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		ct.Name = *name
	}

	if fields != nil {
		fieldsJSON, err := json.Marshal(*fields)
		if err != nil {
			return nil, err
		}
		_, err = s.pool.Exec(ctx,
			`UPDATE content_types SET name = $1, fields = $2, updated_at = $3 WHERE id = $4`,
			ct.Name, fieldsJSON, time.Now(), id)
		if err != nil {
			return nil, err
		}
		defs := make([]FieldDefinition, len(*fields))
		for i, f := range *fields {
			defs[i] = FieldDefinition{
				Name: f.Name, Slug: f.Slug, Type: f.Type,
				Required: f.Required, Options: f.Options, SortOrder: i,
			}
		}
		ct.Fields = defs
	} else {
		_, err = s.pool.Exec(ctx,
			`UPDATE content_types SET name = $1, updated_at = $2 WHERE id = $3`,
			ct.Name, time.Now(), id)
		if err != nil {
			return nil, err
		}
	}

	return ct, nil
}

func (s *ContentTypeService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM content_types WHERE id = $1`, id)
	return err
}
