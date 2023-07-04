package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"server/config"
	"time"
)

/*
redis 工具类
*/
var Rdb *redis.Client
var Cxt = context.Background()

// InitRedisConn 初始化redis客户端
func InitRedisConn() error {

	Rdb = redis.NewClient(&redis.Options{
		Addr:        config.RedisAddr,
		Password:    config.RedisPassword,
		DB:          config.RedisDBNo,
		PoolSize:    10,               // 最大连接数
		DialTimeout: time.Second * 10, // 超时时间
	})
	// 测试连接是否正常
	_, err := Rdb.Ping(Cxt).Result()
	if err != nil {
		panic(err)
	}
	return nil
}

// CloseRedis 关闭redis连接
func CloseRedis() error {
	return Rdb.Close()
}
