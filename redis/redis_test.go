package redis

import (
	"testing"
	"time"
)

const USER = "jhon snow"

func TestRedis(t *testing.T) {
	redisContainer := NewRedisContainer()

	response := redisContainer.Redis.Set("test_user", USER, time.Duration(time.Duration.Minutes(1)))

	if response.Err() != nil {
		ErrPurge(t, redisContainer)
	}

	result := redisContainer.Redis.Get("test_user")

	if USER != result.Val() {
		ErrPurge(t, redisContainer)
	}

	redisContainer.Purge()
}

func ErrPurge(t *testing.T, container *RedisContainer) {
	container.Purge()
	t.Error()
}
