package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/leehai1107/cmm_server/pkg/logger"
	"github.com/leehai1107/cmm_server/service/cmm/model/entity"
	"gorm.io/gorm"
)

type IBookingRepo interface {
	CreateBooking(booking *entity.Booking) error
	GetBookingByID(id uuid.UUID) (*entity.Booking, error)
	GetBookingsByCustomer(customerID uuid.UUID) ([]entity.Booking, error)
	GetBookingsByMeetingRoom(roomID uuid.UUID) ([]entity.Booking, error)
	UpdateBookingStatus(id uuid.UUID, status string) error
	CheckRoomAvailability(roomID uuid.UUID, startTime, endTime time.Time) (bool, error)
	CancelBooking(id uuid.UUID) error
}

type bookingRepo struct {
	db *gorm.DB
}

func NewBookingRepo(db *gorm.DB) IBookingRepo {
	return &bookingRepo{
		db: db,
	}
}

func (r *bookingRepo) CreateBooking(booking *entity.Booking) error {
	logger.Info("CreateBooking repository method called")
	return r.db.Create(booking).Error
}

func (r *bookingRepo) GetBookingByID(id uuid.UUID) (*entity.Booking, error) {
	logger.Info("GetBookingByID repository method called")
	var booking entity.Booking
	err := r.db.Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepo) GetBookingsByCustomer(customerID uuid.UUID) ([]entity.Booking, error) {
	logger.Info("GetBookingsByCustomer repository method called")
	var bookings []entity.Booking
	err := r.db.Where("customer_id = ?", customerID).Order("created_at DESC").Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepo) GetBookingsByMeetingRoom(roomID uuid.UUID) ([]entity.Booking, error) {
	logger.Info("GetBookingsByMeetingRoom repository method called")
	var bookings []entity.Booking
	err := r.db.Where("meeting_room_id = ?", roomID).Order("start_time ASC").Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepo) UpdateBookingStatus(id uuid.UUID, status string) error {
	logger.Info("UpdateBookingStatus repository method called")
	return r.db.Model(&entity.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *bookingRepo) CheckRoomAvailability(roomID uuid.UUID, startTime, endTime time.Time) (bool, error) {
	logger.Info("CheckRoomAvailability repository method called")
	var count int64

	err := r.db.Model(&entity.Booking{}).
		Where("meeting_room_id = ?", roomID).
		Where("status != ?", "cancelled").
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?)",
			endTime, startTime, // Overlaps at the start
			startTime, endTime, // Overlaps at the end
			startTime, endTime, // Completely within
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *bookingRepo) CancelBooking(id uuid.UUID) error {
	logger.Info("CancelBooking repository method called")
	return r.db.Model(&entity.Booking{}).Where("id = ?", id).Update("status", "cancelled").Error
}
