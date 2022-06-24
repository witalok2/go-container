package redis

import (
	"net/url"

	"github.com/go-redis/redis"
	"github.com/ory/dockertest/v3"
	logger "github.com/sirupsen/logrus"
	"github.com/witalok2/go-container/container"
)

type RedisContainer struct {
	pool     *dockertest.Pool
	resource *dockertest.Resource
	URL      *url.URL
	Redis    *redis.Client
}

func NewRedisContainer() *RedisContainer {
	var rd *redis.Client
	redisURL := &url.URL{}

	dockerContainer := container.NewContainer(dockertest.RunOptions{
		Repository:   "redis",
		Env:          []string{},
		ExposedPorts: []string{"6379/tcp"},
	})

	err := dockerContainer.Resource.Expire(30 * 60)
	if err != nil {
		logger.WithError(err).Fatalf("could not expire dockerContainer")
	}

	redisURL.Host = dockerContainer.Resource.Container.NetworkSettings.IPAddress

	err = dockerContainer.Pool.Retry(func() error {
		rd = redis.NewClient(&redis.Options{
			Addr: dockerContainer.Resource.Container.NetworkSettings.IPAddress + ":6379",
		})

		err := rd.Ping().Err()
		if err != nil {
			return err
		}
		return rd.Ping().Err()
	})
	if err != nil {
		logger.WithError(err).Fatalf("could not connect to redis server: %v", err)
	}

	return &RedisContainer{
		pool:     dockerContainer.Pool,
		resource: dockerContainer.Resource,
		URL:      redisURL,
		Redis:    rd,
	}
}

// Purge removes the container and associated volumes from docker
func (c *RedisContainer) Purge() {
	if err := c.pool.Purge(c.resource); err != nil {
		logger.WithError(err).Fatal("Could not purge resource")
	}
}
