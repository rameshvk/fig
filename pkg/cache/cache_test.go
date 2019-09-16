package cache_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/rameshvk/fig/pkg/cache"
	"github.com/rameshvk/fig/pkg/server"
)

func TestCache(t *testing.T) {
	redis, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer redis.Close()

	s := server.NewRedisStore(redis.Addr(), "test-redis")
	fakeTime := time.Now()

	now := func() time.Time { return fakeTime }
	duration := 5 * time.Second

	s = cache.New(s, duration, now)

	ver, config := s.GetSince(-1)
	if ver != -1 || len(config) > 0 {
		t.Fatal("Unexpected config change", ver, config)
	}
	s.Set("boo", `"hoo"`)
	ver, config = s.GetSince(-1)
	if ver != -1 || len(config) > 0 {
		t.Fatal("Unexpected config change", ver, config)
	}

	// update time and redo test
	fakeTime = fakeTime.Add(duration)
	ver, config = s.GetSince(-1)
	if ver != 1 || len(config) != 1 || config["boo"] != `"hoo"` {
		t.Fatal("Unexpected config change", ver, config)
	}

	// update once more
	s.Set("boo", `"woo"`)
	fakeTime = fakeTime.Add(duration)
	ver, config = s.GetSince(-1)
	if ver != 2 || len(config) != 1 || config["boo"] != `"woo"` {
		t.Fatal("Unexpected config change", ver, config)
	}

	if x, cfg := s.GetSince(1); x != 2 || len(cfg) != 1 || cfg["boo"] != `"woo"` {
		t.Fatal("Unexpected pass throgh", x, cfg)
	}
}
