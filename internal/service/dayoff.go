package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
	"github.com/joremysh/fliqt/internal/repository"
)

type DayOffService interface {
	SubmitDayOff(ctx context.Context, record *model.DayOffRecord) error
	ListDayOffs(ctx context.Context, employeeID uint, params *model.ListParams) (*PaginatedResult[model.DayOffRecord], error)
	CancelDayOff(ctx context.Context, id uint, cancellationReason string) error
}

type dayOffService struct {
	repo         repository.DayOff
	employeeRepo repository.Employee
}

func NewDayOffService(repo repository.DayOff, employeeRepo repository.Employee) DayOffService {
	return &dayOffService{
		repo:         repo,
		employeeRepo: employeeRepo,
	}
}

func (s *dayOffService) SubmitDayOff(ctx context.Context, record *model.DayOffRecord) error {
	_, err := s.employeeRepo.GetByID(record.EmployeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("employee not found by id: %d", record.EmployeeID)
		}
		return err
	}

	if err := s.validateDayOff(record); err != nil {
		return err
	}

	exists, err := s.repo.ExistsOverlapping(record.EmployeeID, record.StartTime, record.EndTime)
	if err != nil {
		return err
	}
	if exists {
		return ErrOverlappingDayOff
	}

	return s.repo.Create(record)
}

func (s *dayOffService) ListDayOffs(ctx context.Context, employeeID uint, params *model.ListParams) (*PaginatedResult[model.DayOffRecord], error) {
	records, total, err := s.repo.List(params)
	if err != nil {
		return nil, err
	}

	return &PaginatedResult[model.DayOffRecord]{
		Data:       records,
		TotalCount: total,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (s *dayOffService) CancelDayOff(ctx context.Context, id uint, cancellationReason string) error {
	record, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDayOffNotFound
		}
		return err
	}

	record.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	record.Reason = fmt.Sprintf("%s (Cancelled: %s)", record.Reason, cancellationReason)

	return s.repo.Update(record)
}

// Custom errors
var (
	ErrEmployeeNotFound     = errors.New("employee not found")
	ErrOverlappingDayOff    = errors.New("overlapping day off exists")
	ErrInvalidDayOffType    = errors.New("invalid day off type")
	ErrInvalidDateRange     = errors.New("end date must be after start date")
	ErrPastDateNotAllowed   = errors.New("cannot submit day off for past dates")
	ErrReasonRequired       = errors.New("reason is required")
	ErrDayOffNotFound       = errors.New("day off record not found")
	ErrCantCancelPastDayOff = errors.New("cannot cancel past day off")
)

func (s *dayOffService) validateDayOff(record *model.DayOffRecord) error {
	validTypes := map[string]bool{
		"PTO":            true,
		"sick leave":     true,
		"parental leave": true,
		"bereavement":    true,
	}
	if !validTypes[record.DayOffType] {
		return ErrInvalidDayOffType
	}

	// Validate dates
	if record.StartTime.After(record.EndTime) {
		return ErrInvalidDateRange
	}

	// Validate reason provided
	if strings.TrimSpace(record.Reason) == "" {
		return ErrReasonRequired
	}

	return nil
}
