package controller

import (
	"ratblog/config"
	"ratblog/contract"

	"github.com/redis/go-redis/v9"
)

type Controller struct {
	repo contract.Repository

	redisClient *redis.Client
}

type response struct {
	Message string `json:"message"`
}

var badRequestResponse = response{Message: "Bad Request"}
var notImplementedResponse = response{Message: "Not Implemented"}
var usernameOrPasswordIncorrectResponse = response{Message: "Username or Password is incorrect"}
var userAlreadyExistResponse = response{Message: "User already exist"}
var unauthorizedResponse = response{Message: "Unauthorized"}

func New(repo contract.Repository, cfg config.Config) *Controller {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	return &Controller{
		repo:        repo,
		redisClient: redisClient,
	}
}
