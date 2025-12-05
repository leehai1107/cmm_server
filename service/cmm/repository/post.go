package repository

import (
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type IPostRepo interface {
	// Shop Posts (Public announcements from coffee shops)
	CreateShopPost(post *entity.ShopPost) error
	GetShopPostByID(id uuid.UUID) (*entity.ShopPost, error)
	GetShopPostsByCoffeeShop(shopID uuid.UUID) ([]entity.ShopPost, error)
	GetAllShopPosts() ([]entity.ShopPost, error)
	UpdateShopPost(post *entity.ShopPost) error
	DeleteShopPost(id uuid.UUID) error

	// Internal Posts (Platform announcements)
	CreateInternalPost(post *entity.InternalPost) error
	GetInternalPostByID(id uuid.UUID) (*entity.InternalPost, error)
	GetInternalPostsByVisibility(visibleTo string) ([]entity.InternalPost, error)
	GetAllInternalPosts() ([]entity.InternalPost, error)
	UpdateInternalPost(post *entity.InternalPost) error
	DeleteInternalPost(id uuid.UUID) error
}

type postRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) IPostRepo {
	return &postRepo{
		db: db,
	}
}

// Shop Post methods
func (r *postRepo) CreateShopPost(post *entity.ShopPost) error {
	logger.Info("CreateShopPost repository method called")
	return r.db.Create(post).Error
}

func (r *postRepo) GetShopPostByID(id uuid.UUID) (*entity.ShopPost, error) {
	logger.Info("GetShopPostByID repository method called")
	var post entity.ShopPost
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) GetShopPostsByCoffeeShop(shopID uuid.UUID) ([]entity.ShopPost, error) {
	logger.Info("GetShopPostsByCoffeeShop repository method called")
	var posts []entity.ShopPost
	err := r.db.Where("coffee_shop_id = ?", shopID).Order("published_at DESC").Find(&posts).Error
	return posts, err
}

func (r *postRepo) GetAllShopPosts() ([]entity.ShopPost, error) {
	logger.Info("GetAllShopPosts repository method called")
	var posts []entity.ShopPost
	err := r.db.Order("published_at DESC").Find(&posts).Error
	return posts, err
}

func (r *postRepo) UpdateShopPost(post *entity.ShopPost) error {
	logger.Info("UpdateShopPost repository method called")
	return r.db.Save(post).Error
}

func (r *postRepo) DeleteShopPost(id uuid.UUID) error {
	logger.Info("DeleteShopPost repository method called")
	return r.db.Delete(&entity.ShopPost{}, "id = ?", id).Error
}

// Internal Post methods
func (r *postRepo) CreateInternalPost(post *entity.InternalPost) error {
	logger.Info("CreateInternalPost repository method called")
	return r.db.Create(post).Error
}

func (r *postRepo) GetInternalPostByID(id uuid.UUID) (*entity.InternalPost, error) {
	logger.Info("GetInternalPostByID repository method called")
	var post entity.InternalPost
	err := r.db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) GetInternalPostsByVisibility(visibleTo string) ([]entity.InternalPost, error) {
	logger.Info("GetInternalPostsByVisibility repository method called")
	var posts []entity.InternalPost
	err := r.db.Where("visible_to = ? OR visible_to = ?", visibleTo, "all").
		Order("published_at DESC").Find(&posts).Error
	return posts, err
}

func (r *postRepo) GetAllInternalPosts() ([]entity.InternalPost, error) {
	logger.Info("GetAllInternalPosts repository method called")
	var posts []entity.InternalPost
	err := r.db.Order("published_at DESC").Find(&posts).Error
	return posts, err
}

func (r *postRepo) UpdateInternalPost(post *entity.InternalPost) error {
	logger.Info("UpdateInternalPost repository method called")
	return r.db.Save(post).Error
}

func (r *postRepo) DeleteInternalPost(id uuid.UUID) error {
	logger.Info("DeleteInternalPost repository method called")
	return r.db.Delete(&entity.InternalPost{}, "id = ?", id).Error
}
