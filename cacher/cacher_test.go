package cacher

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	mockGetterErrKey   = "getter_error_key"
	mockCacheGetErrKey = "cache_get_error_key"
	mockCacheSetErrKey = "cache_set_error_key"
)

type mockStruct struct {
	a int
}

type mockGetter struct {
	ret mockStruct
}

func (mg *mockGetter) Get(_ context.Context, key string) (mockStruct, error) {
	if key == mockGetterErrKey {
		return mockStruct{}, errors.New("getter error")
	}
	return mg.ret, nil
}

type mockCache struct {
	cache map[string]mockStruct
}

func (mc mockCache) Get(_ context.Context, key string) (mockStruct, error) {
	if key == mockCacheGetErrKey {
		return mockStruct{}, errors.New("cache Get error")
	}
	val, ok := mc.cache[key]
	if !ok {
		return val, ErrNoKey
	}
	return val, nil
}

func (mc mockCache) Set(_ context.Context, key string, value mockStruct) error {
	if key == mockCacheSetErrKey {
		return errors.New("cache Set error")
	}
	mc.cache[key] = value
	return nil
}

func TestCacherGet(t *testing.T) {
	ctx := context.Background()
	getter := &mockGetter{
		ret: mockStruct{a: 1},
	}
	cacher := NewCacher[string, mockStruct](mockCache{
		cache: make(map[string]mockStruct),
	}, getter.Get)

	t.Run("should return cached value even if it was changed in getter", func(t *testing.T) {
		val, err := cacher.Get(ctx, "key_1")
		require.NoError(t, err)
		assert.Equal(t, 1, val.a)

		getter.ret.a = 2
		val, err = cacher.Get(ctx, "key_2")
		require.NoError(t, err)
		assert.Equal(t, 2, val.a)

		val, err = cacher.Get(ctx, "key_1")
		require.NoError(t, err)
		assert.Equal(t, 1, val.a, "expected cached value")

		// Cleanup
		getter.ret.a = 1
	})

	t.Run("should return ErrGetProblem when cache Get errors", func(t *testing.T) {
		_, err := cacher.Get(ctx, mockCacheGetErrKey)
		assert.ErrorIs(t, err, ErrGetProblem)
	})

	t.Run("should return ErrSetProblem and cached data when cache Set errors", func(t *testing.T) {
		val, err := cacher.Get(ctx, mockCacheSetErrKey)
		require.ErrorIs(t, err, ErrSetProblem)
		assert.Equal(t, 1, val.a)
	})

	t.Run("should return error when getter errors", func(t *testing.T) {
		_, err := cacher.Get(ctx, mockGetterErrKey)
		assert.Error(t, err)
	})
}
