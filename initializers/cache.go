package initializers

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func init() {
	ctx := context.TODO()
	connectRedis(ctx)

	setToRedis(ctx, "name", "redis-test")
	setToRedis(ctx, "name2", "redis-test-2")
	val := getFromRedis(ctx, "name")

	fmt.Printf("First value with name key : %s \n", val)

	values := getAllKeys(ctx, "name*")

	fmt.Printf("All values : %v \n", values)

}

func connectRedis(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pong)

	redisClient = client
}

func setToRedis(ctx context.Context, key, val string) {
	err := redisClient.Set(ctx, key, val, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func getFromRedis(ctx context.Context, key string) string {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return val
}

func getAllKeys(ctx context.Context, key string) []string {
	keys := []string{}

	iter := redisClient.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return keys
}
