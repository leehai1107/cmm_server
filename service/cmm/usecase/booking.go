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

type IBookingUsecase interface {
	CreateBooking(ctx context.Context, customerID uuid.UUID, req request.CreateBooking) (*response.BookingResponse, error)
	GetBooking(ctx context.Context, bookingID uuid.UUID) (*response.BookingResponse, error)
	GetCustomerBookings(ctx context.Context, customerID uuid.UUID) ([]response.BookingResponse, error)
	CancelBooking(ctx context.Context, customerID uuid.UUID, bookingID uuid.UUID) error
	GetRoomBookings(ctx context.Context, roomID uuid.UUID) ([]response.BookingResponse, error)
}

type bookingUsecase struct {
	bookingRepo     repository.IBookingRepo
	meetingRoomRepo repository.IMeetingRoomRepo
	walletRepo      repository.IWalletRepo
	voucherRepo     repository.IVoucherRepo
	transactionRepo repository.ITransactionRepo
}

func NewBookingUsecase(
	bookingRepo repository.IBookingRepo,
	meetingRoomRepo repository.IMeetingRoomRepo,
	walletRepo repository.IWalletRepo,
	voucherRepo repository.IVoucherRepo,
	transactionRepo repository.ITransactionRepo,
) IBookingUsecase {
	return &bookingUsecase{
		bookingRepo:     bookingRepo,
		meetingRoomRepo: meetingRoomRepo,
		walletRepo:      walletRepo,
		voucherRepo:     voucherRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, customerID uuid.UUID, req request.CreateBooking) (*response.BookingResponse, error) {
	log := logger.EnhanceWith(ctx)
	log.Info("CreateBooking usecase called")

	// Validate time
	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return nil, errors.New("end time must be after start time")
	}

	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("cannot book in the past")
	}

	// Get meeting room
	room, err := u.meetingRoomRepo.GetMeetingRoomByID(req.MeetingRoomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("meeting room not found")
		}
		return nil, err
	}

	if !room.Available {
		return nil, errors.New("meeting room is not available")
	}

	// Check availability
	available, err := u.bookingRepo.CheckRoomAvailability(req.MeetingRoomID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("meeting room is already booked for this time slot")
	}

	// Calculate price
	duration := req.EndTime.Sub(req.StartTime).Hours()
	totalPrice := room.PricePerHour * duration

	// Apply voucher if provided
	var voucherID *uuid.UUID
	if req.VoucherCode != "" {
		voucher, err := u.voucherRepo.GetVoucherByCode(req.VoucherCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("invalid voucher code")
			}
			return nil, err
		}

		// Validate voucher
		now := time.Now()
		if now.Before(voucher.ValidFrom) || now.After(voucher.ValidTo) {
			return nil, errors.New("voucher is not valid at this time")
		}
		if voucher.MaxUses > 0 && voucher.UsedCount >= voucher.MaxUses {
			return nil, errors.New("voucher has reached maximum uses")
		}

		// Apply discount
		discount := totalPrice * float64(voucher.DiscountPercent) / 100
		totalPrice -= discount
		voucherID = &voucher.ID

		// Increment voucher usage
		if err := u.voucherRepo.IncrementUsedCount(voucher.ID); err != nil {
			log.Errorw("Failed to increment voucher usage", "error", err)
		}
	}

	// Check wallet balance
	wallet, err := u.walletRepo.GetWalletByUserID(customerID)
	if err != nil {
		return nil, errors.New("wallet not found")
	}
	if wallet.Balance < totalPrice {
		return nil, errors.New("insufficient balance")
	}

	// Create booking
	booking := &entity.Booking{
		ID:            uuid.New(),
		CustomerID:    customerID,
		MeetingRoomID: req.MeetingRoomID,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		TotalPrice:    totalPrice,
		Status:        "booked",
		CreatedAt:     time.Now(),
	}
	if voucherID != nil {
		booking.VoucherID = *voucherID
	}

	if err := u.bookingRepo.CreateBooking(booking); err != nil {
		return nil, err
	}

	// Deduct from wallet
	if err := u.walletRepo.DeductBalance(customerID, totalPrice); err != nil {
		log.Errorw("Failed to deduct balance", "error", err)
		return nil, errors.New("payment failed")
	}

	// Create transaction record
	transaction := &entity.Transaction{
		ID:           uuid.New(),
		UserID:       customerID,
		ServiceID:    1, // 1 for booking
		ServiceRefID: booking.ID,
		Amount:       totalPrice,
		PaidAt:       time.Now(),
		Status:       "completed",
	}
	if err := u.transactionRepo.CreateTransaction(transaction); err != nil {
		log.Errorw("Failed to create transaction", "error", err)
	}

	return &response.BookingResponse{
		ID:            booking.ID,
		CustomerID:    booking.CustomerID,
		MeetingRoomID: booking.MeetingRoomID,
		RoomName:      room.Name,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		TotalPrice:    booking.TotalPrice,
		VoucherID:     booking.VoucherID,
		Status:        booking.Status,
		CreatedAt:     booking.CreatedAt,
	}, nil
}

