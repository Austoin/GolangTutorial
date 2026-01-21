// database/database_redis.go
// Redis 数据库操作 - 详细注释版

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ====== Redis 基础 ======
/*
Redis 是一个高性能的键值存储系统，常用于缓存、消息队列等场景。

主要特点：
1. 内存存储 - 高速读写
2. 丰富的数据结构 - String、Hash、List、Set、ZSet 等
3. 持久化 - 支持 RDB 和 AOF
4. 集群支持 - 支持主从、集群模式

安装 Redis 客户端：
  go get -u github.com/redis/go-redis/v9

数据类型：
  - String: 字符串
  - Hash: 哈希表
  - List: 列表
  - Set: 集合
  - ZSet: 有序集合
*/

// ====== Redis 客户端 ======

// RedisClient Redis 客户端封装
type RedisClient struct {
	client *redis.Client // Redis 客户端实例
	ctx    context.Context
}

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis 地址，如 "localhost:6379"
		Password: password, // 密码（为空表示不需要）
		DB:       db,       // 数据库编号
		PoolSize: 10,       // 连接池大小
	})

	// 创建上下文
	ctx := context.Background()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 Redis 失败: %w", err)
	}

	log.Printf("Redis 连接成功: %s", addr)

	return &RedisClient{
		client: client,
		ctx:    ctx,
	}, nil
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// Client 获取原生客户端
func (r *RedisClient) Client() *redis.Client {
	return r.client
}

// ====== String 操作 ======

// Set 设置字符串
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	// SET key value [EX seconds] [PX milliseconds] [NX|XX]
	// EX: 过期时间（秒）
	// PX: 过期时间（毫秒）
	// NX: 仅在 key 不存在时设置
	// XX: 仅在 key 存在时设置
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// Get 获取字符串
func (r *RedisClient) Get(key string) (string, error) {
	// GET key
	// 返回 key 对应的值，不存在返回 nil
	return r.client.Get(r.ctx, key).Result()
}

// GetInt 获取整数值
func (r *RedisClient) GetInt(key string) (int, error) {
	return r.client.Get(r.ctx, key).Int()
}

// GetFloat 获取浮点值
func (r *RedisClient) GetFloat(key string) (float64, error) {
	return r.client.Get(r.ctx, key).Float64()
}

// Incr 递增
func (r *RedisClient) Incr(key string) (int64, error) {
	// INCR key
	// 将 key 对应的值加 1，返回新的值
	return r.client.Incr(r.ctx, key).Result()
}

// IncrBy 递增指定值
func (r *RedisClient) IncrBy(key string, amount int64) (int64, error) {
	// INCRBY key increment
	return r.client.IncrBy(r.ctx, key, amount).Result()
}

// Decr 递减
func (r *RedisClient) Decr(key string) (int64, error) {
	// DECR key
	return r.client.Decr(r.ctx, key).Result()
}

// SetNX 仅在不存在时设置
func (r *RedisClient) SetNX(key string, value interface{}, expiration time.Duration) (bool, error) {
	// SET key value NX EX seconds
	return r.client.SetNX(r.ctx, key, value, expiration).Result()
}

// SetEX 设置带过期时间
func (r *RedisClient) SetEX(key string, value interface{}, expiration time.Duration) error {
	// SET key value EX seconds
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// MGet 批量获取
func (r *RedisClient) MGet(keys ...string) ([]interface{}, error) {
	// MGET key [key ...]
	return r.client.MGet(r.ctx, keys...).Result()
}

// MSet 批量设置
func (r *RedisClient) MSet(values ...interface{}) error {
	// MSET key value [key value ...]
	return r.client.MSet(r.ctx, values...).Err()
}

// ====== Hash 操作 ======

// HSet 设置哈希字段
func (r *RedisClient) HSet(key, field string, value interface{}) error {
	// HSET key field value
	return r.client.HSet(r.ctx, key, field, value).Err()
}

// HGet 获取哈希字段
func (r *RedisClient) HGet(key, field string) (string, error) {
	// HGET key field
	return r.client.HGet(r.ctx, key, field).Result()
}

// HGetAll 获取所有字段
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	// HGETALL key
	return r.client.HGetAll(r.ctx, key).Result()
}

