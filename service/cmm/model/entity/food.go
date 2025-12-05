package entity

import (
    "time"

    "github.com/google/uuid"
)

type FoodItem struct {
    ID           uuid.UUID `gorm:"primaryKey;column:id"`
    Name         string    `gorm:"column:name;not null"`
    Price        float64   `gorm:"column:price;not null"`
    CoffeeShopID uuid.UUID `gorm:"column:coffee_shop_id;not null"`
    Available    bool      `gorm:"column:available;default:true"`
}

type FoodOrder struct {
    ID           uuid.UUID `gorm:"primaryKey;column:id"`
    CustomerID   uuid.UUID `gorm:"column:customer_id;not null"`
    CoffeeShopID uuid.UUID `gorm:"column:coffee_shop_id;not null"`
    TotalPrice   float64   `gorm:"column:total_price;not null"`
    VoucherID    uuid.UUID `gorm:"column:voucher_id"`
    Status       string    `gorm:"column:status;default:ordered"`
    CreatedAt    time.Time `gorm:"column:created_at;default:now()"`
}

type FoodOrderItem struct {
    ID          int       `gorm:"primaryKey;column:id"`
    FoodOrderID uuid.UUID `gorm:"column:food_order_id;not null"`
    FoodItemID  uuid.UUID `gorm:"column:food_item_id;not null"`
    Quantity    int       `gorm:"column:quantity;not null"`
    UnitPrice   float64   `gorm:"column:unit_price;not null"`
}