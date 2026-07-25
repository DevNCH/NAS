package services

import (
	"errors"

	"github.com/DevNCH/NAS/internal/models"
	"github.com/DevNCH/NAS/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
) *AuthService {

	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(
	username string,
	password string,
	role string,
) error {

	existingUser, _ := s.userRepo.GetByUsername(username)

	if existingUser != nil {
		return errors.New("usuário já existe")
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
	}

	return s.userRepo.Create(&user)
}

func (s *AuthService) Login(
	username string,
	password string,
) (*models.User, error) {

	user, err := s.userRepo.GetByUsername(username)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