// HMSet 批量设置哈希字段
func (r *RedisClient) HMSet(key string, values map[string]interface{}) error {
	// HMSET key field value [field value ...]
	return r.client.HMSet(r.ctx, key, values).Err()
}

// HMGet 批量获取哈希字段
func (r *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	// HMGET key field [field ...]
	return r.client.HMGet(r.ctx, key, fields...).Result()
}

// HIncrBy 字段值递增
func (r *RedisClient) HIncrBy(key, field string, amount int64) (int64, error) {
	// HINCRBY key field increment
	return r.client.HIncrBy(r.ctx, key, field, amount).Result()
}

// HExists 检查字段是否存在
func (r *RedisClient) HExists(key, field string) (bool, error) {
	// HEXISTS key field
	return r.client.HExists(r.ctx, key, field).Result()
}

// HDel 删除字段
func (r *RedisClient) HDel(key string, fields ...string) error {
	// HDEL key field [field ...]
	return r.client.HDel(r.ctx, key, fields...).Err()
}

// HLen 获取字段数量
func (r *RedisClient) HLen(key string) (int64, error) {
	// HLEN key
	return r.client.HLen(r.ctx, key).Result()
}

// ====== List 操作 ======

// LPush 从左侧插入
func (r *RedisClient) LPush(key string, values ...interface{}) (int64, error) {
	// LPUSH key value [value ...]
	return r.client.LPush(r.ctx, key, values...).Result()
}

// RPush 从右侧插入
func (r *RedisClient) RPush(key string, values ...interface{}) (int64, error) {
	// RPUSH key value [value ...]
	return r.client.RPush(r.ctx, key, values...).Result()
}

// LPop 从左侧弹出
func (r *RedisClient) LPop(key string) (string, error) {
	// LPOP key
	return r.client.LPop(r.ctx, key).Result()
}

// RPop 从右侧弹出
func (r *RedisClient) RPop(key string) (string, error) {
	// RPOP key
	return r.client.RPop(r.ctx, key).Result()
}

// LRange 获取列表范围
func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	// LRANGE key start stop
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (r *RedisClient) LLen(key string) (int64, error) {
	// LLEN key
	return r.client.LLen(r.ctx, key).Result()
}

// LIndex 获取指定索引元素
func (r *RedisClient) LIndex(key string, index int64) (string, error) {
	// LINDEX key index
	return r.client.LIndex(r.ctx, key, index).Result()
}

// LSet 设置指定索引元素
func (r *RedisClient) LSet(key string, index int64, value interface{}) error {
	// LSET key index value
	return r.client.LSet(r.ctx, key, index, value).Err()
}

// ====== Set 操作 ======

// SAdd 添加集合成员
func (r *RedisClient) SAdd(key string, members ...interface{}) (int64, error) {
	// SADD key member [member ...]
	return r.client.SAdd(r.ctx, key, members...).Result()
}

// SMembers 获取所有成员
func (r *RedisClient) SMembers(key string) ([]string, error) {
	// SMEMBERS key
	return r.client.SMembers(r.ctx, key).Result()
}

