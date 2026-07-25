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

	admin := &models.User{
		Username:     "admin",
		PasswordHash: "admin123",
		Role:         "admin",
	}

	return repo.Create(admin)
}