package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name        string    `gorm:"type:varchar(50);not null;index"`
	Email       string    `gorm:"type:varchar(100);not null;uniqueIndex"`
	PhoneNumber string    `gorm:"type:varchar(20);not null"`
	Department  string    `gorm:"type:varchar(50);not null"`
	Address     string    `gorm:"type:varchar(255);not null"`
	Salary      int       `gorm:"type:mediumint unsigned;not null"` // Assuming NTD is used here, if decimal points need to be stored, it can be switched to `decimal` or other methods.
	OnboardDate time.Time `gorm:"not null"`
}
