package models

import "gorm.io/gorm"

type GroupsMembers struct {
	gorm.Model
	GroupID uint `gorm:"not null" json:"group_id"`
	UserID  uint `gorm:"not null" json:"user_id"`
}
