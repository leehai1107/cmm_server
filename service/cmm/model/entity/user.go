package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleOwner    UserRole = "owner"
	RoleCustomer UserRole = "customer"
)

type User struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id"`
	FullName     string    `gorm:"column:full_name;not null"`
	Email        string    `gorm:"column:email;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	Role         UserRole  `gorm:"column:role;not null;default:customer"`
	CreatedAt    time.Time `gorm:"column:created_at;default:now()"`
}
