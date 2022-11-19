package cacher

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNoKey      = errors.New("key is not present in the cache")
	ErrGetProblem = errors.New("cannot get item from cache")
	ErrSetProblem = errors.New("cannot set item in cache")
)

type Getter[K comparable, T any] func(ctx context.Context, id K) (T, error)

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
