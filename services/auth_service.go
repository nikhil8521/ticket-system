package services

import (
	"errors"
	"ticket-system/models"
	"ticket-system/repository"
	"ticket-system/utils"
)

type AuthService interface {
	Register(name, email, password string) error
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(name, email, password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateJWT(user.ID)
}
