package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/alicebob/miniredis"

	"github.com/rameshvk/fig/pkg/server"
)

var address = flag.String("http", ":80", "server:port for http listener")
var redis = flag.String("redis", "", "redis server:port")
var staticDir = flag.String("staticdir", "web", "directory for static html, js files")

func main() {
	flag.Parse()

	if *redis == "mini" {
		server, err := miniredis.Run()
		if err != nil {
			log.Fatal("could not start mini redis", err)
		}
		*redis = server.Addr()
	}

	store := server.NewRedisStore(*redis, "all")
	authorized := func(r *http.Request) server.Store {
		return store
	}
	unauthorized := func(r *http.Request) server.Store {
		return nil
	}

	authorize := server.BasicAuth(store, authorized, unauthorized)
	handler := server.Handler(authorize)

	http.Handle("/", http.FileServer(http.Dir(*staticDir)))
	http.Handle("/items", handler)
	http.Handle("/items/", handler)

	log.Fatal(http.ListenAndServe(*address, nil))
}
