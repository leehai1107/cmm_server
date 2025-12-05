package http

import (
	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/websocket"
)

// IChatHandler defines chat-related handler methods
type IChatHandler interface {
	ServeWS(ctx *gin.Context)
}

// ServeWS godoc
// @Summary WebSocket connection
// @Description Establish a WebSocket connection for real-time chat
// @Tags chat
// @Param roomId path string true "Room ID"
// @Success 101 "Switching Protocols to WebSocket"
// @Failure 400 {object} apiwrapper.APIResponse "Bad request"
// @Router /internal/api/v1/chat/ws/{roomId} [get]
func (h *Handler) ServeWS(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	if roomId == "" {
		apiwrapper.SendBadRequest(ctx, "Room ID is required")
		return
	}
	websocket.ServeWs(ctx, roomId)
}
