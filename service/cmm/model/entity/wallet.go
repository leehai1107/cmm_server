package entity

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
    UserID  uuid.UUID `gorm:"primaryKey;column:user_id"`
    Balance float64   `gorm:"column:balance;not null;default:0"`
}

type Topup struct {
    ID        uuid.UUID `gorm:"primaryKey;column:id"`
    UserID    uuid.UUID `gorm:"column:user_id;not null"`
    Amount    float64   `gorm:"column:amount;not null"`
    Method    string    `gorm:"column:method;not null"`
    Status    string    `gorm:"column:status;default:pending"`
    CreatedAt time.Time `gorm:"column:created_at;default:now()"`
}