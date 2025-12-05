package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IVoucherHandler defines voucher-related handler methods
type IVoucherHandler interface {
	CreateVoucher(ctx *gin.Context)
	GetVoucher(ctx *gin.Context)
	GetAllVouchers(ctx *gin.Context)
	GetValidVouchers(ctx *gin.Context)
	UpdateVoucher(ctx *gin.Context)
	DeleteVoucher(ctx *gin.Context)
	ApplyVoucher(ctx *gin.Context)
}

// CreateVoucher godoc
// @Summary Create a new voucher
// @Description Create a voucher (admin only)
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body request.CreateVoucher true "Voucher details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/create [post]
func (h *Handler) CreateVoucher(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.CreateVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	voucher, err := h.voucherUsecase.CreateVoucher(ctx, req)
	if err != nil {
		log.Errorw("Failed to create voucher", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, voucher)
}

// GetVoucher godoc
// @Summary Get voucher details
// @Description Get details of a specific voucher
// @Tags voucher
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/{id} [get]
func (h *Handler) GetVoucher(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	voucherID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid voucher ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid voucher ID")
		return
	}

	voucher, err := h.voucherUsecase.GetVoucher(ctx, voucherID)
	if err != nil {
		log.Errorw("Failed to get voucher", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, voucher)
}

// GetAllVouchers godoc
// @Summary Get all vouchers
// @Description Get all vouchers (admin only)
// @Tags voucher
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/all [get]
func (h *Handler) GetAllVouchers(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	vouchers, err := h.voucherUsecase.GetAllVouchers(ctx)
	if err != nil {
		log.Errorw("Failed to get vouchers", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get vouchers")
		return
	}

	apiwrapper.SendSuccess(ctx, vouchers)
}

// GetValidVouchers godoc
// @Summary Get valid vouchers
// @Description Get all currently valid vouchers
// @Tags voucher
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/valid [get]
func (h *Handler) GetValidVouchers(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	vouchers, err := h.voucherUsecase.GetValidVouchers(ctx)
	if err != nil {
		log.Errorw("Failed to get valid vouchers", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get valid vouchers")
		return
	}

	apiwrapper.SendSuccess(ctx, vouchers)
}

// UpdateVoucher godoc
// @Summary Update a voucher
// @Description Update voucher details (admin only)
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body request.UpdateVoucher true "Voucher details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/update [put]
func (h *Handler) UpdateVoucher(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.UpdateVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.voucherUsecase.UpdateVoucher(ctx, req)
	if err != nil {
		log.Errorw("Failed to update voucher", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Voucher updated successfully"})
}

// DeleteVoucher godoc
// @Summary Delete a voucher
// @Description Delete a voucher (admin only)
// @Tags voucher
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/{id} [delete]
func (h *Handler) DeleteVoucher(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	voucherID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid voucher ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid voucher ID")
		return
	}

	err = h.voucherUsecase.DeleteVoucher(ctx, voucherID)
	if err != nil {
		log.Errorw("Failed to delete voucher", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Voucher deleted successfully"})
}

// ApplyVoucher godoc
// @Summary Apply voucher
// @Description Calculate discount for a voucher code
// @Tags voucher
// @Accept json
// @Produce json
// @Param request body request.ApplyVoucher true "Voucher application"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/voucher/apply [post]
func (h *Handler) ApplyVoucher(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.ApplyVoucher
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	calculation, err := h.voucherUsecase.ApplyVoucher(ctx, req)
	if err != nil {
		log.Errorw("Failed to apply voucher", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, calculation)
}
