package model

import (
	"time"

	"gorm.io/gorm"
)

type Schools struct {
	ID         string `gorm:"primaryKey"`
	SchoolYear string
	Name       string
	Logo       string
	Address    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
