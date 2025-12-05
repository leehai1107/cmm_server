package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IBookingHandler defines booking-related handler methods
type IBookingHandler interface {
	CreateBooking(ctx *gin.Context)
	GetBooking(ctx *gin.Context)
	GetMyBookings(ctx *gin.Context)
	CancelBooking(ctx *gin.Context)
	GetRoomBookings(ctx *gin.Context)
}

// CreateBooking godoc
// @Summary Create a new booking
// @Description Create a booking for a meeting room
// @Tags booking
// @Accept json
// @Produce json
// @Param request body request.CreateBooking true "Booking details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/booking/create [post]
func (h *Handler) CreateBooking(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	// Get customer ID from context (set by auth middleware)
	customerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	customerID := customerIDStr.(uuid.UUID)

	var req request.CreateBooking
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	booking, err := h.bookingUsecase.CreateBooking(ctx, customerID, req)
	if err != nil {
		log.Errorw("Failed to create booking", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, booking)
}

// GetBooking godoc
// @Summary Get booking details
// @Description Get details of a specific booking
// @Tags booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/booking/{id} [get]
func (h *Handler) GetBooking(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	bookingID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid booking ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid booking ID")
		return
	}

	booking, err := h.bookingUsecase.GetBooking(ctx, bookingID)
	if err != nil {
		log.Errorw("Failed to get booking", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, booking)
}

// GetMyBookings godoc
// @Summary Get my bookings
// @Description Get all bookings for the current user
// @Tags booking
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/booking/my-bookings [get]
func (h *Handler) GetMyBookings(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	customerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	customerID := customerIDStr.(uuid.UUID)

	bookings, err := h.bookingUsecase.GetCustomerBookings(ctx, customerID)
	if err != nil {
		log.Errorw("Failed to get bookings", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get bookings")
		return
	}

	apiwrapper.SendSuccess(ctx, bookings)
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel an existing booking
// @Tags booking
// @Accept json
// @Produce json
// @Param id path string true "Booking ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/booking/{id}/cancel [post]
func (h *Handler) CancelBooking(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	customerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	customerID := customerIDStr.(uuid.UUID)

	bookingID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid booking ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid booking ID")
		return
	}

	err = h.bookingUsecase.CancelBooking(ctx, customerID, bookingID)
	if err != nil {
		log.Errorw("Failed to cancel booking", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Booking cancelled successfully"})
}

// GetRoomBookings godoc
// @Summary Get bookings for a room
// @Description Get all bookings for a specific meeting room
// @Tags booking
// @Accept json
// @Produce json
// @Param room_id path string true "Room ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/booking/room/{room_id} [get]
func (h *Handler) GetRoomBookings(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	roomID, err := uuid.Parse(ctx.Param("room_id"))
	if err != nil {
		log.Errorw("Invalid room ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid room ID")
		return
	}

	bookings, err := h.bookingUsecase.GetRoomBookings(ctx, roomID)
	if err != nil {
		log.Errorw("Failed to get room bookings", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get room bookings")
		return
	}

	apiwrapper.SendSuccess(ctx, bookings)
}
