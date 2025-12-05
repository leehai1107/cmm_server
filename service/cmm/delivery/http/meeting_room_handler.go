package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IMeetingRoomHandler defines meeting room-related handler methods
type IMeetingRoomHandler interface {
	CreateMeetingRoom(ctx *gin.Context)
	GetMeetingRoom(ctx *gin.Context)
	GetMeetingRoomsByCoffeeShop(ctx *gin.Context)
	GetAvailableMeetingRooms(ctx *gin.Context)
	UpdateMeetingRoom(ctx *gin.Context)
	DeleteMeetingRoom(ctx *gin.Context)
}

// CreateMeetingRoom godoc
// @Summary Create a new meeting room
// @Description Create a meeting room for a coffee shop
// @Tags meeting-room
// @Accept json
// @Produce json
// @Param request body request.CreateMeetingRoom true "Meeting room details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/meeting-room/create [post]
func (h *Handler) CreateMeetingRoom(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	var req request.CreateMeetingRoom
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	room, err := h.meetingRoomUsecase.CreateMeetingRoom(ctx, ownerID, req)
	if err != nil {
		log.Errorw("Failed to create meeting room", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, room)
}

// GetMeetingRoom godoc
// @Summary Get meeting room details
// @Description Get details of a specific meeting room
// @Tags meeting-room
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/meeting-room/{id} [get]
func (h *Handler) GetMeetingRoom(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	roomID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid room ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid room ID")
		return
	}

	room, err := h.meetingRoomUsecase.GetMeetingRoom(ctx, roomID)
	if err != nil {
		log.Errorw("Failed to get meeting room", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, room)
}

// GetMeetingRoomsByCoffeeShop godoc
// @Summary Get meeting rooms by coffee shop
// @Description Get all meeting rooms for a coffee shop
// @Tags meeting-room
// @Accept json
// @Produce json
// @Param shop_id path string true "Coffee Shop ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/meeting-room/shop/{shop_id} [get]
func (h *Handler) GetMeetingRoomsByCoffeeShop(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	shopID, err := uuid.Parse(ctx.Param("shop_id"))
	if err != nil {
		log.Errorw("Invalid shop ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid shop ID")
		return
	}

	rooms, err := h.meetingRoomUsecase.GetMeetingRoomsByCoffeeShop(ctx, shopID)
	if err != nil {
		log.Errorw("Failed to get meeting rooms", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get meeting rooms")
		return
	}

	apiwrapper.SendSuccess(ctx, rooms)
}

// GetAvailableMeetingRooms godoc
// @Summary Get available meeting rooms
// @Description Get all available meeting rooms for a coffee shop
// @Tags meeting-room
// @Accept json
// @Produce json
// @Param shop_id path string true "Coffee Shop ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/meeting-room/shop/{shop_id}/available [get]
func (h *Handler) GetAvailableMeetingRooms(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	shopID, err := uuid.Parse(ctx.Param("shop_id"))
	if err != nil {
		log.Errorw("Invalid shop ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid shop ID")
		return
	}

	rooms, err := h.meetingRoomUsecase.GetAvailableMeetingRooms(ctx, shopID)
	if err != nil {
		log.Errorw("Failed to get available meeting rooms", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get available meeting rooms")
		return
	}

	apiwrapper.SendSuccess(ctx, rooms)
}

// UpdateMeetingRoom godoc
// @Summary Update a meeting room
// @Description Update meeting room details
// @Tags meeting-room
// @Accept json
// @Produce json
// @Param request body request.UpdateMeetingRoom true "Meeting room details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/meeting-room/update [put]
func (h *Handler) UpdateMeetingRoom(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	var req request.UpdateMeetingRoom
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.meetingRoomUsecase.UpdateMeetingRoom(ctx, ownerID, req)
	if err != nil {
		log.Errorw("Failed to update meeting room", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Meeting room updated successfully"})
}

// DeleteMeetingRoom godoc
// @Summary Delete a meeting room
// @Description Delete a meeting room
// @Tags meeting-room
// @Accept json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/meeting-room/{id} [delete]
func (h *Handler) DeleteMeetingRoom(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	roomID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid room ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid room ID")
		return
	}

	err = h.meetingRoomUsecase.DeleteMeetingRoom(ctx, ownerID, roomID)
	if err != nil {
		log.Errorw("Failed to delete meeting room", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Meeting room deleted successfully"})
}
