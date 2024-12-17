package model

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);index"`
	Email       string `gorm:"type:varchar(100);uniqueIndex"`
	PhoneNumber string `gorm:"type:varchar(20)"`
	Department  string `gorm:"type:varchar(50)"`
	Address     string
	Salary      int
	OnboardDate time.Time
}
