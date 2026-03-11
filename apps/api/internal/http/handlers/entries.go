package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/content"
)

type EntryHandler struct {
	service *content.EntryService
}

func NewEntryHandler(service *content.EntryService) *EntryHandler {
	return &EntryHandler{service: service}
}

type createEntryRequest struct {
	ContentTypeID uuid.UUID              `json:"content_type_id" binding:"required"`
	Title         string                 `json:"title" binding:"required"`
	Slug          string                 `json:"slug"`
	ParentID      *uuid.UUID             `json:"parent_id"`
	Fields        map[string]interface{} `json:"fields"`
}

func (h *EntryHandler) Create(c *gin.Context) {
	var req createEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	authorID := getUserID(c)

	entry, err := h.service.Create(c.Request.Context(), content.CreateEntryInput{
		ContentTypeID: req.ContentTypeID,
		Title:         req.Title,
		Slug:          req.Slug,
		ParentID:      req.ParentID,
		Fields:        req.Fields,
		AuthorID:      authorID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *EntryHandler) List(c *gin.Context) {
	contentTypeID := c.Query("content_type_id")
	parentID := c.Query("parent_id")
	status := c.Query("status")

	filters := content.EntryFilters{Status: status}
	if contentTypeID != "" {
		if id, err := uuid.Parse(contentTypeID); err == nil {
			filters.ContentTypeID = &id
		}
	}
	if parentID != "" {
		if id, err := uuid.Parse(parentID); err == nil {
			filters.ParentID = &id
		}
	}

	entries, err := h.service.List(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func (h *EntryHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

type updateEntryRequest struct {
	Title  *string                 `json:"title"`
	Slug   *string                 `json:"slug"`
	Fields *map[string]interface{} `json:"fields"`
}

func (h *EntryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	entry, err := h.service.Update(c.Request.Context(), id, content.UpdateEntryInput{
		Title:  req.Title,
		Slug:   req.Slug,
		Fields: req.Fields,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.service.Publish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) Unpublish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.service.Unpublish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
