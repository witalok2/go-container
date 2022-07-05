package container

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

type ContainerTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}

func (suite *ContainerTestSuite) SetupSuite() {
	suite.ctx = context.Background()
}

func (suite *ContainerTestSuite) TestMongoContainer() {
	dockerContainer := NewContainer(dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "3.3.12",
		ExposedPorts: []string{"3000"},
		Cmd:          []string{"mongod", "--smallfiles", "--port", "3000"},
	})

	defer dockerContainer.Pool.Purge(dockerContainer.Resource)

	port := dockerContainer.Resource.GetPort("3000/tcp")
	suite.NotEmpty(port)

	err := dockerContainer.Pool.Retry(func() error {
		response, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s", port))
		if err != nil {
			return err
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("could not connect to resource")
		}

		defer response.Body.Close()

		return nil
	})
	suite.Require().Nil(err)
}

func (suite *ContainerTestSuite) TestPostgresContainer() {
	dockerContainer := NewContainer(dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "9.5",
		Env:          []string{"POSTGRES_PASSWORD=secret"},
		ExposedPorts: []string{"5432"},
	})

	defer dockerContainer.Pool.Purge(dockerContainer.Resource)

	port := dockerContainer.Resource.GetPort("5432")
	suite.NotEmpty(port)

	err := dockerContainer.Pool.Retry(func() error {
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://postgres:secret@localhost:%s/postgres?sslmode=disable", port))
		if err != nil {
			return err
		}
		return db.Ping()
	})
	suite.Require().Nil(err)
}

func (suite *ContainerTestSuite) TestRedisContainer() {
	dockerContainer := NewContainer(dockertest.RunOptions{
		Repository:   "redis",
		Env:          []string{},
		ExposedPorts: []string{"6379"},
	})

	defer dockerContainer.Pool.Purge(dockerContainer.Resource)

	err := dockerContainer.Pool.Retry(func() error {
		rd := redis.NewClient(&redis.Options{
			Addr: dockerContainer.Resource.Container.NetworkSettings.IPAddress + ":6379",
		})

		err := rd.Ping().Err()
		if err != nil {
			return err
		}
		return rd.Ping().Err()
	})
	suite.Require().Nil(err)
}
