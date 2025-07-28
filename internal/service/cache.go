package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	PoolSize int    `json:"poolSize"`
}

// CacheService 缓存服务
type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheService 创建缓存服务
func NewCacheService(config CacheConfig) *CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: config.PoolSize,
	})

	return &CacheService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// TestConnection 测试连接
func (c *CacheService) TestConnection() error {
	return c.client.Ping(c.ctx).Err()
}

// Set 设置缓存
func (c *CacheService) Set(key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}

	return c.client.Set(c.ctx, key, data, expiration).Err()
}

// Get 获取缓存
func (c *CacheService) Get(key string, dest interface{}) error {
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("缓存不存在")
		}
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// GetString 获取字符串缓存
func (c *CacheService) GetString(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

// SetString 设置字符串缓存
func (c *CacheService) SetString(key, value string, expiration time.Duration) error {
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

// Delete 删除缓存
func (c *CacheService) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists 检查缓存是否存在
func (c *CacheService) Exists(key string) (bool, error) {
	count, err := c.client.Exists(c.ctx, key).Result()
	return count > 0, err
}

// Expire 设置过期时间
func (c *CacheService) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(c.ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (c *CacheService) TTL(key string) (time.Duration, error) {
	return c.client.TTL(c.ctx, key).Result()
}

// Increment 递增
func (c *CacheService) Increment(key string) (int64, error) {
	return c.client.Incr(c.ctx, key).Result()
}

// IncrementBy 按指定值递增
func (c *CacheService) IncrementBy(key string, value int64) (int64, error) {
	return c.client.IncrBy(c.ctx, key, value).Result()
}

// Decrement 递减
func (c *CacheService) Decrement(key string) (int64, error) {
	return c.client.Decr(c.ctx, key).Result()
}

// DecrementBy 按指定值递减
func (c *CacheService) DecrementBy(key string, value int64) (int64, error) {
	return c.client.DecrBy(c.ctx, key, value).Result()
}

// SetNX 仅当key不存在时设置
func (c *CacheService) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化失败: %v", err)
	}

	return c.client.SetNX(c.ctx, key, data, expiration).Result()
}

// GetSet 设置新值并返回旧值
func (c *CacheService) GetSet(key string, value interface{}) (string, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("序列化失败: %v", err)
	}

	return c.client.GetSet(c.ctx, key, data).Result()
}

// MSet 批量设置
func (c *CacheService) MSet(pairs map[string]interface{}) error {
	args := make([]interface{}, 0, len(pairs)*2)
	for key, value := range pairs {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("序列化失败 %s: %v", key, err)
		}
		args = append(args, key, data)
	}

	return c.client.MSet(c.ctx, args...).Err()
}

// MGet 批量获取
func (c *CacheService) MGet(keys []string) ([]interface{}, error) {
	return c.client.MGet(c.ctx, keys...).Result()
}

// Keys 获取匹配的键
func (c *CacheService) Keys(pattern string) ([]string, error) {
	return c.client.Keys(c.ctx, pattern).Result()
}

// FlushDB 清空当前数据库
func (c *CacheService) FlushDB() error {
	return c.client.FlushDB(c.ctx).Err()
}

// FlushAll 清空所有数据库
func (c *CacheService) FlushAll() error {
	return c.client.FlushAll(c.ctx).Err()
}

// Close 关闭连接
func (c *CacheService) Close() error {
	return c.client.Close()
}

// 验证码相关方法

// SetVerificationCode 设置验证码
func (c *CacheService) SetVerificationCode(target, code string, expiration time.Duration) error {
	key := fmt.Sprintf("verification_code:%s", target)
	return c.SetString(key, code, expiration)
}

// GetVerificationCode 获取验证码
func (c *CacheService) GetVerificationCode(target string) (string, error) {
	key := fmt.Sprintf("verification_code:%s", target)
	return c.GetString(key)
}

// DeleteVerificationCode 删除验证码
func (c *CacheService) DeleteVerificationCode(target string) error {
	key := fmt.Sprintf("verification_code:%s", target)
	return c.Delete(key)
}

// CheckVerificationCodeLimit 检查验证码发送限制
func (c *CacheService) CheckVerificationCodeLimit(target string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("verification_limit:%s", target)

	count, err := c.Increment(key)
	if err != nil {
		return false, err
	}

	if count == 1 {
		// 第一次设置过期时间
		if err := c.Expire(key, window); err != nil {
			return false, err
		}
	}

	return count <= int64(limit), nil
}

// 会话相关方法

// SetSession 设置会话
func (c *CacheService) SetSession(sessionID string, data interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Set(key, data, expiration)
}

// GetSession 获取会话
func (c *CacheService) GetSession(sessionID string, dest interface{}) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Get(key, dest)
}

// DeleteSession 删除会话
func (c *CacheService) DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Delete(key)
}

// RefreshSession 刷新会话过期时间
func (c *CacheService) RefreshSession(sessionID string, expiration time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return c.Expire(key, expiration)
}

// 锁相关方法

// Lock 获取分布式锁
func (c *CacheService) Lock(key string, expiration time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock:%s", key)
	return c.client.SetNX(c.ctx, lockKey, "locked", expiration).Result()
}

// Unlock 释放分布式锁
func (c *CacheService) Unlock(key string) error {
	lockKey := fmt.Sprintf("lock:%s", key)
	return c.Delete(lockKey)
}

// 统计相关方法

// IncrementCounter 递增计数器
func (c *CacheService) IncrementCounter(key string, expiration time.Duration) (int64, error) {
	count, err := c.Increment(key)
	if err != nil {
		return 0, err
	}

	if count == 1 {
		// 第一次设置过期时间
		if err := c.Expire(key, expiration); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// GetCounter 获取计数器值
func (c *CacheService) GetCounter(key string) (int64, error) {
	val, err := c.client.Get(c.ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}

// GetCacheStats 获取缓存统计信息
func (c *CacheService) GetCacheStats() (map[string]interface{}, error) {
	info, err := c.client.Info(c.ctx, "memory", "stats").Result()
	if err != nil {
		return nil, err
	}

	dbSize, err := c.client.DBSize(c.ctx).Result()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"dbSize": dbSize,
		"info":   info,
	}, nil
}
