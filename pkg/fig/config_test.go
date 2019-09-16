package fig_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/rameshvk/fig/pkg/fig"
	"github.com/rameshvk/fig/pkg/server"
)

func Example() {
	store, url, key, secret, cleanup := getStoreAndInfo()
	defer cleanup()

	cfg := fig.Config(url, key, secret, time.Second)

	// set my.setting to the expression:
	//   if (user == "boo") then "hoo" else "woo"
	// The <user> part needs to be provided when getting the config.

	expr := `["if", ["==", ["ref", "user"], "boo"], "hoo", "woo"]`
	store.Set("my.setting", expr)

	// now get the setting and provide user = boo as arg
	v, err := cfg.Get("my.setting", map[string]string{"user": "boo"})

	if v != "hoo" || err != nil {
		panic("unexpected result")
	}

}

func TestConfig(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	store := server.NewRedisStore(s.Addr(), "test-cfg")
	authStore := server.NewRedisStore(s.Addr(), "auth-store")
	unauthorized := func(r *http.Request) server.Store {
		return nil
	}
	authorized := func(r *http.Request) server.Store {
		return store
	}

	ts := httptest.NewServer(server.Handler(server.BasicAuth(authStore, authorized, unauthorized)))
	defer ts.Close()

	server.SetBasicAuthInfo(authStore, "mykey", "mysecret")
	store.Set("boo", `["ref", "boo"]`)

	cfg := fig.Config(ts.URL, "mykey", "mysecret", time.Millisecond)
	v, err := cfg.Get("boo", map[string]string{"boo": "hoo"})

	if v != "hoo" || err != nil {
		t.Fatal("Unexpected config", v, err)
	}
}

func getStoreAndInfo() (server.Store, string, string, string, func()) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	store := server.NewRedisStore(s.Addr(), "test-cfg")
	authStore := server.NewRedisStore(s.Addr(), "auth-store")
	unauthorized := func(r *http.Request) server.Store {
		return nil
	}
	authorized := func(r *http.Request) server.Store {
		return store
	}

	ts := httptest.NewServer(server.Handler(server.BasicAuth(authStore, authorized, unauthorized)))
	server.SetBasicAuthInfo(authStore, "mykey", "mysecret")
	cleanup := func() {
		ts.Close()
		s.Close()
	}

	return store, ts.URL, "mykey", "mysecret", cleanup
}
