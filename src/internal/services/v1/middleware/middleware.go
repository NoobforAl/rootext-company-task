package middleware

import "github.com/redis/go-redis/v9"

type middleware struct {
	redisClient *redis.Client
}

func New(redisClient *redis.Client) *middleware {
	return &middleware{redisClient: redisClient}
}
