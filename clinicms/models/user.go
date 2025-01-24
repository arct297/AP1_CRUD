package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID         uint           `gorm:"primaryKey"`
	Login      string         `gorm:"unique;not null" json:"login"`
	Password   string         `json:"password"`
	IsVerified bool           `json:"is_verified" gorm:"default:false"`
	Role       string         `json:"role"` // "user", "admin"
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
