package entity

import (
    "time"

    "github.com/google/uuid"
)

type Transaction struct {
    ID              uuid.UUID `gorm:"primaryKey;column:id"`
    UserID          uuid.UUID `gorm:"column:user_id;not null"`
    ServiceID       int       `gorm:"column:service_id;not null"`
    ServiceRefID    uuid.UUID `gorm:"column:service_ref_id;not null"`
    Amount          float64   `gorm:"column:amount;not null"`
    PaymentMethodID int       `gorm:"column:payment_method_id"`
    PaidAt          time.Time `gorm:"column:paid_at;default:now()"`
    Status          string    `gorm:"column:status;default:completed"`
}