func (u *bookingUsecase) GetBooking(ctx context.Context, bookingID uuid.UUID) (*response.BookingResponse, error) {
	booking, err := u.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}

	room, _ := u.meetingRoomRepo.GetMeetingRoomByID(booking.MeetingRoomID)
	roomName := ""
	if room != nil {
		roomName = room.Name
	}

	return &response.BookingResponse{
		ID:            booking.ID,
		CustomerID:    booking.CustomerID,
		MeetingRoomID: booking.MeetingRoomID,
		RoomName:      roomName,
		StartTime:     booking.StartTime,
		EndTime:       booking.EndTime,
		TotalPrice:    booking.TotalPrice,
		VoucherID:     booking.VoucherID,
		Status:        booking.Status,
		CreatedAt:     booking.CreatedAt,
	}, nil
}

func (u *bookingUsecase) GetCustomerBookings(ctx context.Context, customerID uuid.UUID) ([]response.BookingResponse, error) {
	bookings, err := u.bookingRepo.GetBookingsByCustomer(customerID)
	if err != nil {
		return nil, err
	}

	var result []response.BookingResponse
	for _, booking := range bookings {
		room, _ := u.meetingRoomRepo.GetMeetingRoomByID(booking.MeetingRoomID)
		roomName := ""
		if room != nil {
			roomName = room.Name
		}

		result = append(result, response.BookingResponse{
			ID:            booking.ID,
			CustomerID:    booking.CustomerID,
			MeetingRoomID: booking.MeetingRoomID,
			RoomName:      roomName,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
			TotalPrice:    booking.TotalPrice,
			VoucherID:     booking.VoucherID,
			Status:        booking.Status,
			CreatedAt:     booking.CreatedAt,
		})
	}

	return result, nil
}

func (u *bookingUsecase) CancelBooking(ctx context.Context, customerID uuid.UUID, bookingID uuid.UUID) error {
	log := logger.EnhanceWith(ctx)

	booking, err := u.bookingRepo.GetBookingByID(bookingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("booking not found")
		}
		return err
	}

	if booking.CustomerID != customerID {
		return errors.New("unauthorized to cancel this booking")
	}

	if booking.Status == "cancelled" {
		return errors.New("booking is already cancelled")
	}

	// Only allow cancellation if booking is at least 24 hours away
	if time.Until(booking.StartTime) < 24*time.Hour {
		return errors.New("cannot cancel booking less than 24 hours before start time")
	}

	// Cancel booking
	if err := u.bookingRepo.CancelBooking(bookingID); err != nil {
		return err
	}

	// Refund to wallet
	if err := u.walletRepo.UpdateBalance(customerID, booking.TotalPrice); err != nil {
		log.Errorw("Failed to refund balance", "error", err)
		return errors.New("failed to process refund")
	}

	// Create refund transaction
	transaction := &entity.Transaction{
		ID:           uuid.New(),
		UserID:       customerID,
		ServiceID:    1,
		ServiceRefID: booking.ID,
		Amount:       -booking.TotalPrice, // Negative for refund
		PaidAt:       time.Now(),
		Status:       "refunded",
	}
	if err := u.transactionRepo.CreateTransaction(transaction); err != nil {
		log.Errorw("Failed to create refund transaction", "error", err)
	}

	return nil
}

func (u *bookingUsecase) GetRoomBookings(ctx context.Context, roomID uuid.UUID) ([]response.BookingResponse, error) {
	bookings, err := u.bookingRepo.GetBookingsByMeetingRoom(roomID)
	if err != nil {
		return nil, err
	}

	var result []response.BookingResponse
	for _, booking := range bookings {
		result = append(result, response.BookingResponse{
			ID:            booking.ID,
			CustomerID:    booking.CustomerID,
			MeetingRoomID: booking.MeetingRoomID,
			StartTime:     booking.StartTime,
			EndTime:       booking.EndTime,
			TotalPrice:    booking.TotalPrice,
			VoucherID:     booking.VoucherID,
			Status:        booking.Status,
			CreatedAt:     booking.CreatedAt,
		})
	}

	return result, nil
}
