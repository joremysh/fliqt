package repository

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
)

// repositories/day_off.go

type DayOff interface {
	Create(record *model.DayOffRecord) error
	GetByID(id uint) (*model.DayOffRecord, error)
	Update(record *model.DayOffRecord) error
	List(params *model.ListParams) ([]model.DayOffRecord, int64, error)
	ExistsOverlapping(employeeID uint, startTime, endTime time.Time) (bool, error)
}

type dayOffRepo struct {
	db *gorm.DB
}

func NewDayOffRepo(db *gorm.DB) DayOff {
	return &dayOffRepo{
		db: db,
	}
}

func (r *dayOffRepo) Create(record *model.DayOffRecord) error {
	return r.db.Create(record).Error
}

func (r *dayOffRepo) GetByID(id uint) (*model.DayOffRecord, error) {
	var record model.DayOffRecord
	err := r.db.Preload("Employee").First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *dayOffRepo) Update(record *model.DayOffRecord) error {
	return r.db.Save(record).Error
}

func (r *dayOffRepo) List(params *model.ListParams) ([]model.DayOffRecord, int64, error) {
	var records []model.DayOffRecord
	var totalCount int64
	query := r.db.Model(&model.DayOffRecord{})
	countQuery := r.db.Model(&model.DayOffRecord{})

	var listFilterColumnNames = map[string]string{"DayOffType": "day_off_type"}
	// Apply filters
	for _, field := range listFilterColumnNames {
		if s, ok := params.Filters[field]; ok {
			condition := listFilterColumnNames[field] + " like ?"
			value := fmt.Sprintf("%s", s)
			query = query.Where(condition, value)
			countQuery = countQuery.Where(condition, value)
		}
	}

	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	} else {
		// Default sorting by start time descending
		query = query.Order("start_time DESC")
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Preload("Employee").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, totalCount, nil
}

func (r *dayOffRepo) ExistsOverlapping(employeeID uint, startTime, endTime time.Time) (bool, error) {
	var count int64

	err := r.db.Model(&model.DayOffRecord{}).
		Where("employee_id = ?", employeeID).
		Where("deleted_at IS NULL"). // Exclude cancelled records
		Where(
			"(start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?) OR (start_time <= ? AND end_time >= ?)",
			startTime, endTime,
			startTime, endTime,
			startTime, endTime,
		).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
