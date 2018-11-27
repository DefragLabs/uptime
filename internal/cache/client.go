package cache

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

// GetClient returns redis client.
func GetClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisDbNo := os.Getenv("REDIS_DB_NUMBER")
	redisPort := os.Getenv("REDIS_PORT")

	redisAddr := redisHost + ":" + redisPort

	redisDbNoInt, err := strconv.Atoi(redisDbNo)

	if err != nil {
		log.Fatal("Redis db no is not a valid integer")
	}

	client := redis.NewClient(
		&redis.Options{
			Addr: redisAddr,
			DB:   redisDbNoInt,
		},
	)

	return client
}
