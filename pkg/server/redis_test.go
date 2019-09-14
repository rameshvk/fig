package server_test

import (
	"github.com/rameshvk/fig/pkg/server"
	"github.com/alicebob/miniredis"

	"reflect"
	"testing"
)


func TestRedisGetSinceEmpty(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	r := server.NewRedisStore(s.Addr(), "test")
	ver, config := r.GetSince(-1)
	if ver != -1 || len(config) != 0 {
		t.Error("Unexpected result", ver, config)
	}
}

func TestRedisSet(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	r := server.NewRedisStore(s.Addr(), "test")

	r.Set("boo", "hoo")
	ver, config := r.GetSince(-1)
	if ver != 1 || len(config) != 1 || config["boo"] != "hoo" {
		t.Error("Unexpected result", ver, config)
	}

	r.Set("boo", "woop")
 	ver, config = r.GetSince(ver)
	if ver != 2 || len(config) != 1 || config["boo"] != "woop" {
		t.Error("Unexpected result", ver, config)
	}	
 	ver, config = r.GetSince(-1)
	if ver != 2 || len(config) != 1 || config["boo"] != "woop" {
		t.Error("Unexpected result", ver, config)
	}	
}

func TestHistory(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal("mini redis failed", err)
	}
	defer s.Close()

	r := server.NewRedisStore(s.Addr(), "test")

	epoch, items := r.History("boo", "")
	if epoch != "" || len(items) != 0 {
		t.Fatal("unexpected", epoch, items)
	}

	r.Set("boo", "hoo")
	r.Set("boo", "hop")
	r.Set("boo", "wop")

	epoch, items = r.History("boo", "")
	if epoch != "0" || !reflect.DeepEqual(items, []string{"wop", "hop", "hoo"}) {
		t.Fatal("unexpected", epoch, items)
	}
}
