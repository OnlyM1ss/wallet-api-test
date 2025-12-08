package handlers

import (
	"net/http"

	"walletapitest/internal/domain/entities"
	"walletapitest/internal/domain/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

type WalletOperationRequest struct {
	WalletID      uuid.UUID `json:"walletId" binding:"required"`
	OperationType string    `json:"operationType" binding:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        int64       `json:"amount" binding:"required,gt=0"`
}

type WalletResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   int64       `json:"balance"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type CreateWalletRequest struct {
	UserID      uuid.UUID `json:"user_id" binding:"required"`
}

func (h *WalletHandler) ProcessOperation(c *gin.Context) {
	var req WalletOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	operationType := entities.OperationType(req.OperationType)
	if operationType != entities.OperationTypeDeposit && operationType != entities.OperationTypeWithdraw {
		c.JSON(http.StatusBadRequest, gin.H{"error": "operationType must be DEPOSIT or WITHDRAW"})
		return
	}

	err := h.walletService.ProcessOperation(
		c.Request.Context(),
		req.WalletID,
		operationType,
		req.Amount,
	)

	if err != nil {
		switch err {
		case services.ErrWalletNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case services.ErrInsufficientFunds:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case services.ErrInvalidOperation, services.ErrInvalidAmount:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"details": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "operation completed successfully",
		"walletId": req.WalletID,
		"operationType": req.OperationType,
		"amount": req.Amount,
	})
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	walletIDStr := c.Param("walletId")
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}

	wallet, err := h.walletService.GetWallet(c.Request.Context(), walletID)
	if err != nil {
		switch err {
		case services.ErrWalletNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"details": err.Error(),
			})
		}
		return
	}

	response := WalletResponse{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: wallet.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

func (h *WalletHandler) CreateWallet(c *gin.Context) {
	var req CreateWalletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	w,err:= h.walletService.CreateWallet(c.Request.Context(), req.UserID)
		if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create wallet"})
		return 
	}
	
	c.JSON(http.StatusOK, w)
}






