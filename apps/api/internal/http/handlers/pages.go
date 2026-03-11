package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/content"
)

type PagesHandler struct {
	pages *content.PageService
	ct    *content.ContentTypeService
}

func NewPagesHandler(pages *content.PageService, ct *content.ContentTypeService) *PagesHandler {
	return &PagesHandler{pages: pages, ct: ct}
}

func (h *PagesHandler) ListTree(c *gin.Context) {
	tree, err := h.pages.ListTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tree)
}

type createPageRequest struct {
	Title    string                 `json:"title" binding:"required"`
	Slug     string                 `json:"slug"`
	ParentID *uuid.UUID             `json:"parent_id"`
	Fields   map[string]interface{} `json:"fields"`
}

func (h *PagesHandler) Create(c *gin.Context) {
	var req createPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	if req.Fields == nil {
		req.Fields = map[string]interface{}{}
	}

	authorID := getUserID(c)

	page, err := h.pages.Create(c.Request.Context(), authorID, req.Title, req.Slug, req.ParentID, req.Fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, page)
}

func (h *PagesHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, err := h.pages.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	ct, _ := h.ct.GetBySlug(c.Request.Context(), "page")

	c.JSON(http.StatusOK, gin.H{
		"page":         page,
		"content_type": ct,
	})
}

type updatePageRequest struct {
	Title  *string                 `json:"title"`
	Slug   *string                 `json:"slug"`
	Fields *map[string]interface{} `json:"fields"`
}

func (h *PagesHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	page, err := h.pages.Update(c.Request.Context(), id, content.UpdateEntryInput{
		Title:  req.Title,
		Slug:   req.Slug,
		Fields: req.Fields,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *PagesHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, err := h.pages.Publish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *PagesHandler) Unpublish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	page, err := h.pages.Unpublish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, page)
}

func (h *PagesHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.pages.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