// SIsMember 检查成员是否存在
func (r *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
	// SISMEMBER key member
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// SCard 获取集合基数（大小）
func (r *RedisClient) SCard(key string) (int64, error) {
	// SCARD key
	return r.client.SCard(r.ctx, key).Result()
}

// SRem 移除成员
func (r *RedisClient) SRem(key string, members ...interface{}) (int64, error) {
	// SREM key member [member ...]
	return r.client.SRem(r.ctx, key, members...).Result()
}

// SInter 求交集
func (r *RedisClient) SInter(keys ...string) ([]string, error) {
	// SINTER key [key ...]
	return r.client.SInter(r.ctx, keys...).Result()
}

// SUnion 求并集
func (r *RedisClient) SUnion(keys ...string) ([]string, error) {
	// SUNION key [key ...]
	return r.client.SUnion(r.ctx, keys...).Result()
}

// SDiff 求差集
func (r *RedisClient) SDiff(keys ...string) ([]string, error) {
	// SDIFF key [key ...]
	return r.client.SDiff(r.ctx, keys...).Result()
}

// ====== ZSet 操作 ======

// ZAdd 添加有序集合成员
func (r *RedisClient) ZAdd(key string, members ...redis.Z) (int64, error) {
	// ZADD key [NX|XX] [CH] [INCR] score member [score member ...]
	return r.client.ZAdd(r.ctx, key, members...).Result()
}

// ZRange 获取范围成员
func (r *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	// ZRANGE key start stop [BYSCORE | REV] [LIMIT offset count] [WITHSCORES]
	return r.client.ZRange(r.ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取范围成员及分数
func (r *RedisClient) ZRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	return r.client.ZRangeWithScores(r.ctx, key, start, stop).Result()
}

// ZScore 获取成员分数
func (r *RedisClient) ZScore(key, member string) (float64, error) {
	// ZSCORE key member
	return r.client.ZScore(r.ctx, key, member).Result()
}

// ZIncrBy 递增成员分数
func (r *RedisClient) ZIncrBy(key string, increment float64, member string) (float64, error) {
	// ZINCRBY key increment member
	return r.client.ZIncrBy(r.ctx, key, increment, member).Result()
}

// ZCard 获取基数
func (r *RedisClient) ZCard(key string) (int64, error) {
	// ZCARD key
	return r.client.ZCard(r.ctx, key).Result()
}

// ZRem 移除成员
func (r *RedisClient) ZRem(key string, members ...interface{}) (int64, error) {
	// ZREM key member [member ...]
	return r.client.ZRem(r.ctx, key, members...).Result()
}

// ZRank 获取成员排名
func (r *RedisClient) ZRank(key, member string) (int64, error) {
	// ZRANK key member
	return r.client.ZRank(r.ctx, key, member).Result()
}

// ====== 键操作 ======

// Exists 检查键是否存在
func (r *RedisClient) Exists(keys ...string) (int64, error) {
	// EXISTS key [key ...]
	return r.client.Exists(r.ctx, keys...).Result()
}

// Del 删除键
func (r *RedisClient) Del(keys ...string) (int64, error) {
	// DEL key [key ...]
	return r.client.Del(r.ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisClient) Expire(key string, expiration time.Duration) (bool, error) {
	// EXPIRE key seconds
	return r.client.Expire(r.ctx, key, expiration).Result()
}

// TTL 获取剩余过期时间
func (r *RedisClient) TTL(key string) (time.Duration, error) {
	// TTL key
	return r.client.TTL(r.ctx, key).Result()
}

// Rename 重命名键
func (r *RedisClient) Rename(key, newkey string) error {
	// RENAME key newkey
	return r.client.Rename(r.ctx, key, newkey).Err()
}

// Type 获取键类型
func (r *RedisClient) Type(key string) (string, error) {
	// TYPE key
	return r.client.Type(r.ctx, key).Result()
}

// ====== 过期操作 ======

// Persist 移除过期时间
func (r *RedisClient) Persist(key string) (bool, error) {
	// PERSIST key
	return r.client.Persist(r.ctx, key).Result()
}

// ====== 管道操作 ======

// Pipeline 管道操作示例
func (r *RedisClient) PipelineExample() {
	// 创建管道
	pipe := r.client.Pipeline()

	// 添加命令到管道
	incr := pipe.Incr(r.ctx, "pipeline_counter")
	pipe.Expire(r.ctx, "pipeline_counter", time.Hour)

	// 执行管道中的所有命令
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		log.Printf("管道执行失败: %v", err)
		return
	}

	// 获取结果
	fmt.Println("管道计数器:", incr.Val())
}

// ====== 事务操作 ======

// TxPipelined 事务管道
func (r *RedisClient) TxPipelinedExample() {
	var incr int64

	// Watch 监控键
	err := r.client.Watch(r.ctx, func(tx *redis.Tx) error {
		// 获取当前值
		n, err := tx.Get(r.ctx, "counter").Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// 事务操作
		_, err = tx.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
			incr = int64(n + 1)
			pipe.Set(r.ctx, "counter", incr, 0)
			return nil
		})

		return err
	}, "counter")

	if err != nil {
		log.Printf("事务失败: %v", err)
	}

	fmt.Println("计数器值:", incr)
}

// ====== 发布订阅 ======

