package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"walletapitest/internal/domain/entities"
	"walletapitest/internal/domain/repositories"
)

type WalletRepositoryImpl struct {
	db *sqlx.DB
}

func NewWalletRepository(db *sqlx.DB) repositories.WalletRepository {
	return &WalletRepositoryImpl{db: db}
}

// ProcessOperationAtomic выполняет атомарную операцию пополнения или списания
func (r *WalletRepositoryImpl) ProcessOperationAtomic(
	ctx context.Context,
	walletID uuid.UUID,
	operationType entities.OperationType,
	amount int64,
) error {
	// Начинаем транзакцию
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Для DEPOSIT - пополнение
	if operationType == entities.OperationTypeDeposit {
		query := `
			UPDATE wallets 
			SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2
			RETURNING balance, user_id
		`
		var newBalance int64
		var userID uuid.UUID
		err := tx.QueryRowContext(ctx, query, amount, walletID).Scan(&newBalance, &userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("wallet not found")
			}
			return err
		}

		// Логируем операцию пополнения
		insertQuery := `
			INSERT INTO operations (id, wallet_id, user_id, operation_type, amount, balance_after, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.ExecContext(ctx, insertQuery,
			uuid.New(),
			walletID,
			userID,
			operationType,
			amount,
			newBalance,
			time.Now(),
		)
		if err != nil {
			return err
		}

	} else if operationType == entities.OperationTypeWithdraw {
		// Для WITHDRAW - списание с проверкой баланса
		query := `
			UPDATE wallets 
			SET balance = balance - $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND balance >= $1
			RETURNING balance, user_id
		`
		var newBalance int64
		var userID uuid.UUID
		err := tx.QueryRowContext(ctx, query, amount, walletID).Scan(&newBalance, &userID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Проверяем, существует ли кошелек
				var exists bool
				checkQuery := `SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1)`
				err = tx.QueryRowContext(ctx, checkQuery, walletID).Scan(&exists)
				if err != nil {
					return err
				}
				if !exists {
					return errors.New("wallet not found")
				}
				// Кошелек существует, но недостаточно средств
				return errors.New("insufficient funds")
			}
			return err
		}

		// Логируем операцию списания
		insertQuery := `
			INSERT INTO operations (id, wallet_id, user_id, operation_type, amount, balance_after, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		_, err = tx.ExecContext(ctx, insertQuery,
			uuid.New(),
			walletID,
			userID,
			operationType,
			amount,
			newBalance,
			time.Now(),
		)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid operation type")
	}

	// Коммитим транзакцию
	err = tx.Commit()
	return err
}

// Также нужно обновить остальные методы репозитория для поддержки транзакций:

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

// GetOperationsHistory - получение истории операций по кошельку
func (r *WalletRepositoryImpl) GetOperationsHistory(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*entities.Operation, error) {
	var operations []*entities.Operation
	query := `
		SELECT * FROM operations 
		WHERE wallet_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &operations, query, walletID, limit, offset)
	if err != nil {
		return nil, err
	}
	
	return operations, nil
}

// GetOperationsHistoryByUser - получение истории операций по пользователю
func (r *WalletRepositoryImpl) GetOperationsHistoryByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entities.Operation, error) {
	var operations []*entities.Operation
	query := `
		SELECT * FROM operations 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`
	
	err := r.db.SelectContext(ctx, &operations, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	
	return operations, nil
}