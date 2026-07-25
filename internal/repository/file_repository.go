package repository

import (
	"database/sql"

	"github.com/DevNCH/NAS/internal/models"
)

type FileRepository struct {
	DB *sql.DB
}

func NewFileRepository(db *sql.DB) *FileRepository {
	return &FileRepository{
		DB: db,
	}
}

func (r *FileRepository) Create(file *models.File) error {

	query := `
	INSERT INTO files (
		filename,
		filepath,
		size,
		uploaded_by
	)
	VALUES (?, ?, ?, ?)
	`

	result, err := r.DB.Exec(
		query,
		file.Filename,
		file.Filepath,
		file.Size,
		file.UploadedBy,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	file.ID = int(id)

	return nil
}

func (r *FileRepository) GetByID(id int) (*models.File, error) {

	query := `
	SELECT
		id,
		filename,
		filepath,
		size,
		uploaded_by,
		created_at
	FROM files
	WHERE id = ?
	`

	var file models.File

	err := r.DB.QueryRow(query, id).Scan(
		&file.ID,
		&file.Filename,
		&file.Filepath,
		&file.Size,
		&file.UploadedBy,
		&file.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *FileRepository) GetAll() ([]models.File, error) {

	query := `
	SELECT
		id,
		filename,
		filepath,
		size,
		uploaded_by,
		created_at
	FROM files
	`

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []models.File

	for rows.Next() {

		var file models.File

		err := rows.Scan(
			&file.ID,
			&file.Filename,
			&file.Filepath,
			&file.Size,
			&file.UploadedBy,
			&file.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
