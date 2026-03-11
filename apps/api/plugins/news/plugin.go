package news

import (
	"github.com/vibezz/cms/internal/content"
	"github.com/vibezz/cms/internal/http/middleware"
	"github.com/vibezz/cms/internal/plugin"
)

const blockTypeNewsList = "news_list"

// Plugin implements the plugin.Plugin interface for the News content type.
type Plugin struct{}

func New() *Plugin {
	return &Plugin{}
}

func (p *Plugin) Name() string {
	return "News"
}

func (p *Plugin) Slug() string {
	return "news"
}

func (p *Plugin) ContentType() plugin.ContentTypeDefinition {
	return plugin.ContentTypeDefinition{
		Name: "News",
		Fields: []content.FieldDefinitionInput{
			{Name: "Headline", Slug: "headline", Type: "text", Required: true},
			{Name: "Excerpt", Slug: "excerpt", Type: "textarea", Required: false},
			{Name: "Body", Slug: "body", Type: "richtext", Required: true},
			{Name: "Featured", Slug: "featured", Type: "boolean", Required: false},
		},
	}
}

func (p *Plugin) Register(deps *plugin.Deps) {
	h := NewHandler(deps)
	deps.Admin.GET("/news", h.List)
	deps.Admin.POST("/news", h.Create)
	deps.Admin.GET("/news/:id", h.Get)
	deps.Admin.PATCH("/news/:id", h.Update)
	deps.Admin.DELETE("/news/:id", middleware.RoleRequired("admin"), h.Delete)
	deps.Admin.POST("/news/:id/publish", middleware.RoleRequired("admin", "editor"), h.Publish)
	deps.Admin.POST("/news/:id/unpublish", middleware.RoleRequired("admin", "editor"), h.Unpublish)

	deps.Public.GET("/news", h.ListPublic)
	deps.Public.GET("/news/:slug", h.GetBySlugPublic)

	deps.RegisterBlockType(plugin.BlockTypeDefinition{
		Slug:  blockTypeNewsList,
		Label: "News List",
		Icon:  "i-heroicons-newspaper",
		ConfigFields: []plugin.BlockConfigField{
			{Name: "Limit", Slug: "limit", Type: "number", Required: false},
			{Name: "Featured only", Slug: "featured_only", Type: "boolean", Required: false},
		},
		DefaultData: map[string]interface{}{
			"limit":         float64(5),
			"featured_only": false,
		},
	})
	deps.RegisterBlockEnricher(blockTypeNewsList, h.EnrichNewsListBlock)
}
