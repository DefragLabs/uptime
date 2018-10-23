package cache

import (
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func getClient() *redis.Client {
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
