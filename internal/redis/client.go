package redis

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// RedisClient abstracts Redis operations for testability.
type RedisClient interface {
	Ping(ctx context.Context) error
	Info(ctx context.Context, sections ...string) (string, error)
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	ObjectIdleTime(ctx context.Context, key string) (time.Duration, error)
	MemoryUsage(ctx context.Context, key string) (int64, error)
	SlowLogGet(ctx context.Context, num int64) ([]SlowLogEntry, error)
	ConfigGet(ctx context.Context, parameter string) (map[string]string, error)
	DBSize(ctx context.Context) (int64, error)
	Close() error
}

// SlowLogEntry represents a single slow log entry from Redis.
type SlowLogEntry struct {
	ID       int64
	Time     time.Time
	Duration time.Duration
	Args     []string
}

// GoRedisClient wraps go-redis/v9 and implements RedisClient.
type GoRedisClient struct {
	client *goredis.Client
}

// NewClient creates a new Redis client connection.
func NewClient(addr, password string, db int) (*GoRedisClient, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &GoRedisClient{client: client}, nil
}

func (c *GoRedisClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *GoRedisClient) Info(ctx context.Context, sections ...string) (string, error) {
	return c.client.Info(ctx, sections...).Result()
}

func (c *GoRedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.Scan(ctx, cursor, match, count).Result()
}

func (c *GoRedisClient) ObjectIdleTime(ctx context.Context, key string) (time.Duration, error) {
	return c.client.ObjectIdleTime(ctx, key).Result()
}

func (c *GoRedisClient) MemoryUsage(ctx context.Context, key string) (int64, error) {
	return c.client.MemoryUsage(ctx, key).Result()
}

func (c *GoRedisClient) SlowLogGet(ctx context.Context, num int64) ([]SlowLogEntry, error) {
	result, err := c.client.SlowLogGet(ctx, num).Result()
	if err != nil {
		return nil, err
	}
	entries := make([]SlowLogEntry, len(result))
	for i, r := range result {
		entries[i] = SlowLogEntry{
			ID:       r.ID,
			Time:     r.Time,
			Duration: r.Duration,
			Args:     r.Args,
		}
	}
	return entries, nil
}

func (c *GoRedisClient) ConfigGet(ctx context.Context, parameter string) (map[string]string, error) {
	return c.client.ConfigGet(ctx, parameter).Result()
}

func (c *GoRedisClient) DBSize(ctx context.Context) (int64, error) {
	return c.client.DBSize(ctx).Result()
}

func (c *GoRedisClient) Close() error {
	return c.client.Close()
}

// FormatBytes returns a human-readable byte size string.
func FormatBytes(bytes int64) string {
	const (
		kb = 1024
		mb = kb * 1024
		gb = mb * 1024
	)
	switch {
	case bytes >= gb:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(gb))
	case bytes >= mb:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(mb))
	case bytes >= kb:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(kb))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
