package repository

import (
	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type IMeetingRoomRepo interface {
	CreateMeetingRoom(room *entity.MeetingRoom) error
	GetMeetingRoomByID(id uuid.UUID) (*entity.MeetingRoom, error)
	GetMeetingRoomsByCoffeeShop(shopID uuid.UUID) ([]entity.MeetingRoom, error)
	GetAvailableMeetingRooms(shopID uuid.UUID) ([]entity.MeetingRoom, error)
	UpdateMeetingRoom(room *entity.MeetingRoom) error
	UpdateRoomAvailability(id uuid.UUID, available bool) error
	DeleteMeetingRoom(id uuid.UUID) error
}

type meetingRoomRepo struct {
	db *gorm.DB
}

func NewMeetingRoomRepo(db *gorm.DB) IMeetingRoomRepo {
	return &meetingRoomRepo{
		db: db,
	}
}

func (r *meetingRoomRepo) CreateMeetingRoom(room *entity.MeetingRoom) error {
	logger.Info("CreateMeetingRoom repository method called")
	return r.db.Create(room).Error
}

func (r *meetingRoomRepo) GetMeetingRoomByID(id uuid.UUID) (*entity.MeetingRoom, error) {
	logger.Info("GetMeetingRoomByID repository method called")
	var room entity.MeetingRoom
	err := r.db.Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *meetingRoomRepo) GetMeetingRoomsByCoffeeShop(shopID uuid.UUID) ([]entity.MeetingRoom, error) {
	logger.Info("GetMeetingRoomsByCoffeeShop repository method called")
	var rooms []entity.MeetingRoom
	err := r.db.Where("coffee_shop_id = ?", shopID).Find(&rooms).Error
	return rooms, err
}

func (r *meetingRoomRepo) GetAvailableMeetingRooms(shopID uuid.UUID) ([]entity.MeetingRoom, error) {
	logger.Info("GetAvailableMeetingRooms repository method called")
	var rooms []entity.MeetingRoom
	err := r.db.Where("coffee_shop_id = ? AND available = ?", shopID, true).Find(&rooms).Error
	return rooms, err
}

func (r *meetingRoomRepo) UpdateMeetingRoom(room *entity.MeetingRoom) error {
	logger.Info("UpdateMeetingRoom repository method called")
	return r.db.Save(room).Error
}

func (r *meetingRoomRepo) UpdateRoomAvailability(id uuid.UUID, available bool) error {
	logger.Info("UpdateRoomAvailability repository method called")
	return r.db.Model(&entity.MeetingRoom{}).Where("id = ?", id).Update("available", available).Error
}

func (r *meetingRoomRepo) DeleteMeetingRoom(id uuid.UUID) error {
	logger.Info("DeleteMeetingRoom repository method called")
	return r.db.Delete(&entity.MeetingRoom{}, "id = ?", id).Error
}
