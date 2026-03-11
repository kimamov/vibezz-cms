package content

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func SeedDefaults(ctx context.Context, pool *pgxpool.Pool) error {
	var count int
	err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM content_types WHERE slug = 'page'`).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Info().Msg("seeding default 'Page' content type")

	fields := []FieldDefinitionInput{
		{Name: "Blocks", Slug: "blocks", Type: "blocks", Required: false},
	}

	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}

	now := time.Now()
	_, err = pool.Exec(ctx,
		`INSERT INTO content_types (id, name, slug, fields, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New(), "Page", "page", fieldsJSON, now, now)
	if err != nil {
		return err
	}

	log.Info().Msg("default 'Page' content type created")
	return nil
}
