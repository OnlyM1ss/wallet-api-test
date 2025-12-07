package entities

import (
	"time"
	
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func NewUser(email, username, password string) *User {
	return &User{
		ID:        uuid.New(),
		Email:     email,
		Username:  username,
		Password:  password, // В реальном приложении хешировать!
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}