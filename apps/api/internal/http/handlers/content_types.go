package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/content"
)

type ContentTypeHandler struct {
	service *content.ContentTypeService
}

func NewContentTypeHandler(service *content.ContentTypeService) *ContentTypeHandler {
	return &ContentTypeHandler{service: service}
}

type createContentTypeRequest struct {
	Name   string                        `json:"name" binding:"required"`
	Slug   string                        `json:"slug" binding:"required"`
	Fields []content.FieldDefinitionInput `json:"fields"`
}

func (h *ContentTypeHandler) Create(c *gin.Context) {
	var req createContentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	ct, err := h.service.Create(c.Request.Context(), req.Name, req.Slug, req.Fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ct)
}

func (h *ContentTypeHandler) List(c *gin.Context) {
	types, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types)
}

func (h *ContentTypeHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ct, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "content type not found"})
		return
	}

	c.JSON(http.StatusOK, ct)
}

type updateContentTypeRequest struct {
	Name   *string                        `json:"name"`
	Fields *[]content.FieldDefinitionInput `json:"fields"`
}

func (h *ContentTypeHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateContentTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	ct, err := h.service.Update(c.Request.Context(), id, req.Name, req.Fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ct)
}

func (h *ContentTypeHandler) Delete(c *gin.Context) {
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
