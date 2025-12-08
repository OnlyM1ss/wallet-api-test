package postgres

import (
	"context"
	"database/sql"
	"walletapitest/internal/domain/entities"
	"walletapitest/internal/domain/repositories"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletRepositoryImpl struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) repositories.WalletRepository {
	return &WalletRepositoryImpl{db: db}
}

func (r *WalletRepositoryImpl) Create(ctx context.Context, wallet *entities.Wallet) error {
	query := `
		INSERT INTO wallets (id, user_id, balance, created_at, updated_at)
		VALUES (:id, :user_id, :balance, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, wallet)
	return err
}

func (r *WalletRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error) {
	var wallet entities.Wallet
	query := `SELECT * FROM wallets WHERE id = $1`

	err := r.db.GetContext(ctx, &wallet, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepositoryImpl) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Wallet, error) {
	var wallets []*entities.Wallet
	query := `SELECT * FROM wallets WHERE user_id = $1 ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &wallets, query, userID)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *WalletRepositoryImpl) Update(ctx context.Context, wallet *entities.Wallet) error {
	query := `
		UPDATE wallets 
		SET balance = :balance, updated_at = :updated_at
		WHERE id = :id
	`

	_, err := r.db.NamedExecContext(ctx, query, wallet)
	return err
}

func (r *WalletRepositoryImpl) UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error {
	query := `UPDATE wallets SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, amount, walletID)
	return err
}

func (r *WalletRepositoryImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r *WalletRepositoryImpl) UpdateBalanceWithTx(ctx context.Context, tx *sqlx.Tx, walletID uuid.UUID, amount int64) error {
	query := `UPDATE wallets SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, amount, walletID)
	return err
}

func (r *WalletRepositoryImpl) FindByIDWithTx(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entities.Wallet, error) {
	var wallet entities.Wallet
	query := `SELECT * FROM wallets WHERE id = $1 FOR UPDATE`

	err := tx.GetContext(ctx, &wallet, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM wallets WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
