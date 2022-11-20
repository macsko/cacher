package cacher

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrNoKey should be returned by Cache if the key is not present.
	// Then the entry will be pulled from the getter and set in the Cache.
	ErrNoKey = errors.New("key is not present in the cache")
	// ErrGetProblem is returned by the Cacher's Get method if the Cache returned not-ErrNoKey error from Get method.
	ErrGetProblem = errors.New("cannot get item from cache")
	// ErrSetProblem is returned by the Cacher's Get method if the Cache returned the error from Set method.
	ErrSetProblem = errors.New("cannot set item in cache")
)

// Getter is a function to be wrapped by the Cacher
type Getter[K comparable, T any] func(ctx context.Context, id K) (T, error)

// Cache is abstract Get/Set interface.
// Appropriate implementation should use ErrNoKey in the Get method.
type Cache[K comparable, T any] interface {
	Get(ctx context.Context, key K) (T, error)
	Set(ctx context.Context, key K, value T) error
}

type Cacher[K comparable, T any] struct {
	cache Cache[K, T]
	f     Getter[K, T]
}

func NewCacher[K comparable, T any](c Cache[K, T], f Getter[K, T]) Cacher[K, T] {
	return Cacher[K, T]{
		cache: c,
		f:     f,
	}
}

// Get caches and returns value received from calling Cacher's function f with context and key argument.
// Even if Cache's Set method fails, the appropriate value should be returned with ErrSetProblem indicated.
func (c Cacher[K, T]) Get(ctx context.Context, key K) (T, error) {
	value, err := c.cache.Get(ctx, key)
	if err != nil && !errors.Is(err, ErrNoKey) {
		return value, fmt.Errorf("%w: %s", ErrGetProblem, err.Error())
	}
	if err == nil {
		return value, nil
	}

	value, err = c.f(ctx, key)
	if err != nil {
		return value, err
	}

	err = c.cache.Set(ctx, key, value)
	if err != nil {
		return value, fmt.Errorf("%w: %s", ErrSetProblem, err.Error())
	}

	return value, nil
}
