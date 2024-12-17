package service

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
	"github.com/joremysh/fliqt/internal/repository"
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
	repo repository.Employee
}

func NewEmployeeService(repo repository.Employee) EmployeeService {
	return &employeeService{
		repo: repo,
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
	employee, err := e.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("employee not found by id: %d", id)
	}
	if err != nil {
		return nil, err
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

	return e.repo.GetByID(existed.ID)
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
