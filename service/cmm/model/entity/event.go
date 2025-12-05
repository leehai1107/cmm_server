package entity

import (
    "time"

    "github.com/google/uuid"
)

type Event struct {
    ID           uuid.UUID `gorm:"primaryKey;column:id"`
    CoffeeShopID uuid.UUID `gorm:"column:coffee_shop_id;not null"`
    Name         string    `gorm:"column:name;not null"`
    Description  string    `gorm:"column:description"`
    StartTime    time.Time `gorm:"column:start_time;not null"`
    EndTime      time.Time `gorm:"column:end_time;not null"`
    Price        float64   `gorm:"column:price;not null"`
}

type EventTicket struct {
    ID         uuid.UUID `gorm:"primaryKey;column:id"`
    CustomerID uuid.UUID `gorm:"column:customer_id;not null"`
    EventID    uuid.UUID `gorm:"column:event_id;not null"`
    VoucherID  uuid.UUID `gorm:"column:voucher_id"`
    PricePaid  float64   `gorm:"column:price_paid;not null"`
    Status     string    `gorm:"column:status;default:booked"`
    BookedAt   time.Time `gorm:"column:booked_at;default:now()"`
}