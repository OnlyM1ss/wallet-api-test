package services

import (
	"context"
	"errors"
	"walletapitest/internal/domain/entities"
	"walletapitest/internal/domain/repositories"

	"github.com/google/uuid"
)

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidOperation  = errors.New("invalid operation type")
	ErrInvalidAmount     = errors.New("invalid amount")
)

type WalletService struct {
	walletRepo repositories.WalletRepository
}

func NewWalletService(walletRepo repositories.WalletRepository) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
	}
}

func (s *WalletService) ProcessOperation(ctx context.Context, walletID uuid.UUID, operationType entities.OperationType, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	// Начинаем транзакцию
	tx, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return err
	}

	// Откатываем транзакцию в случае ошибки
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Логируем ошибку отката, но возвращаем исходную ошибку
			}
		}
	}()

	// Получаем кошелек с блокировкой строки (FOR UPDATE) для предотвращения race conditions
	wallet, err := s.walletRepo.FindByIDWithTx(ctx, tx, walletID)
	if err != nil {
		return err
	}

	if wallet == nil {
		return ErrWalletNotFound
	}

	// Определяем изменение баланса в зависимости от типа операции
	var balanceChange int64
	switch operationType {
	case entities.OperationTypeDeposit:
		balanceChange = amount
	case entities.OperationTypeWithdraw:
		// Проверяем достаточность средств
		if wallet.Balance < amount {
			return ErrInsufficientFunds
		}
		balanceChange = -amount // Отрицательное значение для снятия
	default:
		return ErrInvalidOperation
	}

	// Обновляем баланс в транзакции
	err = s.walletRepo.UpdateBalanceWithTx(ctx, tx, walletID, balanceChange)
	if err != nil {
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *WalletService) GetWallet(ctx context.Context, walletID uuid.UUID) (*entities.Wallet, error) {
	wallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, ErrWalletNotFound
	}

	return wallet, nil
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID) (*entities.Wallet, error) {
	wallet := entities.NewWallet(userID)

	if err := s.walletRepo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]*entities.Wallet, error) {
	return s.walletRepo.FindByUserID(ctx, userID)
}
