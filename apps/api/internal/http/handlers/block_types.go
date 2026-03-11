package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vibezz/cms/internal/plugin"
)

type BlockTypesHandler struct{}

func NewBlockTypesHandler() *BlockTypesHandler {
	return &BlockTypesHandler{}
}

func (h *BlockTypesHandler) List(c *gin.Context) {
	c.JSON(http.StatusOK, plugin.GetAllBlockTypes())
}
