package models

import (
	"gorm.io/gorm"
	"time"
)

type Doctor struct {
	gorm.Model
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Gender          string    `json:"gender"`
	Birthdate       time.Time `json:"birthdate"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	ExperienceYears int       `json:"experience_years"`
	Specialization  string    `json:"specialization"`
	PhotoUrl        string    `json:"photo_url"`
}
