package http

import (
	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IAdminHandler defines admin-related handler methods
type IAdminHandler interface {
	// Add admin-specific methods here
	CreateAccount(ctx *gin.Context)
}

// Add admin-specific handler implementations here
func (h *Handler) CreateAccount(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.CreateAccount
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	// Validate request fields
	if req.FirstName == "" || req.LastName == "" || req.Email == "" || req.Password == "" {
		log.Errorw("Missing required fields",
			"firstName", req.FirstName,
			"lastName", req.LastName,
			"email", req.Email)
		apiwrapper.SendBadRequest(ctx, "All fields are required")
		return
	}

	err := h.adminUsecase.CreateAccount(ctx, req)

	if err != nil {
		log.Errorw("Registration failed", "error", err, "email", req.Email)
		if err.Error() == "email already registered" || err.Error() == "email already exists" {
			apiwrapper.SendBadRequest(ctx, "Email already exists")
			return
		}
		apiwrapper.SendInternalError(ctx, "Registration failed")
		return
	}

	apiwrapper.SendSuccess(ctx, nil)
}
