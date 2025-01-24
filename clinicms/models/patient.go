package models

import (
	"gorm.io/gorm"
	"time"
)

type Patient struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null;unique" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	Name      string         `json:"name"`
	Age       int            `json:"age"`
	Gender    string         `json:"gender"`
	Contact   string         `json:"contact"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
