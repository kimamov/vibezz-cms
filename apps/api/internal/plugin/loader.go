package plugin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/vibezz/cms/internal/config"
	"github.com/vibezz/cms/internal/content"
	"github.com/vibezz/cms/internal/media"
)

// Load ensures each plugin's content type exists, then registers the plugin's routes.
func Load(ctx context.Context, plugins []Plugin, pool *pgxpool.Pool, cfg *config.Config,
	entryService *content.EntryService, ctService *content.ContentTypeService, mediaService *media.Service,
	apiBaseURL string, adminGroup, publicGroup *gin.RouterGroup) {

	for _, p := range plugins {
		ctID, err := ensureContentType(ctx, pool, ctService, p.Slug(), p.ContentType())
		if err != nil {
			log.Error().Err(err).Str("plugin", p.Name()).Msg("failed to ensure plugin content type")
			continue
		}

		deps := &Deps{
			Pool:               pool,
			Config:             cfg,
			EntryService:       entryService,
			ContentTypeService: ctService,
			MediaService:       mediaService,
			ApiBaseURL:         apiBaseURL,
			ContentTypeID:      ctID,
			Admin:              adminGroup,
			Public:             publicGroup,
		}

		p.Register(deps)
		log.Info().Str("plugin", p.Name()).Msg("plugin loaded")
	}
}

func ensureContentType(ctx context.Context, pool *pgxpool.Pool, ctService *content.ContentTypeService, slug string, def ContentTypeDefinition) (uuid.UUID, error) {
	existing, err := ctService.GetBySlug(ctx, slug)
	if err == nil {
		return existing.ID, nil
	}

	fieldsJSON, err := json.Marshal(def.Fields)
	if err != nil {
		return uuid.Nil, err
	}

	id := uuid.New()
	now := time.Now()
	_, err = pool.Exec(ctx,
		`INSERT INTO content_types (id, name, slug, fields, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, def.Name, slug, fieldsJSON, now, now)
	if err != nil {
		return uuid.Nil, err
	}

	log.Info().Str("slug", slug).Msg("created content type for plugin")
	return id, nil
}
