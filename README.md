# Cacher
Simple Go caching library with built-in Memcached support. It wraps a getter function and caches its results.

## Installing
```bash
go get github.com/macsko/cacher
```

## Memcached support
By default, support for Memcached caching method is added. To use it, it is required to [download](https://memcached.org) and run the Memcached.
Cacher uses [rainycape/memcache](https://github.com/rainycape/memcache) library as a client.

## Usage
- Create new Cacher using a Memcached client with 10 seconds expiration and `getter` function
```go
mc, err := memcache.New("10.0.0.1:11211")
...
c := cacher.NewCacher[string, int](
	memcached.NewMemcached[int](mc, 10),
	getter,
)
```
- Get value
```go
value, err := c.Get(ctx, "key")
```

## Using other caching mechanisms
Core (Cacher) code of the library is open to use with other caching mechanisms by passing a structure fulfilling the `Cache` interface:
```go
type Cache[K comparable, T any] interface {
	Get(ctx context.Context, key K) (T, error)
	Set(ctx context.Context, key K, value T) error
}
```
**Note:** ErrNoKey should be returned by Get method if the key is not present.
