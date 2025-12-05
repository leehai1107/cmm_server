package entity

import (
    "time"

    "github.com/google/uuid"
)

type ShopPost struct {
    ID           uuid.UUID `gorm:"primaryKey;column:id"`
    CoffeeShopID uuid.UUID `gorm:"column:coffee_shop_id;not null"`
    Title        string    `gorm:"column:title;not null"`
    Content      string    `gorm:"column:content;not null"`
    PublishedAt  time.Time `gorm:"column:published_at;default:now()"`
}

type InternalPost struct {
    ID          uuid.UUID `gorm:"primaryKey;column:id"`
    Title       string    `gorm:"column:title;not null"`
    Content     string    `gorm:"column:content;not null"`
    CreatedBy   uuid.UUID `gorm:"column:created_by;not null"`
    VisibleTo   string    `gorm:"column:visible_to;not null"`
    PublishedAt time.Time `gorm:"column:published_at;default:now()"`
}