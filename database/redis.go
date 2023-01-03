package database

import (
	"github.com/go-redis/redis/v9"
)

// "fmt"

var RDB *redis.Client

func ConnectRedis() {
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RDB = r
}
