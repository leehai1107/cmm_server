package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// ICoffeeShopHandler defines coffee shop-related handler methods
type ICoffeeShopHandler interface {
	CreateCoffeeShop(ctx *gin.Context)
	GetCoffeeShop(ctx *gin.Context)
	GetMyCoffeeShops(ctx *gin.Context)
	GetAllCoffeeShops(ctx *gin.Context)
	UpdateCoffeeShop(ctx *gin.Context)
	DeleteCoffeeShop(ctx *gin.Context)
	SetCommissionRate(ctx *gin.Context)
	GetCommissionRate(ctx *gin.Context)
}

// CreateCoffeeShop godoc
// @Summary Create a new coffee shop
// @Description Create a coffee shop (owner only)
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Param request body request.CreateCoffeeShop true "Coffee shop details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/create [post]
func (h *Handler) CreateCoffeeShop(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	var req request.CreateCoffeeShop
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	shop, err := h.coffeeShopUsecase.CreateCoffeeShop(ctx, ownerID, req)
	if err != nil {
		log.Errorw("Failed to create coffee shop", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, shop)
}

// GetCoffeeShop godoc
// @Summary Get coffee shop details
// @Description Get details of a specific coffee shop
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/{id} [get]
func (h *Handler) GetCoffeeShop(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid shop ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid shop ID")
		return
	}

	shop, err := h.coffeeShopUsecase.GetCoffeeShop(ctx, shopID)
	if err != nil {
		log.Errorw("Failed to get coffee shop", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, shop)
}

// GetMyCoffeeShops godoc
// @Summary Get my coffee shops
// @Description Get all coffee shops owned by current user
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/my-shops [get]
func (h *Handler) GetMyCoffeeShops(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	shops, err := h.coffeeShopUsecase.GetCoffeeShopsByOwner(ctx, ownerID)
	if err != nil {
		log.Errorw("Failed to get coffee shops", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get coffee shops")
		return
	}

	apiwrapper.SendSuccess(ctx, shops)
}

// GetAllCoffeeShops godoc
// @Summary Get all coffee shops
// @Description Get all coffee shops
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/all [get]
func (h *Handler) GetAllCoffeeShops(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	shops, err := h.coffeeShopUsecase.GetAllCoffeeShops(ctx)
	if err != nil {
		log.Errorw("Failed to get coffee shops", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get coffee shops")
		return
	}

	apiwrapper.SendSuccess(ctx, shops)
}

// UpdateCoffeeShop godoc
// @Summary Update a coffee shop
// @Description Update coffee shop details
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Param request body request.UpdateCoffeeShop true "Coffee shop details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/update [put]
func (h *Handler) UpdateCoffeeShop(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	var req request.UpdateCoffeeShop
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.coffeeShopUsecase.UpdateCoffeeShop(ctx, ownerID, req)
	if err != nil {
		log.Errorw("Failed to update coffee shop", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Coffee shop updated successfully"})
}

// DeleteCoffeeShop godoc
// @Summary Delete a coffee shop
// @Description Delete a coffee shop
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Param id path string true "Coffee Shop ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/{id} [delete]
func (h *Handler) DeleteCoffeeShop(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	shopID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid shop ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid shop ID")
		return
	}

	err = h.coffeeShopUsecase.DeleteCoffeeShop(ctx, ownerID, shopID)
	if err != nil {
		log.Errorw("Failed to delete coffee shop", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Coffee shop deleted successfully"})
}

// SetCommissionRate godoc
// @Summary Set commission rate
// @Description Set commission rate for a coffee shop (admin only)
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Param request body request.SetCommissionRate true "Commission rate details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/commission/set [post]
func (h *Handler) SetCommissionRate(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.SetCommissionRate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.coffeeShopUsecase.SetCommissionRate(ctx, req)
	if err != nil {
		log.Errorw("Failed to set commission rate", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Commission rate set successfully"})
}

// GetCommissionRate godoc
// @Summary Get commission rate
// @Description Get commission rate for a coffee shop
// @Tags coffee-shop
// @Accept json
// @Produce json
// @Param shop_id path string true "Coffee Shop ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/coffee-shop/commission/{shop_id} [get]
func (h *Handler) GetCommissionRate(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	shopID, err := uuid.Parse(ctx.Param("shop_id"))
	if err != nil {
		log.Errorw("Invalid shop ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid shop ID")
		return
	}

	rate, err := h.coffeeShopUsecase.GetCommissionRate(ctx, shopID)
	if err != nil {
		log.Errorw("Failed to get commission rate", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, rate)
}
