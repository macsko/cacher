package memcached

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/macsko/cacher/cacher"

	"github.com/rainycape/memcache"
)

// Memcached is a wrapper for rainycape/memcache library Client and key expiration.
// The provided Go's structures of type T are stored as the JSON.
type Memcached[T any] struct {
	client *memcache.Client
	// Expiration is the cache expiration time, in seconds: either a relative
	// time from now (up to 1 month), or an absolute Unix epoch time.
	// Zero means the Item has no expiration time.
	// Description is taken from memcache.Item structure.
	expiration int32
}

func NewMemcached[T any](c *memcache.Client, expiration int32) Memcached[T] {
	return Memcached[T]{
		client:     c,
		expiration: expiration,
	}
}

func (m Memcached[T]) Get(_ context.Context, key string) (T, error) {
	var value T

	item, err := m.client.Get(key)
	if errors.Is(err, memcache.ErrCacheMiss) {
		return value, cacher.ErrNoKey
	}
	if err != nil {
		return value, err
	}

	err = json.Unmarshal(item.Value, &value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (m Memcached[T]) Set(_ context.Context, key string, value T) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return m.client.Add(&memcache.Item{
		Key:        key,
		Value:      v,
		Expiration: m.expiration,
	})
}
