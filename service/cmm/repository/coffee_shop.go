package repository

import (
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type ICoffeeShopRepo interface {
	CreateCoffeeShop(shop *entity.CoffeeShop) error
	GetCoffeeShopByID(id uuid.UUID) (*entity.CoffeeShop, error)
	GetCoffeeShopsByOwner(ownerID uuid.UUID) ([]entity.CoffeeShop, error)
	GetAllCoffeeShops() ([]entity.CoffeeShop, error)
	UpdateCoffeeShop(shop *entity.CoffeeShop) error
	DeleteCoffeeShop(id uuid.UUID) error

	// Commission Rate methods
	SetCommissionRate(rate *entity.CommissionRate) error
	GetCommissionRate(shopID uuid.UUID) (*entity.CommissionRate, error)
}

type coffeeShopRepo struct {
	db *gorm.DB
}

func NewCoffeeShopRepo(db *gorm.DB) ICoffeeShopRepo {
	return &coffeeShopRepo{
		db: db,
	}
}

func (r *coffeeShopRepo) CreateCoffeeShop(shop *entity.CoffeeShop) error {
	logger.Info("CreateCoffeeShop repository method called")
	return r.db.Create(shop).Error
}

func (r *coffeeShopRepo) GetCoffeeShopByID(id uuid.UUID) (*entity.CoffeeShop, error) {
	logger.Info("GetCoffeeShopByID repository method called")
	var shop entity.CoffeeShop
	err := r.db.Where("id = ?", id).First(&shop).Error
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

func (r *coffeeShopRepo) GetCoffeeShopsByOwner(ownerID uuid.UUID) ([]entity.CoffeeShop, error) {
	logger.Info("GetCoffeeShopsByOwner repository method called")
	var shops []entity.CoffeeShop
	err := r.db.Where("owner_id = ?", ownerID).Find(&shops).Error
	return shops, err
}

func (r *coffeeShopRepo) GetAllCoffeeShops() ([]entity.CoffeeShop, error) {
	logger.Info("GetAllCoffeeShops repository method called")
	var shops []entity.CoffeeShop
	err := r.db.Find(&shops).Error
	return shops, err
}

func (r *coffeeShopRepo) UpdateCoffeeShop(shop *entity.CoffeeShop) error {
	logger.Info("UpdateCoffeeShop repository method called")
	return r.db.Save(shop).Error
}

func (r *coffeeShopRepo) DeleteCoffeeShop(id uuid.UUID) error {
	logger.Info("DeleteCoffeeShop repository method called")
	return r.db.Delete(&entity.CoffeeShop{}, "id = ?", id).Error
}

func (r *coffeeShopRepo) SetCommissionRate(rate *entity.CommissionRate) error {
	logger.Info("SetCommissionRate repository method called")
	return r.db.Save(rate).Error
}

func (r *coffeeShopRepo) GetCommissionRate(shopID uuid.UUID) (*entity.CommissionRate, error) {
	logger.Info("GetCommissionRate repository method called")
	var rate entity.CommissionRate
	err := r.db.Where("coffee_shop_id = ?", shopID).First(&rate).Error
	if err != nil {
		return nil, err
	}
	return &rate, nil
}
