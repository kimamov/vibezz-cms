package news

import (
	"time"

	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/content"
)

// NewsItem is the plugin's model for a news entry with typed fields.
type NewsItem struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Path        string     `json:"path"`
	Status      string     `json:"status"`
	Headline    string     `json:"headline"`
	Excerpt     string     `json:"excerpt"`
	Body        string     `json:"body"`
	Featured    bool       `json:"featured"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// FromEntry maps a generic entry to the NewsItem model.
func FromEntry(e *content.Entry) NewsItem {
	ni := NewsItem{
		ID:          e.ID,
		Title:       e.Title,
		Slug:        e.Slug,
		Path:        e.Path,
		Status:      e.Status,
		PublishedAt: e.PublishedAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
	if e.Fields != nil {
		if v, ok := e.Fields["headline"].(string); ok {
			ni.Headline = v
		}
		if v, ok := e.Fields["excerpt"].(string); ok {
			ni.Excerpt = v
		}
		if v, ok := e.Fields["body"].(string); ok {
			ni.Body = v
		}
		if v, ok := e.Fields["featured"].(bool); ok {
			ni.Featured = v
		}
	}
	return ni
}
