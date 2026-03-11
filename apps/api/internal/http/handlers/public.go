package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/content"
	"github.com/vibezz/cms/internal/media"
	"github.com/vibezz/cms/internal/plugin"
)

type PublicHandler struct {
	entries      *content.EntryService
	contentTypes *content.ContentTypeService
	media        *media.Service
	pages        *content.PageService
	apiBaseURL   string
}

func NewPublicHandler(entries *content.EntryService, contentTypes *content.ContentTypeService, mediaSvc *media.Service, pages *content.PageService, apiBaseURL string) *PublicHandler {
	return &PublicHandler{entries: entries, contentTypes: contentTypes, media: mediaSvc, pages: pages, apiBaseURL: apiBaseURL}
}

func (h *PublicHandler) resolveMediaID(c *gin.Context, idStr string) map[string]interface{} {
	mediaID, err := uuid.Parse(idStr)
	if err != nil {
		return nil
	}
	file, err := h.media.GetByID(c.Request.Context(), mediaID)
	if err != nil {
		return nil
	}
	return map[string]interface{}{
		"id":        file.ID.String(),
		"filename":  file.Filename,
		"mime_type": file.MimeType,
		"size":      file.Size,
		"url":       fmt.Sprintf("%s/api/public/media/%s", h.apiBaseURL, file.ID.String()),
	}
}

func (h *PublicHandler) enrichMediaFields(c *gin.Context, entry *content.Entry) {
	ct, err := h.contentTypes.GetByID(c.Request.Context(), entry.ContentTypeID)
	if err != nil {
		return
	}

	fieldTypes := map[string]string{}
	for _, f := range ct.Fields {
		fieldTypes[f.Slug] = f.Type
	}

	for slug, fType := range fieldTypes {
		val, ok := entry.Fields[slug]
		if !ok || val == nil {
			continue
		}

		switch fType {
		case "media":
			if idStr, ok := val.(string); ok && idStr != "" {
				if resolved := h.resolveMediaID(c, idStr); resolved != nil {
					entry.Fields[slug] = resolved
				}
			}
		case "blocks":
			blocks, ok := val.([]interface{})
			if !ok {
				continue
			}
			for i, raw := range blocks {
				block, ok := raw.(map[string]interface{})
				if !ok {
					continue
				}
				blockType, _ := block["type"].(string)
				data, _ := block["data"].(map[string]interface{})
				if data == nil {
					continue
				}
				if blockType == "image" {
					if mediaID, ok := data["media_id"].(string); ok && mediaID != "" {
						if resolved := h.resolveMediaID(c, mediaID); resolved != nil {
							data["media"] = resolved
						}
					}
				}
				if enricher := plugin.GetBlockEnricher(blockType); enricher != nil {
					if enriched, err := enricher(c.Request.Context(), data); err == nil {
						data = enriched
					}
				}
				block["data"] = data
				blocks[i] = block
			}
			entry.Fields[slug] = blocks
		}
	}
}

func (h *PublicHandler) ResolveRoute(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		path = "/"
	}

	entry, err := h.entries.GetByPath(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	h.enrichMediaFields(c, entry)

	c.JSON(http.StatusOK, entry)
}

func (h *PublicHandler) GetNavigation(c *gin.Context) {
	nav, err := h.entries.GetNavigation(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nav)
}

// GetPageByPath resolves only page entries by path (optimized for the common case).
func (h *PublicHandler) GetPageByPath(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		path = "/"
	}

	entry, err := h.pages.GetByPath(c.Request.Context(), path)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	h.enrichMediaFields(c, entry)
	c.JSON(http.StatusOK, entry)
}
