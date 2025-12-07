package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// User структура
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// New создает нового пользователя
func New(email, username, password string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Ошибки
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrEmailExists     = errors.New("email already exists")
)

// Repository интерфейс
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// Service структура
type Service struct {
	repo Repository
}

// NewService создает новый экземпляр Service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// NewMockService создает мок-сервис для разработки без БД
func NewMockService() *Service {
	return &Service{
		repo: &mockRepository{},
	}
}

// Create создает нового пользователя
func (s *Service) Create(ctx context.Context, email, username, password string) (*User, error) {
	existing, _ := s.repo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailExists
	}

	user := New(email, username, password)

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID возвращает пользователя по ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// Authenticate аутентифицирует пользователя
func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	if user.Password != password {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// mockRepository для разработки без БД
type mockRepository struct {
	users  map[uuid.UUID]*User
	emails map[string]*User
}

func (m *mockRepository) Create(ctx context.Context, user *User) error {
	if m.users == nil {
		m.users = make(map[uuid.UUID]*User)
		m.emails = make(map[string]*User)
	}
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if m.users == nil {
		return nil, nil
	}
	return m.users[id], nil
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	if m.emails == nil {
		return nil, nil
	}
	return m.emails[email], nil
}

func (m *mockRepository) Update(ctx context.Context, user *User) error {
	if m.users == nil {
		m.users = make(map[uuid.UUID]*User)
		m.emails = make(map[string]*User)
	}
	m.users[user.ID] = user
	m.emails[user.Email] = user
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.users == nil {
		return nil
	}
	if user, exists := m.users[id]; exists {
		delete(m.emails, user.Email)
		delete(m.users, id)
	}
	return nil
}