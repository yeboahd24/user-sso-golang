package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email       string `gorm:"uniqueIndex;not null"`
	Password    string `gorm:"not null"`
	SSOProvider string
	SSOEmail    string
	SSOLinkedAt *time.Time
}
