package server_test

import (
	"github.com/alicebob/miniredis"
	"github.com/rameshvk/fig/pkg/fig"
	"github.com/rameshvk/fig/pkg/server"

	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestRedis(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	suite := Suite{server.NewRedisStore(s.Addr(), "test-redis")}
	suite.Run(t)
}

func TestAuthorizedHandler(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	store := server.NewRedisStore(s.Addr(), "test-handler")
	authStore := server.NewRedisStore(s.Addr(), "auth-store")
	unauthorized := func(r *http.Request) server.Store {
		return nil
	}
	authorized := func(r *http.Request) server.Store {
		return store
	}

	ts := httptest.NewServer(server.Handler(server.BasicAuth(authStore, authorized, unauthorized)))
	defer ts.Close()

	server.SetBasicAuthInfo(authStore, "authorized_key", "secret")
	suite := Suite{fig.New(ts.URL).WithKey("authorized_key", "secret")}
	suite.Run(t)
	t.Run("MalformedJSON", suite.testMalformedJSON)
}

func TestUnauthorizedHandler(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	store := server.NewRedisStore(s.Addr(), "test-handler")
	unauthorized := func(r *http.Request) server.Store {
		return nil
	}
	authorized := func(r *http.Request) server.Store {
		return store
	}

	ts := httptest.NewServer(server.Handler(server.BasicAuth(store, authorized, unauthorized)))
	defer ts.Close()

	c := fig.New(ts.URL)

	mustPanic := func(cause string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("did not panic:", cause)
			}
		}()
		fn()
	}

	mustPanic("GetSince", func() {
		c.GetSince(-1)
	})
	mustPanic("Set", func() {
		c.Set("boo", `"true"`)
	})
	mustPanic("History", func() {
		c.History("boo", "")
	})
}

func TestUnauthorizedHandlerWrongKey(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	store := server.NewRedisStore(s.Addr(), "test-handler")
	authStore := server.NewRedisStore(s.Addr(), "auth-store")
	unauthorized := func(r *http.Request) server.Store {
		return nil
	}
	authorized := func(r *http.Request) server.Store {
		return store
	}

	ts := httptest.NewServer(server.Handler(server.BasicAuth(authStore, authorized, unauthorized)))
	defer ts.Close()

	server.SetBasicAuthInfo(authStore, "authorized_key", "secret")

	c := fig.New(ts.URL).WithKey("authorized_key", "wrong")

	mustPanic := func(cause string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("did not panic:", cause)
			}
		}()
		fn()
	}

	mustPanic("GetSince", func() {
		c.GetSince(-1)
	})
	mustPanic("Set", func() {
		c.Set("boo", `"true"`)
	})
	mustPanic("History", func() {
		c.History("boo", "")
	})
}

type Suite struct {
	server.Store
}

func (s Suite) Run(t *testing.T) {
	t.Run("GetSinceEmpty", s.testGetSinceEmpty)
	t.Run("Set", s.testSet)
	t.Run("History", s.testHistory)
}

func (s Suite) testGetSinceEmpty(t *testing.T) {
	ver, config := s.GetSince(-1)
	if ver != -1 || len(config) != 0 {
		t.Error("Unexpected result", ver, config)
	}
}

func (s Suite) testSet(t *testing.T) {
	s.Set("boo", `"hoo"`)
	ver, config := s.GetSince(-1)
	if ver != 1 || len(config) != 1 || config["boo"] != `"hoo"` {
		t.Error("Unexpected result", ver, config)
	}

	s.Set("boo", `"woop"`)
	ver, config = s.GetSince(ver)
	if ver != 2 || len(config) != 1 || config["boo"] != `"woop"` {
		t.Error("Unexpected result", ver, config)
	}
	ver, config = s.GetSince(-1)
	if ver != 2 || len(config) != 1 || config["boo"] != `"woop"` {
		t.Error("Unexpected result", ver, config)
	}
}

func (s Suite) testHistory(t *testing.T) {
	ver, _ := s.GetSince(-1)
	epoch, items := s.History("boop", "")
	if epoch != "" || len(items) != 0 {
		t.Fatal("unexpected", epoch, items)
	}

	s.Set("boop", `"hoo"`)
	s.Set("boop", `"hop"`)
	s.Set("boop", `"wop"`)

	epoch, items = s.History("boop", "")
	if epoch != strconv.Itoa(ver) || !reflect.DeepEqual(items, []string{`"wop"`, `"hop"`, `"hoo"`}) {
		t.Fatal("unexpected", epoch, items)
	}
}

func (s Suite) testMalformedJSON(t *testing.T) {
	mustPanic := func(cause string, fn func()) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("did not panic:", cause)
			}
		}()
		fn()
	}
	mustPanic("malformed json", func() {
		s.Set("boo", "hoo")
	})
	mustPanic("empty array", func() {
		s.Set("boo", "[]")
	})
	mustPanic("no objects", func() {
		s.Set("boo", "{}")
	})
}
