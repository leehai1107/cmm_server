package entity

import (
    "time"

    "github.com/google/uuid"
)

type Voucher struct {
    ID              uuid.UUID `gorm:"primaryKey;column:id"`
    Code            string    `gorm:"column:code;unique;not null"`
    DiscountPercent int       `gorm:"column:discount_percent;not null"`
    MaxUses         int       `gorm:"column:max_uses"`
    UsedCount       int       `gorm:"column:used_count;default:0"`
    ServiceID       int       `gorm:"column:service_id"`
    ValidFrom       time.Time `gorm:"column:valid_from"`
    ValidTo         time.Time `gorm:"column:valid_to"`
}