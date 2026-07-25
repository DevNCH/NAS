package repository

import (
	"database/sql"
	"github.com/DevNCH/NAS/internal/models"
)

type FileRepository struct {
	db *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{db: db}
}

func (r *FileRepository) Create(file *models.File) error {

    return nil
}

func (r *FileRepository) FindByID(id int64) (*models.File, error) {
	return nil, nil
}

func (r *FileRepository) FindByFilename(filename string) (*models.File, error) {
	return nil, nil
}

func (r *FileRepository) Update(file *models.File) error {
	return nil
}

func (r *FileRepository) Delete(id int64) error {
	return nil
}

func (r *FileRepository) FindByUser(userID int64) ([]models.File, error) {
	return nil, nil
}