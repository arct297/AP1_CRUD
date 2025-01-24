package models

import (
	"gorm.io/gorm"
	"time"
)

type Doctor struct {
	ID              uint           `gorm:"primaryKey"`
	UserID          uint           `gorm:"not null;unique" json:"user_id"`
	User            User           `gorm:"foreignKey:UserID" json:"user"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Gender          string         `json:"gender"`
	Birthdate       time.Time      `json:"birthdate"`
	Email           string         `json:"email"`
	PhoneNumber     string         `json:"phone_number"`
	ExperienceYears int            `json:"experience_years"`
	Specialization  string         `json:"specialization"`
	PhotoUrl        string         `json:"photo_url"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