// PubSub 发布订阅示例
func (r *RedisClient) PubSubExample() {
	// 订阅
	pubsub := r.client.Subscribe(r.ctx, "mychannel")

	// 等待订阅成功
	_, err := pubsub.Receive(r.ctx)
	if err != nil {
		log.Printf("订阅失败: %v", err)
		return
	}

	// 获取通道
	ch := pubsub.Channel()

	// 发布消息
	r.client.Publish(r.ctx, "mychannel", "Hello, Redis!")

	// 接收消息
	msg := <-ch
	fmt.Println("收到消息:", msg.Payload)

	// 关闭订阅
	pubsub.Close()
}

// ====== 分布式锁 ======

// Lock 尝试获取分布式锁
func (r *RedisClient) Lock(key string, value string, expiration time.Duration) (bool, error) {
	// SET key value NX EX seconds
	return r.client.SetNX(r.ctx, key, value, expiration).Result()
}

// Unlock 释放锁（使用 Lua 脚本确保原子性）
func (r *RedisClient) Unlock(key, value string) error {
	// Lua 脚本：检查值并删除
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	return script.Run(r.ctx, r.client, []string{key}, value).Err()
}

// ====== 缓存示例 ======

// CacheUser 缓存用户信息
func (r *RedisClient) CacheUser(userID int, userData map[string]interface{}, expiration time.Duration) error {
	key := fmt.Sprintf("user:%d", userID)
	return r.HMSet(key, userData)
}

// GetCachedUser 获取缓存的用户信息
func (r *RedisClient) GetCachedUser(userID int) (map[string]string, error) {
	key := fmt.Sprintf("user:%d", userID)
	return r.HGetAll(key)
}

// ====== 主函数 ======

func main() {
	fmt.Println("=== Redis 数据库操作示例 ===")

	// 1. 创建 Redis 客户端
	client, err := NewRedisClient("localhost:6379", "", 0)
	if err != nil {
		log.Fatalf("连接 Redis 失败: %v", err)
	}
	defer client.Close()

	// 2. String 操作示例
	fmt.Println("\n--- String 操作 ---")

	// 设置和获取
	client.Set("name", "Alice", time.Hour)
	name, _ := client.Get("name")
	fmt.Printf("name = %s\n", name)

	// 计数器
	client.Set("counter", 0, 0)
	client.Incr("counter")
	client.IncrBy("counter", 5)
	counter, _ := client.GetInt("counter")
	fmt.Printf("counter = %d\n", counter)

	// 3. Hash 操作示例
	fmt.Println("\n--- Hash 操作 ---")

	client.HSet("user:1", "name", "Bob")
	client.HSet("user:1", "age", "25")
	client.HSet("user:1", "email", "bob@example.com")

	user, _ := client.HGetAll("user:1")
	fmt.Printf("user:1 = %v\n", user)

	// 4. List 操作示例
	fmt.Println("\n--- List 操作 ---")

	client.RPush("tasks", "task1")
	client.RPush("tasks", "task2")
	client.RPush("tasks", "task3")

	tasks, _ := client.LRange("tasks", 0, -1)
	fmt.Printf("tasks = %v\n", tasks)

	// 5. Set 操作示例
	fmt.Println("\n--- Set 操作 ---")

	client.SAdd("tags", "go", "redis", "database")
	client.SAdd("tags", "go") // 重复添加无效

	tags, _ := client.SMembers("tags")
	fmt.Printf("tags = %v\n", tags)

	// 6. ZSet 操作示例
	fmt.Println("\n--- ZSet 操作 ---")

	client.ZAdd("leaderboard",
		redis.Z{Score: 100, Member: "Alice"},
		redis.Z{Score: 200, Member: "Bob"},
		redis.Z{Score: 150, Member: "Charlie"},
	)

	leaderboard, _ := client.ZRangeWithScores("leaderboard", 0, -1)
	fmt.Println("leaderboard:")
	for _, z := range leaderboard {
		fmt.Printf("  %s: %.0f\n", z.Member, z.Score)
	}

	// 7. 键操作示例
	fmt.Println("\n--- 键操作 ---")

	exists, _ := client.Exists("name")
	fmt.Printf("name exists: %d\n", exists)

	// 8. 清理测试数据
	client.Del("name", "counter", "user:1", "tasks", "tags", "leaderboard")

	fmt.Println("\nRedis 操作示例完成")
}
