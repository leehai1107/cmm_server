package entity

type Service struct {
    ID          int    `gorm:"primaryKey;column:id"`
    Name        string `gorm:"column:name;unique;not null"`
    Description string `gorm:"column:description"`
}