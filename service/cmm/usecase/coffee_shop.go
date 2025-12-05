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

type ICoffeeShopUsecase interface {
	CreateCoffeeShop(ctx context.Context, ownerID uuid.UUID, req request.CreateCoffeeShop) (*response.CoffeeShopResponse, error)
	GetCoffeeShop(ctx context.Context, shopID uuid.UUID) (*response.CoffeeShopResponse, error)
	GetCoffeeShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]response.CoffeeShopResponse, error)
	GetAllCoffeeShops(ctx context.Context) ([]response.CoffeeShopResponse, error)
	UpdateCoffeeShop(ctx context.Context, ownerID uuid.UUID, req request.UpdateCoffeeShop) error
	DeleteCoffeeShop(ctx context.Context, ownerID uuid.UUID, shopID uuid.UUID) error

	SetCommissionRate(ctx context.Context, req request.SetCommissionRate) error
	GetCommissionRate(ctx context.Context, shopID uuid.UUID) (*response.CommissionRateResponse, error)
}

type coffeeShopUsecase struct {
	coffeeShopRepo repository.ICoffeeShopRepo
}

func NewCoffeeShopUsecase(coffeeShopRepo repository.ICoffeeShopRepo) ICoffeeShopUsecase {
	return &coffeeShopUsecase{
		coffeeShopRepo: coffeeShopRepo,
	}
}

func (u *coffeeShopUsecase) CreateCoffeeShop(ctx context.Context, ownerID uuid.UUID, req request.CreateCoffeeShop) (*response.CoffeeShopResponse, error) {
	logger.EnhanceWith(ctx).Info("CreateCoffeeShop usecase called")

	shop := &entity.CoffeeShop{
		ID:          uuid.New(),
		OwnerID:     ownerID,
		Name:        req.Name,
		Location:    req.Location,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := u.coffeeShopRepo.CreateCoffeeShop(shop); err != nil {
		return nil, err
	}

	return &response.CoffeeShopResponse{
		ID:          shop.ID,
		OwnerID:     shop.OwnerID,
		Name:        shop.Name,
		Location:    shop.Location,
		Description: shop.Description,
		CreatedAt:   shop.CreatedAt,
	}, nil
}

func (u *coffeeShopUsecase) GetCoffeeShop(ctx context.Context, shopID uuid.UUID) (*response.CoffeeShopResponse, error) {
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coffee shop not found")
		}
		return nil, err
	}

	return &response.CoffeeShopResponse{
		ID:          shop.ID,
		OwnerID:     shop.OwnerID,
		Name:        shop.Name,
		Location:    shop.Location,
		Description: shop.Description,
		CreatedAt:   shop.CreatedAt,
	}, nil
}

func (u *coffeeShopUsecase) GetCoffeeShopsByOwner(ctx context.Context, ownerID uuid.UUID) ([]response.CoffeeShopResponse, error) {
	shops, err := u.coffeeShopRepo.GetCoffeeShopsByOwner(ownerID)
	if err != nil {
		return nil, err
	}

	var result []response.CoffeeShopResponse
	for _, shop := range shops {
		result = append(result, response.CoffeeShopResponse{
			ID:          shop.ID,
			OwnerID:     shop.OwnerID,
			Name:        shop.Name,
			Location:    shop.Location,
			Description: shop.Description,
			CreatedAt:   shop.CreatedAt,
		})
	}

	return result, nil
}

func (u *coffeeShopUsecase) GetAllCoffeeShops(ctx context.Context) ([]response.CoffeeShopResponse, error) {
	shops, err := u.coffeeShopRepo.GetAllCoffeeShops()
	if err != nil {
		return nil, err
	}

	var result []response.CoffeeShopResponse
	for _, shop := range shops {
		result = append(result, response.CoffeeShopResponse{
			ID:          shop.ID,
			OwnerID:     shop.OwnerID,
			Name:        shop.Name,
			Location:    shop.Location,
			Description: shop.Description,
			CreatedAt:   shop.CreatedAt,
		})
	}

	return result, nil
}

func (u *coffeeShopUsecase) UpdateCoffeeShop(ctx context.Context, ownerID uuid.UUID, req request.UpdateCoffeeShop) error {
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("coffee shop not found")
		}
		return err
	}

	if shop.OwnerID != ownerID {
		return errors.New("unauthorized to update this coffee shop")
	}

	if req.Name != "" {
		shop.Name = req.Name
	}
	if req.Location != "" {
		shop.Location = req.Location
	}
	if req.Description != "" {
		shop.Description = req.Description
	}

	return u.coffeeShopRepo.UpdateCoffeeShop(shop)
}

func (u *coffeeShopUsecase) DeleteCoffeeShop(ctx context.Context, ownerID uuid.UUID, shopID uuid.UUID) error {
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("coffee shop not found")
		}
		return err
	}

	if shop.OwnerID != ownerID {
		return errors.New("unauthorized to delete this coffee shop")
	}

	return u.coffeeShopRepo.DeleteCoffeeShop(shopID)
}

func (u *coffeeShopUsecase) SetCommissionRate(ctx context.Context, req request.SetCommissionRate) error {
	logger.EnhanceWith(ctx).Info("SetCommissionRate usecase called")

	// Verify coffee shop exists
	_, err := u.coffeeShopRepo.GetCoffeeShopByID(req.CoffeeShopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("coffee shop not found")
		}
		return err
	}

	// Check if rate already exists
	existingRate, err := u.coffeeShopRepo.GetCommissionRate(req.CoffeeShopID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existingRate != nil {
		// Update existing rate
		existingRate.RatePercent = req.RatePercent
		return u.coffeeShopRepo.SetCommissionRate(existingRate)
	}

	// Create new rate
	rate := &entity.CommissionRate{
		CoffeeShopID: req.CoffeeShopID,
		RatePercent:  req.RatePercent,
	}

	return u.coffeeShopRepo.SetCommissionRate(rate)
}

func (u *coffeeShopUsecase) GetCommissionRate(ctx context.Context, shopID uuid.UUID) (*response.CommissionRateResponse, error) {
	rate, err := u.coffeeShopRepo.GetCommissionRate(shopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("commission rate not found")
		}
		return nil, err
	}

	return &response.CommissionRateResponse{
		ID:           rate.ID,
		CoffeeShopID: rate.CoffeeShopID,
		RatePercent:  rate.RatePercent,
	}, nil
}
