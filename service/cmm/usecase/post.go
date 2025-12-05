package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
	"github.com/leehai1107/cmm_server/service/cmm/model/response"
	"github.com/leehai1107/cmm_server/service/cmm/repository"
	"gorm.io/gorm"
)

type IPostUsecase interface {
	// Shop Posts
	CreateShopPost(ctx context.Context, ownerID uuid.UUID, req request.CreateShopPost) (*response.ShopPostResponse, error)
	GetShopPost(ctx context.Context, postID uuid.UUID) (*response.ShopPostResponse, error)
	GetShopPostsByCoffeeShop(ctx context.Context, shopID uuid.UUID) ([]response.ShopPostResponse, error)
	GetAllShopPosts(ctx context.Context) ([]response.ShopPostResponse, error)
	UpdateShopPost(ctx context.Context, ownerID uuid.UUID, req request.UpdateShopPost) error
	DeleteShopPost(ctx context.Context, ownerID uuid.UUID, postID uuid.UUID) error

	// Internal Posts
	CreateInternalPost(ctx context.Context, adminID uuid.UUID, req request.CreateInternalPost) (*response.InternalPostResponse, error)
	GetInternalPost(ctx context.Context, postID uuid.UUID) (*response.InternalPostResponse, error)
	GetInternalPostsForUser(ctx context.Context, role string) ([]response.InternalPostResponse, error)
	GetAllInternalPosts(ctx context.Context) ([]response.InternalPostResponse, error)
	UpdateInternalPost(ctx context.Context, req request.UpdateInternalPost) error
	DeleteInternalPost(ctx context.Context, postID uuid.UUID) error
}

type postUsecase struct {
	postRepo       repository.IPostRepo
	coffeeShopRepo repository.ICoffeeShopRepo
	userRepo       repository.IUserRepo
}

func NewPostUsecase(
	postRepo repository.IPostRepo,
	coffeeShopRepo repository.ICoffeeShopRepo,
	userRepo repository.IUserRepo,
) IPostUsecase {
	return &postUsecase{
		postRepo:       postRepo,
		coffeeShopRepo: coffeeShopRepo,
		userRepo:       userRepo,
	}
}

// Shop Post methods
func (u *postUsecase) CreateShopPost(ctx context.Context, ownerID uuid.UUID, req request.CreateShopPost) (*response.ShopPostResponse, error) {
	logger.EnhanceWith(ctx).Info("CreateShopPost usecase called")

	// Verify coffee shop exists and belongs to owner
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(req.CoffeeShopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coffee shop not found")
		}
		return nil, err
	}

	if shop.OwnerID != ownerID {
		return nil, errors.New("unauthorized to create post for this coffee shop")
	}

	post := &entity.ShopPost{
		ID:           uuid.New(),
		CoffeeShopID: req.CoffeeShopID,
		Title:        req.Title,
		Content:      req.Content,
		PublishedAt:  time.Now(),
	}

	if err := u.postRepo.CreateShopPost(post); err != nil {
		return nil, err
	}

	return &response.ShopPostResponse{
		ID:           post.ID,
		CoffeeShopID: post.CoffeeShopID,
		ShopName:     shop.Name,
		Title:        post.Title,
		Content:      post.Content,
		PublishedAt:  post.PublishedAt,
	}, nil
}

func (u *postUsecase) GetShopPost(ctx context.Context, postID uuid.UUID) (*response.ShopPostResponse, error) {
	post, err := u.postRepo.GetShopPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	shop, _ := u.coffeeShopRepo.GetCoffeeShopByID(post.CoffeeShopID)
	shopName := ""
	if shop != nil {
		shopName = shop.Name
	}

	return &response.ShopPostResponse{
		ID:           post.ID,
		CoffeeShopID: post.CoffeeShopID,
		ShopName:     shopName,
		Title:        post.Title,
		Content:      post.Content,
		PublishedAt:  post.PublishedAt,
	}, nil
}

func (u *postUsecase) GetShopPostsByCoffeeShop(ctx context.Context, shopID uuid.UUID) ([]response.ShopPostResponse, error) {
	posts, err := u.postRepo.GetShopPostsByCoffeeShop(shopID)
	if err != nil {
		return nil, err
	}

	shop, _ := u.coffeeShopRepo.GetCoffeeShopByID(shopID)
	shopName := ""
	if shop != nil {
		shopName = shop.Name
	}

	var result []response.ShopPostResponse
	for _, post := range posts {
		result = append(result, response.ShopPostResponse{
			ID:           post.ID,
			CoffeeShopID: post.CoffeeShopID,
			ShopName:     shopName,
			Title:        post.Title,
			Content:      post.Content,
			PublishedAt:  post.PublishedAt,
		})
	}

	return result, nil
}

