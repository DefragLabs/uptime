package cache

import (
	"os"
	"github.com/go-redis/redis"
)

func getClient() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisDbNo := os.Getenv("REDIS_DB_NUMBER")
	redisPort := os.Getenv("REDIS_PORT")

	redisAddr := redisHost + ":" + redisPort

	client := redis.NewClient(
		&redis.Options{
			Addr: redisAddr,
			Db: redisDbNo,
		}
	)

	return client
}
