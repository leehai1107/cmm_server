package request

import (
	"time"

	"github.com/google/uuid"
)

// Booking requests
type CreateBooking struct {
	MeetingRoomID uuid.UUID `json:"meeting_room_id" binding:"required"`
	StartTime     time.Time `json:"start_time" binding:"required"`
	EndTime       time.Time `json:"end_time" binding:"required"`
	VoucherCode   string    `json:"voucher_code,omitempty"`
}

type CancelBooking struct {
	BookingID uuid.UUID `json:"booking_id" binding:"required"`
}

// Coffee Shop requests
type CreateCoffeeShop struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Description string `json:"description"`
}

type UpdateCoffeeShop struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

type SetCommissionRate struct {
	CoffeeShopID uuid.UUID `json:"coffee_shop_id" binding:"required"`
	RatePercent  int       `json:"rate_percent" binding:"required,min=0,max=100"`
}

// Meeting Room requests
type CreateMeetingRoom struct {
	CoffeeShopID uuid.UUID `json:"coffee_shop_id" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Capacity     int       `json:"capacity" binding:"required,min=1"`
	PricePerHour float64   `json:"price_per_hour" binding:"required,min=0"`
}

type UpdateMeetingRoom struct {
	ID           uuid.UUID `json:"id" binding:"required"`
	Name         string    `json:"name"`
	Capacity     int       `json:"capacity" binding:"min=1"`
	PricePerHour float64   `json:"price_per_hour" binding:"min=0"`
	Available    *bool     `json:"available"`
}

// Wallet requests
type TopupWallet struct {
	Amount float64 `json:"amount" binding:"required,min=1"`
	Method string  `json:"method" binding:"required"`
}

type ConfirmTopup struct {
	TopupID uuid.UUID `json:"topup_id" binding:"required"`
}

// Voucher requests
type CreateVoucher struct {
	Code            string     `json:"code" binding:"required"`
	DiscountPercent int        `json:"discount_percent" binding:"required,min=1,max=100"`
	MaxUses         int        `json:"max_uses" binding:"min=0"`
	ServiceID       int        `json:"service_id"`
	ValidFrom       time.Time  `json:"valid_from" binding:"required"`
	ValidTo         time.Time  `json:"valid_to" binding:"required"`
}

type UpdateVoucher struct {
	ID              uuid.UUID  `json:"id" binding:"required"`
	Code            string     `json:"code"`
	DiscountPercent int        `json:"discount_percent" binding:"min=1,max=100"`
	MaxUses         int        `json:"max_uses" binding:"min=0"`
	ValidFrom       *time.Time `json:"valid_from"`
	ValidTo         *time.Time `json:"valid_to"`
}

type ApplyVoucher struct {
	VoucherCode string  `json:"voucher_code" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,min=0"`
}

// Post requests
type CreateShopPost struct {
	CoffeeShopID uuid.UUID `json:"coffee_shop_id" binding:"required"`
	Title        string    `json:"title" binding:"required"`
	Content      string    `json:"content" binding:"required"`
}

type UpdateShopPost struct {
	ID      uuid.UUID `json:"id" binding:"required"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
}

type CreateInternalPost struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	VisibleTo string `json:"visible_to" binding:"required"` // "all", "admin", "owner", "customer"
}

type UpdateInternalPost struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	VisibleTo string    `json:"visible_to"`
}
