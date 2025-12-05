package http

import (
	"github.com/gin-gonic/gin"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IUserHandler defines user-related handler methods
type IUserHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user and return a token
// @Tags user
// @Accept json
// @Produce json
// @Param request body request.Login true "Login credentials"
// @Success 200 {object} apiwrapper.APIResponse "Success response with token"
// @Failure 400 {object} apiwrapper.APIResponse "Bad request"
// @Failure 401 {object} apiwrapper.APIResponse "Unauthorized"
// @Router /internal/api/v1/user/login [post]
func (h *Handler) Login(ctx *gin.Context) {
	var req request.Login
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	token, err := h.userUsecase.Login(ctx, req)
	if err != nil {
		logger.EnhanceWith(ctx).Errorw("Login failed", "error", err, "email", req.Email)
		apiwrapper.SendUnauthorized(ctx, "Login failed")
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{
		"token": token,
	})
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags user
// @Accept json
// @Produce json
// @Param request body request.Register true "Registration information"
// @Success 200 {object} apiwrapper.APIResponse "Success response"
// @Failure 400 {object} apiwrapper.APIResponse "Bad request"
// @Failure 500 {object} apiwrapper.APIResponse "Internal server error"
// @Router /internal/api/v1/user/register [post]
func (h *Handler) Register(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.Register
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

	err := h.userUsecase.Register(ctx, req)

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
