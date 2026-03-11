package news

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/content"
	"github.com/vibezz/cms/internal/plugin"
)

const pathPrefix = "/news"

type Handler struct {
	deps *plugin.Deps
}

func NewHandler(deps *plugin.Deps) *Handler {
	return &Handler{deps: deps}
}

func (h *Handler) getAuthorID(c *gin.Context) (uuid.UUID, bool) {
	id, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	uid, ok := id.(uuid.UUID)
	return uid, ok
}

func (h *Handler) List(c *gin.Context) {
	entries, err := h.deps.EntryService.List(c.Request.Context(), content.EntryFilters{
		ContentTypeID: &h.deps.ContentTypeID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

type createNewsRequest struct {
	Title   string                 `json:"title" binding:"required"`
	Slug    string                 `json:"slug" binding:"required"`
	Fields  map[string]interface{} `json:"fields"`
}

func (h *Handler) Create(c *gin.Context) {
	authorID, ok := h.getAuthorID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	if req.Fields == nil {
		req.Fields = map[string]interface{}{}
	}

	pathPrefix := pathPrefix
	entry, err := h.deps.EntryService.Create(c.Request.Context(), content.CreateEntryInput{
		ContentTypeID: h.deps.ContentTypeID,
		Title:         req.Title,
		Slug:          strings.TrimSpace(strings.Trim(req.Slug, "/")),
		PathPrefix:    &pathPrefix,
		Fields:        req.Fields,
		AuthorID:      authorID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entry)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.deps.EntryService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if entry.ContentTypeID != h.deps.ContentTypeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

type updateNewsRequest struct {
	Title  *string                 `json:"title"`
	Slug   *string                 `json:"slug"`
	Fields *map[string]interface{} `json:"fields"`
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.deps.EntryService.GetByID(c.Request.Context(), id)
	if err != nil || entry.ContentTypeID != h.deps.ContentTypeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	var req updateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	updated, err := h.deps.EntryService.Update(c.Request.Context(), id, content.UpdateEntryInput{
		Title:  req.Title,
		Slug:   req.Slug,
		Fields: req.Fields,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.deps.EntryService.GetByID(c.Request.Context(), id)
	if err != nil || entry.ContentTypeID != h.deps.ContentTypeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	if err := h.deps.EntryService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.deps.EntryService.GetByID(c.Request.Context(), id)
	if err != nil || entry.ContentTypeID != h.deps.ContentTypeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	entry, err = h.deps.EntryService.Publish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *Handler) Unpublish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.deps.EntryService.GetByID(c.Request.Context(), id)
	if err != nil || entry.ContentTypeID != h.deps.ContentTypeID {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	entry, err = h.deps.EntryService.Unpublish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *Handler) ListPublic(c *gin.Context) {
	entries, err := h.deps.EntryService.List(c.Request.Context(), content.EntryFilters{
		ContentTypeID: &h.deps.ContentTypeID,
		Status:        "published",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *Handler) GetBySlugPublic(c *gin.Context) {
	slug := c.Param("slug")
	entry, err := h.deps.EntryService.GetByContentTypeAndSlug(c.Request.Context(), h.deps.ContentTypeID, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

// EnrichNewsListBlock implements plugin.BlockEnricher: fetches news entries and adds "items" to block data.
func (h *Handler) EnrichNewsListBlock(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	limit := 5
	if v, ok := data["limit"]; ok {
		switch n := v.(type) {
		case float64:
			limit = int(n)
		case int:
			limit = n
		}
	}
	if limit <= 0 || limit > 50 {
		limit = 5
	}

	featuredOnly := false
	if v, ok := data["featured_only"]; ok {
		if b, ok := v.(bool); ok {
			featuredOnly = b
		}
	}

	entries, err := h.deps.EntryService.List(ctx, content.EntryFilters{
		ContentTypeID: &h.deps.ContentTypeID,
		Status:        "published",
	})
	if err != nil {
		return data, err
	}

	if featuredOnly {
		filtered := entries[:0]
		for _, e := range entries {
			if v, ok := e.Fields["featured"].(bool); ok && v {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if len(entries) > limit {
		entries = entries[:limit]
	}

	items := make([]map[string]interface{}, 0, len(entries))
	for _, e := range entries {
		items = append(items, map[string]interface{}{
			"id":           e.ID.String(),
			"title":        e.Title,
			"slug":         e.Slug,
			"path":         e.Path,
			"headline":     e.Fields["headline"],
			"excerpt":      e.Fields["excerpt"],
			"published_at": e.PublishedAt,
		})
	}

	out := make(map[string]interface{}, len(data)+1)
	for k, v := range data {
		out[k] = v
	}
	out["items"] = items
	return out, nil
}
