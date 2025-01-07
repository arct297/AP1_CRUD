package models

import (
	"gorm.io/gorm"
)

type Patient struct {
	gorm.Model
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Gender  string `json:"gender"`
	Contact string `json:"contact"`
	Address string `json:"address"`
}
