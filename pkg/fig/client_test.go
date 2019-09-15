package fig_test

import (
	"github.com/alicebob/miniredis"
	"github.com/rameshvk/fig/pkg/fig"
	"github.com/rameshvk/fig/pkg/server"

	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFigGetSinceEmpty(t *testing.T) {
	c, cleanup := startServer(t)
	defer cleanup()

	ver, config := c.GetSince(-1)
	if ver != -1 || len(config) != 0 {
		t.Error("Unexpected result", ver, config)
	}
}

func TestFigSet(t *testing.T) {
	c, cleanup := startServer(t)
	defer cleanup()

	c.Set("boo", `"hoo"`)
	ver, config := c.GetSince(-1)
	if ver != 1 || len(config) != 1 || config["boo"] != `"hoo"` {
		t.Error("Unexpected result", ver, config)
	}

	c.Set("boo", `"woop"`)
	ver, config = c.GetSince(ver)
	if ver != 2 || len(config) != 1 || config["boo"] != `"woop"` {
		t.Error("Unexpected result", ver, config)
	}
	ver, config = c.GetSince(-1)
	if ver != 2 || len(config) != 1 || config["boo"] != `"woop"` {
		t.Error("Unexpected result", ver, config)
	}
}

func TestFigSetBadRequest(t *testing.T) {
	c, cleanup := startServer(t)
	defer cleanup()

	mustPanic := func(cause string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("did not panic:", cause)
			}
		}()
		fn()
	}
	mustPanic("malformed json", func() {
		c.Set("boo", "hoo")
	})
	mustPanic("empty array", func() {
		c.Set("boo", "[]")
	})
	mustPanic("no objects", func() {
		c.Set("boo", "{}")
	})
}

func TestHistory(t *testing.T) {
	c, cleanup := startServer(t)
	defer cleanup()

	epoch, items := c.History("boo", "")
	if epoch != "" || len(items) != 0 {
		t.Fatal("unexpected", epoch, items)
	}

	c.Set("boo", `"hoo"`)
	c.Set("boo", `"hop"`)
	c.Set("boo", `"wop"`)

	epoch, items = c.History("boo", "")
	if epoch != "0" || !reflect.DeepEqual(items, []string{`"wop"`, `"hop"`, `"hoo"`}) {
		t.Fatal("unexpected", epoch, items)
	}
}

func startServer(t *testing.T) (*fig.Client, func()) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}

	store := server.NewRedisStore(s.Addr(), "test")
	ts := httptest.NewServer(server.Handler(func(r *http.Request) server.Store {
		return store
	}))
	return fig.New(ts.URL), func() { ts.Close(); s.Close() }
}
