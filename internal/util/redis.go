package util

import (
	"github.com/go-redis/redis/v8"

	"fliqt/config"
)

func NewClient(cfg *config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}
