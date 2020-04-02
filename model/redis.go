package model

import (
	"log"
	"time"

	redis "github.com/go-redis/redis"
)

// RedisOpen 是否连接 redis
var RedisOpen = false

// RedisCli redis客户端
var RedisCli client

type client interface {
	Ping() *redis.StatusCmd
	Close() error
	// Expire 设置key超时时间
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	// HExists 判断是否存在
	HExists(key, field string) *redis.BoolCmd
	// HMset 设置值
	HMSet(key string, values ...interface{}) *redis.BoolCmd
	// HMGet 获取值
	HMGet(key string, fields ...string) *redis.SliceCmd
	// HScan 分页查询
	HScan(key string, cursor uint64, match string, count int64) *redis.ScanCmd
	// HLen hashmap字段数
	HLen(key string) *redis.IntCmd
	// Sort 排序
	Sort(key string, sort *redis.Sort) *redis.StringSliceCmd
	// SAdd 添加值到集合
	SAdd(key string, members ...interface{}) *redis.IntCmd
	// SIsMember 是否是集合成员
	SIsMember(key string, member interface{}) *redis.BoolCmd
	// Pipeline 管道
	Pipeline() redis.Pipeliner
	Watch(fn func(*redis.Tx) error, keys ...string) error
	Get(key string) *redis.StringCmd
}

// SetRedis 设置redis
func SetRedis() {
	log.Println("启动redis")
	if conf.RedisCluster == "true" {
		// clusterIsOpen = true
		RedisCli = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{conf.RedisHost + ":" + conf.RedisPort},
			Password: conf.RedisPassword,
		})
		_, err := RedisCli.Ping().Result()
		if err != nil {
			log.Printf("连接 redis cluster：%s 失败,原因：%v\n", conf.RedisHost+":"+conf.RedisPort, err)
		} else {
			RedisOpen = true
		}

	} else {
		RedisCli = redis.NewClient(&redis.Options{
			Addr:     conf.RedisHost + ":" + conf.RedisPort,
			Password: conf.RedisPassword,
		})
		_, err := RedisCli.Ping().Result()
		if err != nil {
			log.Printf("连接 redis：%s 失败,原因：%v\n", conf.RedisHost+":"+conf.RedisPort, err)
		} else {
			RedisOpen = true
		}
	}
}
