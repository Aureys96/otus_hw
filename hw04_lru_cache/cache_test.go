package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
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

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		_, isInCache := c.Get("aaa")
		require.True(t, isInCache)
		_, isInCache = c.Get("bbb")
		require.True(t, isInCache)
		_, isInCache = c.Get("ccc")
		require.True(t, isInCache)

		c.Clear()

		_, isInCache = c.Get("aaa")
		require.False(t, isInCache)
		_, isInCache = c.Get("bbb")
		require.False(t, isInCache)
		_, isInCache = c.Get("ccc")
		require.False(t, isInCache)

	})

	t.Run("queue capacity exceeded test", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		_, isInCache := c.Get("ccc")
		require.True(t, isInCache)
		_, isInCache = c.Get("bbb")
		require.True(t, isInCache)
		_, isInCache = c.Get("ddd")
		require.True(t, isInCache)

		_, isInCache = c.Get("aaa")
		require.False(t, isInCache)

	})

	t.Run("least frequent key replaced", func(t *testing.T) {
		c := NewCache(4)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)
		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		_, isInCache := c.Get("ccc")
		require.True(t, isInCache)
		wasInCache = c.Set("bbb", 201)
		require.True(t, wasInCache)
		val, isInCache := c.Get("bbb")
		require.True(t, isInCache)
		require.Equal(t, 201, val)
		_, isInCache = c.Get("ddd")
		require.True(t, isInCache)
		wasInCache = c.Set("bbb", 202)
		require.True(t, wasInCache)
		val, isInCache = c.Get("bbb")
		require.True(t, isInCache)
		require.Equal(t, 202, val)
		wasInCache = c.Set("eee", 500)
		require.False(t, wasInCache)

		val, isInCache = c.Get("eee")
		require.True(t, isInCache)
		require.Equal(t, 500, val)
		_, isInCache = c.Get("aaa")
		require.False(t, isInCache)

	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
