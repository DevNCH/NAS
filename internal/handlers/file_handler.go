package handlers

import (
	"net/http"
	"strconv"

	"github.com/DevNCH/NAS/internal/middleware"
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

func (h *FileHandler) Upload(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "arquivo não enviado",
		})
		return
	}

	session, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "não autenticado",
		})
		return
	}

	filepath := "storage/uploads/" + file.Filename

	err = c.SaveUploadedFile(file, filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.fileService.Upload(
		file.Filename,
		filepath,
		file.Size,
		session.UserID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "upload realizado",
	})
}

func (h *FileHandler) Download(
	c *gin.Context,
) {

	id, err := strconv.Atoi(
		c.Param("id"),
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id inválido",
		})
		return
	}

	file, err := h.fileService.GetFile(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "arquivo não encontrado",
		})
		return
	}

	c.FileAttachment(
		file.Filepath,
		file.Filename,
	)
}
