package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vibezz/cms/internal/content"
)

type PublicHandler struct {
	entries      *content.EntryService
	contentTypes *content.ContentTypeService
}

func NewPublicHandler(entries *content.EntryService, contentTypes *content.ContentTypeService) *PublicHandler {
	return &PublicHandler{entries: entries, contentTypes: contentTypes}
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
