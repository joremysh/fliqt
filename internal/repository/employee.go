package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
)

type Employee interface {
	Create(employee *model.Employee) error
	GetByID(id uint) (*model.Employee, error)
	GetByEmail(email string) (*model.Employee, error)
	Update(employee *model.Employee) error
	Delete(id uint) error
	List(params *model.ListParams) ([]model.Employee, int64, error)
}

type employeeRepo struct {
	gdb *gorm.DB
}

func (r *employeeRepo) GetByEmail(email string) (*model.Employee, error) {
	employee := &model.Employee{}
	err := r.gdb.First(employee, &model.Employee{Email: email}).Error
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func NewEmployeeRepo(gdb *gorm.DB) Employee {
	return &employeeRepo{gdb: gdb}
}

func (r *employeeRepo) Create(employee *model.Employee) error {
	return r.gdb.Create(employee).Error
}

func (r *employeeRepo) GetByID(id uint) (*model.Employee, error) {
	var employee model.Employee
	err := r.gdb.First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *employeeRepo) Update(employee *model.Employee) error {
	return r.gdb.Save(employee).Error
}

func (r *employeeRepo) Delete(id uint) error {
	return r.gdb.Delete(&model.Employee{}, id).Error
}

func (r *employeeRepo) List(params *model.ListParams) ([]model.Employee, int64, error) {
	query := r.gdb
	countQuery := query

	var listFilterColumnNames = []string{"name", "email", "department"}
	// Apply filters
	for _, field := range listFilterColumnNames {
		if s, ok := params.Filters[field]; ok {
			condition := field + " like ?"
			value := fmt.Sprintf("%s", s)
			query = query.Where(condition, value)
			countQuery = countQuery.Where(condition, value)
		}
	}

	var totalCount int64
	if err := countQuery.Model(&model.Employee{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if params.SortBy != "" {
		order := params.SortBy
		if params.SortOrder == "desc" {
			order += " DESC"
		}
		query = query.Order(order)
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	var employees []model.Employee
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, err
	}
	return employees, totalCount, nil
}
