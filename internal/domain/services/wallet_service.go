package services

import (
	"context"
	"errors"
	"time"
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

	wallet, err := s.walletRepo.FindByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet == nil {
		return ErrWalletNotFound
	}

	switch operationType {
	case entities.OperationTypeDeposit:
		wallet.Balance += amount
	case entities.OperationTypeWithdraw:
		if wallet.Balance < amount {
			return ErrInsufficientFunds
		}
		wallet.Balance -= amount
	default:
		return ErrInvalidOperation
	}

	// Update timestamp
	wallet.UpdatedAt = time.Now()
	return s.walletRepo.Update(ctx, wallet)
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
