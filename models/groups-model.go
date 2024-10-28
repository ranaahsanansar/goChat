package models

import "gorm.io/gorm"

type Groups struct {
	gorm.Model
	Name      string `gorm:"unique" json:"name"`
	CreatedBy string `json:"created_by"`
}
