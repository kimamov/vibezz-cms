package plugin

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vibezz/cms/internal/config"
	"github.com/vibezz/cms/internal/content"
	"github.com/vibezz/cms/internal/media"
)

// ContentTypeDefinition describes the content type a plugin provides.
// Slug is taken from Plugin.Slug().
type ContentTypeDefinition struct {
	Name   string
	Fields []content.FieldDefinitionInput
}

// BlockConfigField describes a config field for a block type (e.g. "limit" for news_list).
type BlockConfigField struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Type     string `json:"type"` // text, number, boolean
	Required bool   `json:"required"`
}

// Plugin is the interface that custom plugins implement.
type Plugin interface {
	// Name returns the human-readable plugin name (e.g. "News").
	Name() string
	// Slug returns the unique slug used for the content type and URL prefix (e.g. "news").
	Slug() string
	// ContentType returns the content type definition (name + fields).
	ContentType() ContentTypeDefinition
	// Register mounts the plugin's admin and public routes and runs any init logic.
	Register(deps *Deps)
}

// Deps provides plugins with access to core services and their resolved content type.
type Deps struct {
	Pool               *pgxpool.Pool
	Config             *config.Config
	EntryService       *content.EntryService
	ContentTypeService *content.ContentTypeService
	MediaService       *media.Service
	ApiBaseURL         string
	ContentTypeID      uuid.UUID
	Admin              *gin.RouterGroup
	Public             *gin.RouterGroup
}

// RegisterBlockType registers a block type that can be added to pages (e.g. "news_list").
func (d *Deps) RegisterBlockType(def BlockTypeDefinition) {
	registerBlockType(def)
}

// RegisterBlockEnricher registers an enricher for a block type for public API responses.
func (d *Deps) RegisterBlockEnricher(slug string, enricher BlockEnricher) {
	registerBlockEnricher(slug, enricher)
}
