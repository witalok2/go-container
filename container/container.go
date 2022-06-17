package container

import (
	"runtime"
	"time"

	"github.com/ory/dockertest/v3"
	logger "github.com/sirupsen/logrus"
)

type Container struct {
	pool     *dockertest.Pool     // API connection to host docker
	resource *dockertest.Resource // The docker container
}

func NewContainer(runOptions dockertest.RunOptions) *Container {
	poolEndpoint := ``
	if runtime.GOOS == "windows" {
		poolEndpoint = "npipe:////./pipe/docker_engine"
	}

	pool, err := dockertest.NewPool(poolEndpoint)
	if err != nil {
		logger.WithError(err).Fatal("could not connect to docker")
	}

	pool.MaxWait = 30 * time.Second

	resource, err := pool.RunWithOptions(&runOptions)
	if err != nil {
		logger.WithError(err).Fatal("could not start postgres container")
	}

	err = resource.Expire(30 * 60)
	if err != nil {
		logger.WithError(err).Fatal("could not set container expiration")
	}

	return &Container{
		pool: pool,
		resource: resource,
	}
}
