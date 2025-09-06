package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"game-backend/config"
)

var RedisClient *redis.Client

// InitRedis 初始化Redis连接
func InitRedis() error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.AppConfig.Redis.GetRedisAddr(),
		Password:     config.AppConfig.Redis.Password,
		DB:           config.AppConfig.Redis.DB,
		PoolSize:     config.AppConfig.Redis.PoolSize,
		MinIdleConns: 5,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	log.Println("Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Close()
}

// GetRedisClient 获取Redis客户端
func GetRedisClient() *redis.Client {
	return RedisClient
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return RedisClient.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	result, err := RedisClient.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置键过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RedisClient.Expire(ctx, key, expiration).Err()
}

// HSet 设置哈希字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希字段值
func HGet(ctx context.Context, key, field string) (string, error) {
	return RedisClient.HGet(ctx, key, field).Result()
}

// HGetAll 获取所有哈希字段
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return RedisClient.HGetAll(ctx, key).Result()
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return RedisClient.HDel(ctx, key, fields...).Err()
}

// LPush 从列表左侧推入元素
func LPush(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.LPush(ctx, key, values...).Err()
}

// RPush 从列表右侧推入元素
func RPush(ctx context.Context, key string, values ...interface{}) error {
	return RedisClient.RPush(ctx, key, values...).Err()
}

// LPop 从列表左侧弹出元素
func LPop(ctx context.Context, key string) (string, error) {
	return RedisClient.LPop(ctx, key).Result()
}

// RPop 从列表右侧弹出元素
func RPop(ctx context.Context, key string) (string, error) {
	return RedisClient.RPop(ctx, key).Result()
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
	return RedisClient.LLen(ctx, key).Result()
}

// SAdd 向集合添加成员
func SAdd(ctx context.Context, key string, members ...interface{}) error {
	return RedisClient.SAdd(ctx, key, members...).Err()
}

// SRem 从集合移除成员
func SRem(ctx context.Context, key string, members ...interface{}) error {
	return RedisClient.SRem(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func SMembers(ctx context.Context, key string) ([]string, error) {
	return RedisClient.SMembers(ctx, key).Result()
}

// SIsMember 检查成员是否在集合中
func SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return RedisClient.SIsMember(ctx, key, member).Result()
}

// ZAdd 向有序集合添加成员
func ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	return RedisClient.ZAdd(ctx, key, members...).Err()
}

// ZRem 从有序集合移除成员
func ZRem(ctx context.Context, key string, members ...interface{}) error {
	return RedisClient.ZRem(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围内的成员
func ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RedisClient.ZRange(ctx, key, start, stop).Result()
}

// ZRevRange 获取有序集合范围内成员（倒序）
func ZRevRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return RedisClient.ZRevRange(ctx, key, start, stop).Result()
}

// ZRank 获取成员在有序集合中的排名
func ZRank(ctx context.Context, key string, member string) (int64, error) {
	return RedisClient.ZRank(ctx, key, member).Result()
}

// ZScore 获取成员在有序集合中的分数
func ZScore(ctx context.Context, key string, member string) (float64, error) {
	return RedisClient.ZScore(ctx, key, member).Result()
}
