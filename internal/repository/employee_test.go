package repository

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/model"
	"github.com/joremysh/fliqt/pkg/database"
)

var (
	gdb *gorm.DB
	err error
)

func TestMain(m *testing.M) {
	dsn := flag.String("dsn", "user:password@tcp(localhost:3306)/hrs?collation=utf8_unicode_ci&parseTime=true&loc=Asia%2FTaipei&multiStatements=true", "database connection string")
	flag.Parse()

	gdb, err = database.NewDatabase(*dsn)
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(m.Run())
}

func TestEmployeeRepo_Create(t *testing.T) {
	tx := gdb.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})

	repo := NewEmployeeRepo(gdb)
	employee := mockEmployee()
	err = repo.Create(employee)
	require.NoError(t, err)
	require.NotNil(t, employee.ID)
	// jsonBytes, _ := json.Marshal(employee)
	// t.Log(string(jsonBytes))

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
	employee := mockEmployee()
	err = repo.Create(employee)
	require.NoError(t, err)
	require.NotNil(t, employee.ID)

	check, err := repo.GetByID(employee.ID)
	require.NoError(t, err)
	require.Equal(t, employee.ID, check.ID)
	require.Equal(t, employee.Name, check.Name)
	require.Equal(t, employee.Email, check.Email)
}
