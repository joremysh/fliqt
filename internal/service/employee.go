package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
	"github.com/joremysh/fliqt/internal/repository"
	"github.com/joremysh/fliqt/pkg/cache"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee *model.Employee) (*model.Employee, error)
	GetEmployee(ctx context.Context, id uint) (*model.Employee, error)
	DeleteEmployee(ctx context.Context, id uint) error
	UpdateEmployee(ctx context.Context, employee *model.Employee) (*model.Employee, error)
	ListEmployees(ctx context.Context, params *model.ListParams) (*PaginatedResult[model.Employee], error)
}

type PaginatedResult[T any] struct {
	Data       []T
	TotalCount int64
	Page       int
	PageSize   int
}

type employeeService struct {
	repo        repository.Employee
	redisClient *cache.RedisClient
}

func NewEmployeeService(repo repository.Employee, redisClient *cache.RedisClient) EmployeeService {
	return &employeeService{
		repo:        repo,
		redisClient: redisClient,
	}
}

func (e employeeService) CreateEmployee(ctx context.Context, employee *model.Employee) (*model.Employee, error) {
	_, err := e.repo.GetByEmail(employee.Email)
	if err == nil {
		return nil, fmt.Errorf("employee email already exists: %s", employee.Email)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err = e.repo.Create(employee); err != nil {
		return nil, err
	}

	created, err := e.repo.GetByID(employee.ID)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (e employeeService) GetEmployee(ctx context.Context, id uint) (*model.Employee, error) {
	cacheKey := fmt.Sprintf("employee:%d", id)
	var employee *model.Employee
	err := e.redisClient.Get(ctx, cacheKey, employee)
	if err == nil {
		return employee, nil
	}

	employee, err = e.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("employee not found by id: %d", id)
	}
	if err != nil {
		return nil, err
	}

	err = e.redisClient.Set(ctx, cacheKey, employee, 1*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("Failed to cache employee: %s", err.Error())
	}
	return employee, nil
}

func (e employeeService) DeleteEmployee(ctx context.Context, id uint) error {
	return e.repo.Delete(id)
}

func (e employeeService) UpdateEmployee(ctx context.Context, employee *model.Employee) (*model.Employee, error) {
	existed, err := e.repo.GetByID(employee.ID)
	if err != nil {
		return nil, err
	}

	if existed.Email != employee.Email {
		if _, err := e.repo.GetByEmail(employee.Email); err == nil {
			return nil, fmt.Errorf("employee email already exists: %s", employee.Email)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	existed.Email = employee.Email
	existed.PhoneNumber = employee.PhoneNumber
	existed.Address = employee.Address
	existed.Department = employee.Department
	existed.Salary = employee.Salary

	if err := e.repo.Update(existed); err != nil {
		return nil, err
	}

	updated, err := e.repo.GetByID(existed.ID)
	if err != nil {
		return nil, err
	}
	cacheKey := fmt.Sprintf("employee:%d", employee.ID)
	_ = e.redisClient.Delete(ctx, cacheKey) // todo: make update in transaction and rollback if delete cache failed

	return updated, nil
}

func (e employeeService) ListEmployees(ctx context.Context, params *model.ListParams) (*PaginatedResult[model.Employee], error) {
	results, totalCount, err := e.repo.List(params)
	if err != nil {
		return nil, err
	}
	return &PaginatedResult[model.Employee]{
		Data:       results,
		TotalCount: totalCount,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}
