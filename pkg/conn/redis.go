package conn

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	*redis.Client
}

func (r *Redis) ConnectRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_URL"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASS"),
	})
	r.Client = client
	if client != nil {
		log.Println("Redis Connected")
	} else {
		log.Println("Failed Redis Connection")
	}
}

func (r *Redis) CloseRedis() {
	if err := r.Close(); err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Redis Connection Closed")
	}
}

func (r *Redis) RedisWrite() error {
	return r.Set(context.Background(), "example-key", "example-value", 0).Err()
}

func (r *Redis) RedisRead() (string, error) {
	value, err := r.Get(context.Background(), "example-key").Result()
	return value, err
}
