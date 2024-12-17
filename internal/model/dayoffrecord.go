package model

import (
	"time"

	"gorm.io/gorm"
)

type DayOffRecord struct {
	gorm.Model
	EmployeeID uint
	Employee   Employee `gorm:"foreignKey:EmployeeID"`
	DayOffType string   `gorm:"type:varchar(50)"`
	Reason     string
	StartTime  time.Time `gorm:"index"`
	EndTime    time.Time
}
