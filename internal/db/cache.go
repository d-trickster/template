package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"template/internal/config"
	"template/internal/logging"
	"time"

	"github.com/redis/go-redis/v9"
)

type CachedDB struct {
	database DB
	redis    *redis.Client
	ttl      time.Duration

	log *slog.Logger
}

func NewChachedDB(database DB, cfg *config.CacheConfig) (*CachedDB, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return &CachedDB{
		database: database,
		redis:    client,
		ttl:      cfg.TTL,
		log:      slog.Default().With(logging.ComponentAttr("cache")),
	}, nil
}

func (c *CachedDB) CreateUser(ctx context.Context, user *User) error {
	err := c.database.CreateUser(ctx, user)
	if err == nil {
		c.storeUser(ctx, user)
	}
	return err
}

func (c *CachedDB) GetUserByID(ctx context.Context, id string) (*User, error) {
	user, err := c.getUser(ctx, id)
	if err == nil {
		return user, nil
	}

	user, err = c.database.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.storeUser(ctx, user)

	return user, nil
}

func (c *CachedDB) getUser(ctx context.Context, id string) (*User, error) {
	data, err := c.redis.GetEx(ctx, "user"+id, c.ttl).Bytes()
	if err == nil {
		var user User
		if err = json.Unmarshal(data, &user); err != nil {
			c.log.Error("Failed to unmarshal cached user", logging.ErrAttr(err))
		} else {
			return &user, nil
		}
	} else if err != redis.Nil {
		c.log.Error("Redis GET failed", logging.ErrAttr(err))
	}
	return nil, err
}

func (c *CachedDB) storeUser(ctx context.Context, user *User) {
	key := "user" + user.ID
	data, err := json.Marshal(user)
	if err != nil {
		c.log.Error("Failed to marshal user", logging.ErrAttr(err))
	} else {
		if err := c.redis.Set(ctx, key, data, c.ttl).Err(); err != nil {
			c.log.Error("Failed to store user", logging.ErrAttr(err))
		}
	}
}
