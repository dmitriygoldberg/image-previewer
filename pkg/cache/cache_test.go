package cache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic by queue size", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100) // [aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [bbb, aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300) // [ccc, bbb, aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 400) // purge aaa -> [ddd, ccc, bbb]
		require.False(t, wasInCache)

		val, ok := c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic by queue frequency", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100) // [aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200) // [bbb, aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300) // [ccc, bbb, aaa]
		require.False(t, wasInCache)

		wasInCache = c.Set("aaa", 500) // [aaa, ccc, bbb]
		require.True(t, wasInCache)

		val, ok := c.Get("ccc") // [ccc, aaa, bbb]
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("aaa") // [aaa, ccc, bbb]
		require.True(t, ok)
		require.Equal(t, 500, val)

		wasInCache = c.Set("ddd", 1000) // purge bbb -> // [ddd, aaa, ccc]
		require.False(t, wasInCache)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)

		wasInCache = c.Set("ppp", 0) // purge ccc -> // [ppp, ddd, aaa]
		require.False(t, wasInCache)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(5)

		c.Set("aaa", 5)
		c.Set("bbb", 7)

		val, ok := c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 7, val)

		c.Clear()

		_, ok = c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})
}
