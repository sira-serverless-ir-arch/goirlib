package cache

import (
	"github.com/sira-serverless-ir-arch/goirlib/cache"
	"testing"
	"time"
)

func TestRcu(t *testing.T) {
	rcu := cache.NewCacheRCU[int]()
	rcu.Put("a", 1)
	rcu.Put("b", 2)
	rcu.Put("c", 3)

	time.Sleep(5 * time.Second)

	v, _ := rcu.Get("a")
	if *v != 1 {
		t.Errorf("Expected 1, got %d", *v)
	}
	v, _ = rcu.Get("b")
	if *v != 2 {
		t.Errorf("Expected 2, got %d", *v)
	}
	v, _ = rcu.Get("c")
	if *v != 3 {
		t.Errorf("Expected 3, got %d", *v)
	}

	rcu.Put("b", 12)
	time.Sleep(5 * time.Second)

	v, _ = rcu.Get("b")
	if *v != 12 {
		t.Errorf("Expected 12, got %d", *v)
	}
}
