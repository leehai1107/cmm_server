package entity

import (
    "time"

    "github.com/google/uuid"
)

type CoffeeShop struct {
    ID          uuid.UUID `gorm:"primaryKey;column:id"`
    OwnerID     uuid.UUID `gorm:"column:owner_id;not null"`
    Name        string    `gorm:"column:name;not null"`
    Location    string    `gorm:"column:location;not null"`
    Description string    `gorm:"column:description"`
    CreatedAt   time.Time `gorm:"column:created_at;default:now()"`
}

type CommissionRate struct {
    ID           int       `gorm:"primaryKey;column:id"`
    CoffeeShopID uuid.UUID `gorm:"column:coffee_shop_id;unique;not null"`
    RatePercent  int       `gorm:"column:rate_percent;not null"`
}