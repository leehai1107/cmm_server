package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IWalletHandler defines wallet-related handler methods
type IWalletHandler interface {
	GetWallet(ctx *gin.Context)
	CreateTopup(ctx *gin.Context)
	ConfirmTopup(ctx *gin.Context)
	GetTopupHistory(ctx *gin.Context)
	GetTransactionHistory(ctx *gin.Context)
}

// GetWallet godoc
// @Summary Get wallet balance
// @Description Get wallet balance for the current user
// @Tags wallet
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/wallet [get]
func (h *Handler) GetWallet(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	userID := userIDStr.(uuid.UUID)

	wallet, err := h.walletUsecase.GetWallet(ctx, userID)
	if err != nil {
		log.Errorw("Failed to get wallet", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, wallet)
}

// CreateTopup godoc
// @Summary Create a top-up request
// @Description Create a top-up request for wallet
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body request.TopupWallet true "Top-up details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/wallet/topup [post]
func (h *Handler) CreateTopup(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	userID := userIDStr.(uuid.UUID)

	var req request.TopupWallet
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	topup, err := h.walletUsecase.CreateTopup(ctx, userID, req)
	if err != nil {
		log.Errorw("Failed to create top-up", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, topup)
}

// ConfirmTopup godoc
// @Summary Confirm a top-up
// @Description Confirm a top-up transaction (admin only)
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body request.ConfirmTopup true "Top-up confirmation"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/wallet/topup/confirm [post]
func (h *Handler) ConfirmTopup(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.ConfirmTopup
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.walletUsecase.ConfirmTopup(ctx, req)
	if err != nil {
		log.Errorw("Failed to confirm top-up", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Top-up confirmed successfully"})
}

// GetTopupHistory godoc
// @Summary Get top-up history
// @Description Get all top-up history for the current user
// @Tags wallet
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/wallet/topup/history [get]
func (h *Handler) GetTopupHistory(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	userID := userIDStr.(uuid.UUID)

	topups, err := h.walletUsecase.GetTopupHistory(ctx, userID)
	if err != nil {
		log.Errorw("Failed to get top-up history", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get top-up history")
		return
	}

	apiwrapper.SendSuccess(ctx, topups)
}

// GetTransactionHistory godoc
// @Summary Get transaction history
// @Description Get all transactions for the current user
// @Tags wallet
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/wallet/transactions [get]
func (h *Handler) GetTransactionHistory(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	userIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	userID := userIDStr.(uuid.UUID)

	transactions, err := h.walletUsecase.GetTransactionHistory(ctx, userID)
	if err != nil {
		log.Errorw("Failed to get transaction history", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get transaction history")
		return
	}

	apiwrapper.SendSuccess(ctx, transactions)
}
