package cache

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

func Test_Put(t *testing.T) {
	t.Parallel()

	t.Run("put item into cache", func(t *testing.T) {
		t.Parallel()
		cache := NewLRUCacheStore(10)

		cache.Put("key1", 1)
		cache.Put("key2", 2)
		cache.Put("key3", 3)

		cacheList := cache.List()

		require.True(t, cmp.Equal(
			[]any{3, 2, 1},
			cacheList,
		))
	})

	t.Run("put more items than capacity", func(t *testing.T) {
		t.Parallel()

		cache := NewLRUCacheStore(2)

		cache.Put("key1", 1)
		cache.Put("key2", 2)
		cache.Put("key3", 3)

		cacheList := cache.List()

		require.True(t, cmp.Equal(
			[]any{3, 2},
			cacheList,
		))
	})

	t.Run("put existing item updates its position", func(t *testing.T) {
		t.Parallel()

		cache := NewLRUCacheStore(3)

		cache.Put("key1", 1)
		cache.Put("key2", 2)
		cache.Put("key1", 3)

		cacheList := cache.List()

		require.Empty(t, cmp.Diff(
			[]any{3, 2},
			cacheList,
		))
	})
}

func Test_Get(t *testing.T) {
	t.Parallel()

	t.Run("get existing item from cache", func(t *testing.T) {
		cache := NewLRUCacheStore(10)

		cache.Put("key1", 1)

		got, exists := cache.Get("key1")

		require.True(t, exists)
		require.Equal(t, any(1), got)
	})

	t.Run("get existing items with two records", func(t *testing.T) {
		t.Parallel()

		cache := NewLRUCacheStore(10)

		cache.Put("key1", 1)
		cache.Put("key2", 2)

		got, exists := cache.Get("key1")
		require.True(t, exists)
		require.Equal(t, any(1), got)

		cacheList := cache.List()
		require.Empty(t, cmp.Diff(
			[]any{1, 2},
			cacheList,
		))
	})
}
