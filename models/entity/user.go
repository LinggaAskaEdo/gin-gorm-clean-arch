package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	Name     string    `gorm:"NOT NULL"`
	Email    string    `gorm:"NOT NULL"`
	Password string    `gorm:"NOT NULL"`
	Age      uint8     `gorm:"NOT NULL"`
	Birthday time.Time `gorm:"NOT NULL"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return "user"
}
