package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/tiskae/go-social/internal/store"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Hour

func (u *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cachedKey := fmt.Sprintf("user-%v", userID)

	data, err := u.rdb.Get(ctx, cachedKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}

	}

	return &user, nil
}

func (u *UserStore) Set(ctx context.Context, user *store.User) error {
	cachedKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return u.rdb.SetEX(ctx, cachedKey, json, UserExpTime).Err()
}
