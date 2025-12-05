package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"github.com/leehai1107/cmm_server/service/cmm/model/request"
	"github.com/leehai1107/cmm_server/service/cmm/model/response"
	"github.com/leehai1107/cmm_server/service/cmm/repository"
	"gorm.io/gorm"
)

type IMeetingRoomUsecase interface {
	CreateMeetingRoom(ctx context.Context, ownerID uuid.UUID, req request.CreateMeetingRoom) (*response.MeetingRoomResponse, error)
	GetMeetingRoom(ctx context.Context, roomID uuid.UUID) (*response.MeetingRoomResponse, error)
	GetMeetingRoomsByCoffeeShop(ctx context.Context, shopID uuid.UUID) ([]response.MeetingRoomResponse, error)
	GetAvailableMeetingRooms(ctx context.Context, shopID uuid.UUID) ([]response.MeetingRoomResponse, error)
	UpdateMeetingRoom(ctx context.Context, ownerID uuid.UUID, req request.UpdateMeetingRoom) error
	DeleteMeetingRoom(ctx context.Context, ownerID uuid.UUID, roomID uuid.UUID) error
}

type meetingRoomUsecase struct {
	meetingRoomRepo repository.IMeetingRoomRepo
	coffeeShopRepo  repository.ICoffeeShopRepo
}

func NewMeetingRoomUsecase(
	meetingRoomRepo repository.IMeetingRoomRepo,
	coffeeShopRepo repository.ICoffeeShopRepo,
) IMeetingRoomUsecase {
	return &meetingRoomUsecase{
		meetingRoomRepo: meetingRoomRepo,
		coffeeShopRepo:  coffeeShopRepo,
	}
}

func (u *meetingRoomUsecase) CreateMeetingRoom(ctx context.Context, ownerID uuid.UUID, req request.CreateMeetingRoom) (*response.MeetingRoomResponse, error) {
	logger.EnhanceWith(ctx).Info("CreateMeetingRoom usecase called")

	// Verify coffee shop exists and belongs to owner
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(req.CoffeeShopID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("coffee shop not found")
		}
		return nil, err
	}

	if shop.OwnerID != ownerID {
		return nil, errors.New("unauthorized to create room for this coffee shop")
	}

	room := &entity.MeetingRoom{
		ID:           uuid.New(),
		CoffeeShopID: req.CoffeeShopID,
		Name:         req.Name,
		Capacity:     req.Capacity,
		PricePerHour: req.PricePerHour,
		Available:    true,
	}

	if err := u.meetingRoomRepo.CreateMeetingRoom(room); err != nil {
		return nil, err
	}

	return &response.MeetingRoomResponse{
		ID:           room.ID,
		CoffeeShopID: room.CoffeeShopID,
		ShopName:     shop.Name,
		Name:         room.Name,
		Capacity:     room.Capacity,
		PricePerHour: room.PricePerHour,
		Available:    room.Available,
	}, nil
}

func (u *meetingRoomUsecase) GetMeetingRoom(ctx context.Context, roomID uuid.UUID) (*response.MeetingRoomResponse, error) {
	room, err := u.meetingRoomRepo.GetMeetingRoomByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("meeting room not found")
		}
		return nil, err
	}

	shop, _ := u.coffeeShopRepo.GetCoffeeShopByID(room.CoffeeShopID)
	shopName := ""
	if shop != nil {
		shopName = shop.Name
	}

	return &response.MeetingRoomResponse{
		ID:           room.ID,
		CoffeeShopID: room.CoffeeShopID,
		ShopName:     shopName,
		Name:         room.Name,
		Capacity:     room.Capacity,
		PricePerHour: room.PricePerHour,
		Available:    room.Available,
	}, nil
}

func (u *meetingRoomUsecase) GetMeetingRoomsByCoffeeShop(ctx context.Context, shopID uuid.UUID) ([]response.MeetingRoomResponse, error) {
	rooms, err := u.meetingRoomRepo.GetMeetingRoomsByCoffeeShop(shopID)
	if err != nil {
		return nil, err
	}

	shop, _ := u.coffeeShopRepo.GetCoffeeShopByID(shopID)
	shopName := ""
	if shop != nil {
		shopName = shop.Name
	}

	var result []response.MeetingRoomResponse
	for _, room := range rooms {
		result = append(result, response.MeetingRoomResponse{
			ID:           room.ID,
			CoffeeShopID: room.CoffeeShopID,
			ShopName:     shopName,
			Name:         room.Name,
			Capacity:     room.Capacity,
			PricePerHour: room.PricePerHour,
			Available:    room.Available,
		})
	}

	return result, nil
}

func (u *meetingRoomUsecase) GetAvailableMeetingRooms(ctx context.Context, shopID uuid.UUID) ([]response.MeetingRoomResponse, error) {
	rooms, err := u.meetingRoomRepo.GetAvailableMeetingRooms(shopID)
	if err != nil {
		return nil, err
	}

	shop, _ := u.coffeeShopRepo.GetCoffeeShopByID(shopID)
	shopName := ""
	if shop != nil {
		shopName = shop.Name
	}

	var result []response.MeetingRoomResponse
	for _, room := range rooms {
		result = append(result, response.MeetingRoomResponse{
			ID:           room.ID,
			CoffeeShopID: room.CoffeeShopID,
			ShopName:     shopName,
			Name:         room.Name,
			Capacity:     room.Capacity,
			PricePerHour: room.PricePerHour,
			Available:    room.Available,
		})
	}

	return result, nil
}

func (u *meetingRoomUsecase) UpdateMeetingRoom(ctx context.Context, ownerID uuid.UUID, req request.UpdateMeetingRoom) error {
	room, err := u.meetingRoomRepo.GetMeetingRoomByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("meeting room not found")
		}
		return err
	}

	// Verify ownership
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(room.CoffeeShopID)
	if err != nil {
		return err
	}
	if shop.OwnerID != ownerID {
		return errors.New("unauthorized to update this meeting room")
	}

	// Update fields
	if req.Name != "" {
		room.Name = req.Name
	}
	if req.Capacity > 0 {
		room.Capacity = req.Capacity
	}
	if req.PricePerHour > 0 {
		room.PricePerHour = req.PricePerHour
	}
	if req.Available != nil {
		room.Available = *req.Available
	}

	return u.meetingRoomRepo.UpdateMeetingRoom(room)
}

func (u *meetingRoomUsecase) DeleteMeetingRoom(ctx context.Context, ownerID uuid.UUID, roomID uuid.UUID) error {
	room, err := u.meetingRoomRepo.GetMeetingRoomByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("meeting room not found")
		}
		return err
	}

	// Verify ownership
	shop, err := u.coffeeShopRepo.GetCoffeeShopByID(room.CoffeeShopID)
	if err != nil {
		return err
	}
	if shop.OwnerID != ownerID {
		return errors.New("unauthorized to delete this meeting room")
	}

	return u.meetingRoomRepo.DeleteMeetingRoom(roomID)
}