func (u *postUsecase) GetAllShopPosts(ctx context.Context) ([]response.ShopPostResponse, error) {
	posts, err := u.postRepo.GetAllShopPosts()
	if err != nil {
		return nil, err
	}

	var result []response.ShopPostResponse
	for _, post := range posts {
		shop, _ := u.coffeeShopRepo.GetCoffeeShopByID(post.CoffeeShopID)
		shopName := ""
		if shop != nil {
			shopName = shop.Name
		}

		result = append(result, response.ShopPostResponse{
			ID:           post.ID,
			CoffeeShopID: post.CoffeeShopID,
			ShopName:     shopName,
			Title:        post.Title,
			Content:      post.Content,
			PublishedAt:  post.PublishedAt,
		})
	}

	return result, nil
}

func (u *postUsecase) UpdateShopPost(ctx context.Context, ownerID uuid.UUID, req request.UpdateShopPost) error {
	post, err := u.postRepo.GetShopPostByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// Verify ownership
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(post.CoffeeShopID)
	if err != nil {
		return err
	}
	if shop.OwnerID != ownerID {
		return errors.New("unauthorized to update this post")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	return u.postRepo.UpdateShopPost(post)
}

func (u *postUsecase) DeleteShopPost(ctx context.Context, ownerID uuid.UUID, postID uuid.UUID) error {
	post, err := u.postRepo.GetShopPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	// Verify ownership
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(post.CoffeeShopID)
	if err != nil {
		return err
	}
	if shop.OwnerID != ownerID {
		return errors.New("unauthorized to delete this post")
	}

	return u.postRepo.DeleteShopPost(postID)
}

// Internal Post methods
func (u *postUsecase) CreateInternalPost(ctx context.Context, adminID uuid.UUID, req request.CreateInternalPost) (*response.InternalPostResponse, error) {
	logger.EnhanceWith(ctx).Info("CreateInternalPost usecase called")

	// Validate visibleTo value
	validVisibility := map[string]bool{"all": true, "admin": true, "owner": true, "customer": true}
	if !validVisibility[req.VisibleTo] {
		return nil, errors.New("invalid visible_to value")
	}

	post := &entity.InternalPost{
		ID:          uuid.New(),
		Title:       req.Title,
		Content:     req.Content,
		CreatedBy:   adminID,
		VisibleTo:   req.VisibleTo,
		PublishedAt: time.Now(),
	}

	if err := u.postRepo.CreateInternalPost(post); err != nil {
		return nil, err
	}

	creator, _ := u.userRepo.GetUserByEmail("") // Would need a GetUserByID method
	creatorName := ""
	if creator != nil {
		creatorName = creator.FullName
	}

	return &response.InternalPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		CreatedBy:   post.CreatedBy,
		CreatorName: creatorName,
		VisibleTo:   post.VisibleTo,
		PublishedAt: post.PublishedAt,
	}, nil
}

func (u *postUsecase) GetInternalPost(ctx context.Context, postID uuid.UUID) (*response.InternalPostResponse, error) {
	post, err := u.postRepo.GetInternalPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	return &response.InternalPostResponse{
		ID:          post.ID,
		Title:       post.Title,
		Content:     post.Content,
		CreatedBy:   post.CreatedBy,
		VisibleTo:   post.VisibleTo,
		PublishedAt: post.PublishedAt,
	}, nil
}

func (u *postUsecase) GetInternalPostsForUser(ctx context.Context, role string) ([]response.InternalPostResponse, error) {
	posts, err := u.postRepo.GetInternalPostsByVisibility(role)
	if err != nil {
		return nil, err
	}

	var result []response.InternalPostResponse
	for _, post := range posts {
		result = append(result, response.InternalPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Content:     post.Content,
			CreatedBy:   post.CreatedBy,
			VisibleTo:   post.VisibleTo,
			PublishedAt: post.PublishedAt,
		})
	}

	return result, nil
}

func (u *postUsecase) GetAllInternalPosts(ctx context.Context) ([]response.InternalPostResponse, error) {
	posts, err := u.postRepo.GetAllInternalPosts()
	if err != nil {
		return nil, err
	}

	var result []response.InternalPostResponse
	for _, post := range posts {
		result = append(result, response.InternalPostResponse{
			ID:          post.ID,
			Title:       post.Title,
			Content:     post.Content,
			CreatedBy:   post.CreatedBy,
			VisibleTo:   post.VisibleTo,
			PublishedAt: post.PublishedAt,
		})
	}

	return result, nil
}

func (u *postUsecase) UpdateInternalPost(ctx context.Context, req request.UpdateInternalPost) error {
	post, err := u.postRepo.GetInternalPostByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.VisibleTo != "" {
		validVisibility := map[string]bool{"all": true, "admin": true, "owner": true, "customer": true}
		if !validVisibility[req.VisibleTo] {
			return errors.New("invalid visible_to value")
		}
		post.VisibleTo = req.VisibleTo
	}

	return u.postRepo.UpdateInternalPost(post)
}

func (u *postUsecase) DeleteInternalPost(ctx context.Context, postID uuid.UUID) error {
	_, err := u.postRepo.GetInternalPostByID(postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("post not found")
		}
		return err
	}

	return u.postRepo.DeleteInternalPost(postID)
}
