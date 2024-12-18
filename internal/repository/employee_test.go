package repository

import (
	"log"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
	testingx "github.com/joremysh/fliqt/pkg/testing"
)

var (
	gdb      *gorm.DB
	err      error
	pool     *dockertest.Pool
	resource *dockertest.Resource
)

func TestMain(m *testing.M) {
	pool, resource, gdb, err = testingx.NewMysqlInDocker()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = Migrate(gdb)
	if err != nil {
		log.Fatal(err.Error())
	}

	m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestEmployeeRepo_Create(t *testing.T) {
	tx := gdb.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	repo := NewEmployeeRepo(gdb)
	employee := MockEmployee()
	err = repo.Create(employee)
	require.NoError(t, err)
	require.NotNil(t, employee.ID)

	check := &model.Employee{}
	err = gdb.First(check, &model.Employee{Email: employee.Email}).Error
	require.NoError(t, err)
	require.Equal(t, employee.Name, check.Name)
	require.Equal(t, employee.PhoneNumber, check.PhoneNumber)
}

func TestEmployeeRepo_GetByID(t *testing.T) {
	tx := gdb.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	repo := NewEmployeeRepo(gdb)
	employee := MockEmployee()
	err = repo.Create(employee)
	require.NoError(t, err)
	require.NotNil(t, employee.ID)

	check, err := repo.GetByID(employee.ID)
	require.NoError(t, err)
	require.Equal(t, employee.ID, check.ID)
	require.Equal(t, employee.Name, check.Name)
	require.Equal(t, employee.Email, check.Email)
}

func TestEmployeeRepo_Update(t *testing.T) {
	tx := gdb.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	repo := NewEmployeeRepo(gdb)
	employee := MockEmployee()
	err = repo.Create(employee)
	require.NoError(t, err)

	employee.Name = gofakeit.Name()
	err = repo.Update(employee)
	require.NoError(t, err)
	check, err := repo.GetByID(employee.ID)
	require.NoError(t, err)
	require.Equal(t, employee.ID, check.ID)
	require.Equal(t, employee.Name, check.Name)
	require.Equal(t, employee.Email, check.Email)
}
