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

func (s *WalletService) ProcessOperation(
	ctx context.Context,
	walletID uuid.UUID,
	operationType entities.OperationType,
	amount int64,
) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := s.walletRepo.ProcessOperationAtomic(ctx, walletID, operationType, amount)
	
	if err != nil {
		switch err {
		case ErrWalletNotFound:
			return ErrWalletNotFound
		case ErrInsufficientFunds:
			return ErrInsufficientFunds
		case ErrInvalidOperation:
			return ErrInvalidOperation
		}
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
