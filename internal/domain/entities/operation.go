package entities

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationTypeDeposit  OperationType = "DEPOSIT"
	OperationTypeWithdraw OperationType = "WITHDRAW"
)

type Operation struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	WalletID      uuid.UUID     `json:"wallet_id" db:"wallet_id"`
	OperationType OperationType `json:"operation_type" db:"operation_type"`
	Amount        int64         `json:"amount" db:"amount"`
	BalanceAfter  int64         `json:"balance_after" db:"balance_after"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
}

func NewOperation(walletID uuid.UUID, operationType OperationType, amount, balanceAfter int64) *Operation {
	return &Operation{
		ID:            uuid.New(),
		WalletID:      walletID,
		OperationType: operationType,
		Amount:        amount,
		BalanceAfter:  balanceAfter,
		CreatedAt:     time.Now(),
	}
}