package services

import (
	"github.com/DevNCH/NAS/internal/models"
	"github.com/DevNCH/NAS/internal/repository"
)

type FileService struct {
	fileRepo *repository.FileRepository
}

func NewFileService(
	fileRepo *repository.FileRepository,
) *FileService {

	return &FileService{
		fileRepo: fileRepo,
	}
}

func (s *FileService) ListFiles() ([]models.File, error) {
	return s.fileRepo.GetAll()
}

func (s *FileService) GetFile(
	id int,
) (*models.File, error) {

	return s.fileRepo.GetByID(id)
}

func (s *FileService) Upload(
	filename string,
	filepath string,
	size int64,
	userID int,
) error {

	file := models.File{
		Filename:   filename,
		Filepath:   filepath,
		Size:       size,
		UploadedBy: userID,
	}

	return s.fileRepo.Create(&file)
}
