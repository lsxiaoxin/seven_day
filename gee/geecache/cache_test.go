package geecache

import (
	"testing"
)

func TestCacheAddAndGet(t *testing.T) {
	c := &cache{cacheBytes: int64(100)}

	c.add("key1", NewByteView("hello"))
	c.add("key2", NewByteView("world"))

	if v, ok := c.get("key1"); !ok || v.String() != "hello" {
		t.Fatalf("expected get(key1)=hello, but got %s", v.String())
	}

	if v, ok := c.get("key2"); !ok || v.String() != "world" {
		t.Fatalf("expected get(key2)=world, but got %s", v.String())
	}

	if _, ok := c.get("key3"); ok {
		t.Fatalf("expected miss on key3, but hit")
	}
}

// 测试自动初始化 LRU
func TestLazyInit(t *testing.T) {
	c := &cache{cacheBytes: int64(10)} // 初始 lru=nil

	// 调用 add 时自动初始化 lru
	c.add("a", NewByteView("A"))
	if c.lru == nil {
		t.Fatal("expected lru to be initialized after first add")
	}
}

// 测试并发安全（简单验证不会 panic）
func TestConcurrentAccess(t *testing.T) {
	c := &cache{cacheBytes: int64(50)}

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := "k" + string(rune('a'+i))
			c.add(key, NewByteView("val"))
			c.get(key)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}