package cache

import (
	"os"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/sagikazarmark/slog-shim"
)

// simple cache implemented using ristretto cache library
type InMemoryCache struct {
	cache *ristretto.Cache
}

var (
	inMemoryCache *InMemoryCache
	inMemoryCacheInitOnce sync.Once
	credentialCache *InMemoryCache
	credentialCacheInitOnce sync.Once
	connectionCache *InMemoryCache
	connectionCacheInitOnce sync.Once
)

func initializeInMemoryCache(config *ristretto.Config) *InMemoryCache {
	if config == nil {
		config = &ristretto.Config{
			NumCounters: 100000,   // number of keys to track frequency
			MaxCost:     67108864, // maximum cost of cache (64mb).
			BufferItems: 64,       // number of keys per Get buffer.
		}
	}
	cache, err := ristretto.NewCache(config)
	if err != nil {
		slog.Error("error initializing in-memory cache", "error", err)
		os.Exit(1)
	}

	inMemoryCache = &InMemoryCache{cache}
	return inMemoryCache
}

func GetCache() *InMemoryCache {
	inMemoryCacheInitOnce.Do(func() {
		initializeInMemoryCache(nil)
	})
	return inMemoryCache
}

func initializeCredentialCache() {
	credCacheConfig := &ristretto.Config{
		NumCounters: 100000,   // number of keys to track frequency
		MaxCost:     67108864, // maximum cost of cache (64mb).
		BufferItems: 64,       // number of keys per Get buffer.
	}

	credCache, err := ristretto.NewCache(credCacheConfig)
	if err != nil {
		slog.Error("error initializing in-memory cache for credentials", "error", err)
		os.Exit(1)
	}

	credentialCache = &InMemoryCache{credCache}
}

func initializeConnectionCache() {
	connCacheConfig := &ristretto.Config{
		NumCounters: 100000,   // number of keys to track frequency
		MaxCost:     67108864, // maximum cost of cache (64mb).
		BufferItems: 64,       // number of keys per Get buffer.
	}

	connCache, err := ristretto.NewCache(connCacheConfig)
	if err != nil {
		slog.Error("error initializing in-memory cache for connections", "error", err)
		os.Exit(1)
	}

	connectionCache = &InMemoryCache{connCache}
}

func GetCredentialCache() *InMemoryCache {
	credentialCacheInitOnce.Do(func() {
		initializeCredentialCache()
	})
	return credentialCache
}

func GetConnectionCache() *InMemoryCache {
	connectionCacheInitOnce.Do(func() {
		initializeConnectionCache()
	})
	return connectionCache
}

func ResetCredentialCache() {
	credentialCache = nil
	initializeCredentialCache()
}

func ResetConnectionCache() {
	connectionCache = nil
	initializeConnectionCache()
}

func ResetInMemoryCache() {
	inMemoryCache = nil
	initializeInMemoryCache(nil)
}

func ResetAllCache() {
	ResetCredentialCache()
	ResetConnectionCache()
	ResetInMemoryCache()
}

func (cache *InMemoryCache) SetWithTTL(key string, value interface{}, ttl time.Duration) bool {
	res := cache.cache.SetWithTTL(key, value, 1, ttl)

	// wait for value to pass through buffers
	time.Sleep(10 * time.Millisecond)
	return res
}

func (cache *InMemoryCache) Get(key string) (interface{}, bool) {
	return cache.cache.Get(key)
}

func (cache *InMemoryCache) Delete(key string) {
	cache.cache.Del(key)
}
