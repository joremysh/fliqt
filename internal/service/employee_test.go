package service

import (
	"context"
	"log"
	"testing"

	"github.com/go-redis/redismock/v9"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/joremysh/fliqt/internal/repository"
	"github.com/joremysh/fliqt/pkg/cache"
	testingx "github.com/joremysh/fliqt/pkg/testing"
)

var (
	gdb      *gorm.DB
	err      error
	pool     *dockertest.Pool
	resource *dockertest.Resource
	client   *redis.Client
	repo     repository.Employee
)

func TestMain(m *testing.M) {
	pool, resource, gdb, err = testingx.NewMysqlInDocker()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = repository.Migrate(gdb)
	if err != nil {
		log.Fatal(err.Error())
	}
	client, _ = redismock.NewClientMock()

	m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestEmployeeService_CreateEmployee(t *testing.T) {
	tx := gdb.Begin()
	t.Cleanup(func() {
		tx.Rollback()
	})
	repo = repository.NewEmployeeRepo(tx)

	svc := NewEmployeeService(repo, &cache.RedisClient{Client: client})
	employee := repository.MockEmployee()
	ctx := context.Background()
	created, err := svc.CreateEmployee(ctx, employee)
	require.NoError(t, err)
	require.Equal(t, employee.Name, created.Name)
	require.Equal(t, employee.Email, created.Email)
	require.Equal(t, employee.PhoneNumber, created.PhoneNumber)
}
