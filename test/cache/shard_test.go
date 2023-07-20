package cache

import (
	"github.com/sira-serverless-ir-arch/goirlib/cache"
	"testing"
)

func TestShard(t *testing.T) {
	sh := cache.NewShardMap[int](3)
	iMap, _ := sh.Get("a")
	iMap.Set("as", 1)
	iMap, _ = sh.Get("b")
	iMap.Set("bs", 2)
	iMap, _ = sh.Get("c")
	iMap.Set("cs", 3)

	//time.Sleep(5 * time.Second)

	iMap, _ = sh.Get("a")
	v, _ := iMap.Get("as")
	if *v != 1 {
		t.Errorf("Expected 1, got %d", v)
	}

	iMap, _ = sh.Get("b")
	v, _ = iMap.Get("bs")
	if *v != 2 {
		t.Errorf("Expected 2, got %d", v)
	}

	iMap, _ = sh.Get("c")
	v, _ = iMap.Get("cs")
	if *v != 3 {
		t.Errorf("Expected 3, got %d", v)
	}

	iMap, _ = sh.Get("d")
	iMap.Set("ds", 4)

	//time.Sleep(5 * time.Second)

	iMap, _ = sh.Get("d")
	v, _ = iMap.Get("ds")
	if *v != 4 {
		t.Errorf("Expected 4, got %d", v)
	}

}
