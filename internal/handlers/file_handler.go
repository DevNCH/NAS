package handlers

import (
	"net/http"

	"github.com/DevNCH/NAS/internal/services"
	"github.com/gin-gonic/gin"
)

type FileHandler struct {
	fileService *services.FileService
}

func NewFileHandler(
	fileService *services.FileService,
) *FileHandler {

	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) ListFiles(c *gin.Context) {

	files, err := h.fileService.ListFiles()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, files)
}
