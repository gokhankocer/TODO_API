package database

import (
	"os"

	"github.com/go-redis/redis/v9"
)

// "fmt"

var RDB *redis.Client

func ConnectRedis() {
	r := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	RDB = r
}
