package postgres

import (
	"context"
	"database/sql"
	"walletapitest/internal/domain/entities"
	"walletapitest/internal/domain/repositories"
	
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) repositories.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, username, password, created_at, updated_at)
		VALUES (:id, :email, :username, :password, :created_at, :updated_at)
	`
	
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE id = $1`
	
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT * FROM users WHERE email = $1`
	
	err := r.db.GetContext(ctx, &user, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE users 
		SET email = :email, username = :username, password = :password, updated_at = :updated_at
		WHERE id = :id
	`
	
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	var users []*entities.User
	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	return users, err
}