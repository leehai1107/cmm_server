package entity

import (
    "time"

    "github.com/google/uuid"
)

type MeetingRoom struct {
    ID           uuid.UUID `gorm:"primaryKey;column:id"`
    CoffeeShopID uuid.UUID `gorm:"column:coffee_shop_id;not null"`
    Name         string    `gorm:"column:name;not null"`
    Capacity     int       `gorm:"column:capacity;not null"`
    PricePerHour float64   `gorm:"column:price_per_hour;not null"`
    Available    bool      `gorm:"column:available;default:true"`
}

type Booking struct {
    ID            uuid.UUID `gorm:"primaryKey;column:id"`
    CustomerID    uuid.UUID `gorm:"column:customer_id;not null"`
    MeetingRoomID uuid.UUID `gorm:"column:meeting_room_id;not null"`
    StartTime     time.Time `gorm:"column:start_time;not null"`
    EndTime       time.Time `gorm:"column:end_time;not null"`
    TotalPrice    float64   `gorm:"column:total_price;not null"`
    VoucherID     uuid.UUID `gorm:"column:voucher_id"`
    Status        string    `gorm:"column:status;default:booked"`
    CreatedAt     time.Time `gorm:"column:created_at;default:now()"`
}