package repository

import ("database/sql"
	"github.com/DevNCH/NAS/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {

	_, err := r.db.Exec(
		`INSERT INTO users (username, password_hash, role)
		 VALUES (?, ?, ?)`,
		user.Username,
		user.PasswordHash,
		user.Role,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id int64) (*models.User, error) {
	return nil, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	return nil, nil
}	

func (r *UserRepository) Update(user *models.User) error {

    return nil
}

func (r *UserRepository) Delete(id int64) error {

    return nil
}

func (r *UserRepository) HasUsers() (bool, error) {

	var count int

	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}