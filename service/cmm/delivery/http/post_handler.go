package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/apiwrapper"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
)

// IPostHandler defines post-related handler methods
type IPostHandler interface {
	// Shop Posts
	CreateShopPost(ctx *gin.Context)
	GetShopPost(ctx *gin.Context)
	GetShopPostsByCoffeeShop(ctx *gin.Context)
	GetAllShopPosts(ctx *gin.Context)
	UpdateShopPost(ctx *gin.Context)
	DeleteShopPost(ctx *gin.Context)

	// Internal Posts
	CreateInternalPost(ctx *gin.Context)
	GetInternalPost(ctx *gin.Context)
	GetMyInternalPosts(ctx *gin.Context)
	GetAllInternalPosts(ctx *gin.Context)
	UpdateInternalPost(ctx *gin.Context)
	DeleteInternalPost(ctx *gin.Context)
}

// Shop Post handlers

// CreateShopPost godoc
// @Summary Create a shop post
// @Description Create a post for a coffee shop
// @Tags post
// @Accept json
// @Produce json
// @Param request body request.CreateShopPost true "Post details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/shop/create [post]
func (h *Handler) CreateShopPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	var req request.CreateShopPost
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	post, err := h.postUsecase.CreateShopPost(ctx, ownerID, req)
	if err != nil {
		log.Errorw("Failed to create shop post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, post)
}

// GetShopPost godoc
// @Summary Get shop post details
// @Description Get details of a specific shop post
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/shop/{id} [get]
func (h *Handler) GetShopPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	postID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid post ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid post ID")
		return
	}

	post, err := h.postUsecase.GetShopPost(ctx, postID)
	if err != nil {
		log.Errorw("Failed to get shop post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, post)
}

// GetShopPostsByCoffeeShop godoc
// @Summary Get shop posts by coffee shop
// @Description Get all posts for a coffee shop
// @Tags post
// @Accept json
// @Produce json
// @Param shop_id path string true "Coffee Shop ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/shop/coffee-shop/{shop_id} [get]
func (h *Handler) GetShopPostsByCoffeeShop(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	shopID, err := uuid.Parse(ctx.Param("shop_id"))
	if err != nil {
		log.Errorw("Invalid shop ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid shop ID")
		return
	}

	posts, err := h.postUsecase.GetShopPostsByCoffeeShop(ctx, shopID)
	if err != nil {
		log.Errorw("Failed to get shop posts", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get shop posts")
		return
	}

	apiwrapper.SendSuccess(ctx, posts)
}

// GetAllShopPosts godoc
// @Summary Get all shop posts
// @Description Get all shop posts
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/shop/all [get]
func (h *Handler) GetAllShopPosts(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	posts, err := h.postUsecase.GetAllShopPosts(ctx)
	if err != nil {
		log.Errorw("Failed to get shop posts", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get shop posts")
		return
	}

	apiwrapper.SendSuccess(ctx, posts)
}

// UpdateShopPost godoc
// @Summary Update a shop post
// @Description Update shop post details
// @Tags post
// @Accept json
// @Produce json
// @Param request body request.UpdateShopPost true "Post details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/shop/update [put]
func (h *Handler) UpdateShopPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	var req request.UpdateShopPost
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.postUsecase.UpdateShopPost(ctx, ownerID, req)
	if err != nil {
		log.Errorw("Failed to update shop post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Shop post updated successfully"})
}

// DeleteShopPost godoc
// @Summary Delete a shop post
// @Description Delete a shop post
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/shop/{id} [delete]
func (h *Handler) DeleteShopPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	ownerIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	ownerID := ownerIDStr.(uuid.UUID)

	postID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid post ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid post ID")
		return
	}

	err = h.postUsecase.DeleteShopPost(ctx, ownerID, postID)
	if err != nil {
		log.Errorw("Failed to delete shop post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Shop post deleted successfully"})
}

// Internal Post handlers

// CreateInternalPost godoc
// @Summary Create an internal post
// @Description Create an internal post (admin only)
// @Tags post
// @Accept json
// @Produce json
// @Param request body request.CreateInternalPost true "Post details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/internal/create [post]
func (h *Handler) CreateInternalPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	adminIDStr, exists := ctx.Get("user_id")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	adminID := adminIDStr.(uuid.UUID)

	var req request.CreateInternalPost
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	post, err := h.postUsecase.CreateInternalPost(ctx, adminID, req)
	if err != nil {
		log.Errorw("Failed to create internal post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, post)
}

// GetInternalPost godoc
// @Summary Get internal post details
// @Description Get details of a specific internal post
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/internal/{id} [get]
func (h *Handler) GetInternalPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	postID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid post ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid post ID")
		return
	}

	post, err := h.postUsecase.GetInternalPost(ctx, postID)
	if err != nil {
		log.Errorw("Failed to get internal post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, post)
}

// GetMyInternalPosts godoc
// @Summary Get internal posts for current user
// @Description Get internal posts visible to current user role
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/internal/my-posts [get]
func (h *Handler) GetMyInternalPosts(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	roleStr, exists := ctx.Get("role")
	if !exists {
		apiwrapper.SendUnauthorized(ctx, "User not authenticated")
		return
	}
	role := roleStr.(string)

	posts, err := h.postUsecase.GetInternalPostsForUser(ctx, role)
	if err != nil {
		log.Errorw("Failed to get internal posts", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get internal posts")
		return
	}

	apiwrapper.SendSuccess(ctx, posts)
}

// GetAllInternalPosts godoc
// @Summary Get all internal posts
// @Description Get all internal posts (admin only)
// @Tags post
// @Accept json
// @Produce json
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/internal/all [get]
func (h *Handler) GetAllInternalPosts(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	posts, err := h.postUsecase.GetAllInternalPosts(ctx)
	if err != nil {
		log.Errorw("Failed to get internal posts", "error", err)
		apiwrapper.SendInternalError(ctx, "Failed to get internal posts")
		return
	}

	apiwrapper.SendSuccess(ctx, posts)
}

// UpdateInternalPost godoc
// @Summary Update an internal post
// @Description Update internal post details (admin only)
// @Tags post
// @Accept json
// @Produce json
// @Param request body request.UpdateInternalPost true "Post details"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/internal/update [put]
func (h *Handler) UpdateInternalPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	var req request.UpdateInternalPost
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Errorw("Invalid request format", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid request format")
		return
	}

	err := h.postUsecase.UpdateInternalPost(ctx, req)
	if err != nil {
		log.Errorw("Failed to update internal post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Internal post updated successfully"})
}

// DeleteInternalPost godoc
// @Summary Delete an internal post
// @Description Delete an internal post (admin only)
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "Post ID"
// @Success 200 {object} apiwrapper.APIResponse
// @Failure 400 {object} apiwrapper.APIResponse
// @Router /api/v1/post/internal/{id} [delete]
func (h *Handler) DeleteInternalPost(ctx *gin.Context) {
	log := logger.EnhanceWith(ctx)

	postID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		log.Errorw("Invalid post ID", "error", err)
		apiwrapper.SendBadRequest(ctx, "Invalid post ID")
		return
	}

	err = h.postUsecase.DeleteInternalPost(ctx, postID)
	if err != nil {
		log.Errorw("Failed to delete internal post", "error", err)
		apiwrapper.SendBadRequest(ctx, err.Error())
		return
	}

	apiwrapper.SendSuccess(ctx, gin.H{"message": "Internal post deleted successfully"})
}
