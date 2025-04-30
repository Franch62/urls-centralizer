package models

import "gorm.io/gorm"

type URL struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	Source string `json:"source" binding:"required"`
	URL    string `json:"url" binding:"required" gorm:"uniqueIndex"`
	gorm.Model
}
