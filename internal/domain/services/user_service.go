package services

import (
	"context"
	"errors"
	"walletapitest/internal/domain/entities"
	"walletapitest/internal/domain/repositories"
	
	"github.com/google/uuid"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists     = errors.New("email already exists")
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, email, username, password string) (*entities.User, error) {
	// Проверяем существование пользователя с таким email
	existing, _ := s.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailExists
	}
	
	user := entities.NewUser(email, username, password)
	
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, ErrUserNotFound
	}
	
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	
	if user == nil {
		return nil, ErrUserNotFound
	}
	
	// В реальном приложении сравниваем хеши паролей
	if user.Password != password {
		return nil, ErrInvalidPassword
	}
	
	return user, nil
}