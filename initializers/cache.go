package initializers

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisInstance struct {
	redis *redis.Client
}

var Cache RedisInstance

func test() {
	ctx := context.Background()
	ConnectRedis(ctx)

	SetToRedis(ctx, "name", "redis-test")
	SetToRedis(ctx, "name2", "redis-test-2")
	val := GetFromRedis(ctx, "name")

	fmt.Printf("First value with name key : %s \n", val)

	values := getAllKeys(ctx, "name*")

	fmt.Printf("All values : %v \n", values)

}

func ConnectRedis(ctx context.Context) {
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

	Cache = RedisInstance{
		redis: client,
	}
}

func SetToRedis(ctx context.Context, key, val string) {
	err := Cache.redis.Set(ctx, key, val, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func GetFromRedis(ctx context.Context, key string) string {
	val, err := Cache.redis.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
	}

	return val
}

func getAllKeys(ctx context.Context, key string) []string {
	keys := []string{}

	iter := Cache.redis.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	return keys
}
