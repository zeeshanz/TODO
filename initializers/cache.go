package initializers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisInstance struct {
	redis *redis.Client
}

var Cache RedisInstance

func testRedis(ctx context.Context) {
	SetToRedis(ctx, "key1", "Canada")
	SetToRedis(ctx, "key2", "Ottawa")
	val1, _ := GetFromRedis(ctx, "key1")
	val2, _ := GetFromRedis(ctx, "key2")
	fmt.Printf("First value with key `key1` should be Canada: %s \n", val1)
	fmt.Printf("First value with key `key2` should be Ottawa: %s \n", val2)
	values := getAllKeys(ctx, "key*")
	fmt.Printf("All keys : %v \n", values)
}

func pingRedis(ctx context.Context) {
	fmt.Println("PING")
	pong, err := Cache.redis.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(pong)
}

func ConnectRedis(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	Cache = RedisInstance{
		redis: client,
	}

	pingRedis(ctx)
	testRedis(ctx)
}

/*
 * Set values in Redis with a fixed time out.
 */
func SetToRedis(ctx context.Context, key, val string) {
	err := Cache.redis.Set(ctx, key, val, 1*time.Hour).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func GetFromRedis(ctx context.Context, key string) (string, error) {
	val, err := Cache.redis.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return val, nil
}

func DeleteFromRedis(ctx context.Context, key string) error {
	err := Cache.redis.Del(ctx, key)
	if err != nil {
		fmt.Println(err)
	}

	return err.Err()
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
