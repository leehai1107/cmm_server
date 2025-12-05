package response

import (
	"time"

	"github.com/google/uuid"
)

// Booking responses
type BookingResponse struct {
	ID            uuid.UUID `json:"id"`
	CustomerID    uuid.UUID `json:"customer_id"`
	MeetingRoomID uuid.UUID `json:"meeting_room_id"`
	RoomName      string    `json:"room_name,omitempty"`
	ShopName      string    `json:"shop_name,omitempty"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	TotalPrice    float64   `json:"total_price"`
	VoucherID     uuid.UUID `json:"voucher_id,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// Coffee Shop responses
type CoffeeShopResponse struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CommissionRateResponse struct {
	ID           int       `json:"id"`
	CoffeeShopID uuid.UUID `json:"coffee_shop_id"`
	RatePercent  int       `json:"rate_percent"`
}

// Meeting Room responses
type MeetingRoomResponse struct {
	ID           uuid.UUID `json:"id"`
	CoffeeShopID uuid.UUID `json:"coffee_shop_id"`
	ShopName     string    `json:"shop_name,omitempty"`
	Name         string    `json:"name"`
	Capacity     int       `json:"capacity"`
	PricePerHour float64   `json:"price_per_hour"`
	Available    bool      `json:"available"`
}

// Wallet responses
type WalletResponse struct {
	UserID  uuid.UUID `json:"user_id"`
	Balance float64   `json:"balance"`
}

type TopupResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Voucher responses
type VoucherResponse struct {
	ID              uuid.UUID `json:"id"`
	Code            string    `json:"code"`
	DiscountPercent int       `json:"discount_percent"`
	MaxUses         int       `json:"max_uses"`
	UsedCount       int       `json:"used_count"`
	ServiceID       int       `json:"service_id,omitempty"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidTo         time.Time `json:"valid_to"`
}

type VoucherCalculation struct {
	OriginalAmount  float64 `json:"original_amount"`
	DiscountAmount  float64 `json:"discount_amount"`
	FinalAmount     float64 `json:"final_amount"`
	DiscountPercent int     `json:"discount_percent"`
	VoucherCode     string  `json:"voucher_code"`
}

// Post responses
type ShopPostResponse struct {
	ID           uuid.UUID `json:"id"`
	CoffeeShopID uuid.UUID `json:"coffee_shop_id"`
	ShopName     string    `json:"shop_name,omitempty"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	PublishedAt  time.Time `json:"published_at"`
}

type InternalPostResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatorName string    `json:"creator_name,omitempty"`
	VisibleTo   string    `json:"visible_to"`
	PublishedAt time.Time `json:"published_at"`
}

// Transaction responses
type TransactionResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	ServiceID       int       `json:"service_id"`
	ServiceRefID    uuid.UUID `json:"service_ref_id"`
	Amount          float64   `json:"amount"`
	PaymentMethodID int       `json:"payment_method_id,omitempty"`
	PaidAt          time.Time `json:"paid_at"`
	Status          string    `json:"status"`
}
