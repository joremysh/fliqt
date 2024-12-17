package repository

import (
	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
)

type Seed struct {
	Name string
	Run  func(*gorm.DB) error
}

func All() []Seed {
	seeds := make([]Seed, 10)
	for i := 0; i < len(seeds); i++ {
		employee := mockEmployee()
		seeds[i].Name = employee.Name
		seeds[i].Run = func(gdb *gorm.DB) error {
			repo := NewEmployeeRepo(gdb)
			return repo.Create(employee)
		}
	}

	return seeds
}

var departments = []string{"Sales", "Financial", "Design", "Engineering", "General affairs"}

func mockEmployee() *model.Employee {
	return &model.Employee{
		Name:        gofakeit.Name(),
		Email:       gofakeit.Email(),
		PhoneNumber: gofakeit.Phone(),
		Department:  departments[gofakeit.IntRange(0, len(departments)-1)],
		Address:     gofakeit.Address().Address,
		Salary:      gofakeit.IntRange(50000, 200000),
		OnboardDate: gofakeit.Date(),
	}
}
