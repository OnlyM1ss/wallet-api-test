package repositories

import (
	"context"
	"walletapitest/internal/domain/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet *entities.Wallet) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Wallet, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Wallet, error)
	Update(ctx context.Context, wallet *entities.Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, amount int64) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Transaction methods
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	UpdateBalanceWithTx(ctx context.Context, tx *sqlx.Tx, walletID uuid.UUID, amount int64) error
	FindByIDWithTx(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entities.Wallet, error)

	ProcessOperationAtomic(ctx context.Context, walletID uuid.UUID, operationType entities.OperationType, amount int64) error
}
