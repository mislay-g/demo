package storage

import (
	"context"
	"demo/src/config"
	"demo/src/pkg/ormx"
	"fmt"
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func InitDB(cfg *config.DBConfig) error {
	db, err := ormx.New(cfg)
	if err == nil {
		DB = db
	}
	return err
}

var Redis interface {
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	Close() error
	Ping(ctx context.Context) *redis.StatusCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

func InitRedis(cfg *config.RedisConfig) (func(), error) {
	redisOptions := &redis.Options{
		Addr:     cfg.Address,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	Redis = redis.NewClient(redisOptions)

	return func() {
		fmt.Println("redis exiting")
		err := Redis.Close()
		if err != nil {
			return
		}
	}, nil
}
