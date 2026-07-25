package repository

import (
	"database/sql"

	"github.com/DevNCH/NAS/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
        INSERT INTO users (
            username,
            password_hash,
            role
        ) VALUES (?, ?, ?)
    `

	result, err := r.DB.Exec(
		query,
		user.Username,
		user.PasswordHash,
		user.Role,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)

	return nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
        SELECT
            id,
            username,
            password_hash,
            role,
            created_at
        FROM users
        WHERE username = ?
    `

	var user models.User

	err := r.DB.QueryRow(
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {

	query := `
	SELECT
		id,
		username,
		password_hash,
		role,
		created_at
	FROM users
	WHERE id = ?
	`

	var user models.User

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) HasUsers() (bool, error) {

	var count int

	err := r.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
