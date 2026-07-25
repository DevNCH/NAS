package database

import (
	"github.com/DevNCH/NAS/internal/models"
	"github.com/DevNCH/NAS/internal/repository"
)

func Seed() error {

	repo := repository.NewUserRepository(GetDB())

	exists, err := repo.HasUsers()
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)

	admin := &models.User{
		Username:     "admin",
		PasswordHash: string(hash),
		Role:         "admin",
	}

	return repo.Create(admin)
}
