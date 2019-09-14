// Package server implements the config server
//
// The server consists of an API handler and a store.
package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Store is the storage interface for the server.  See NewReids/Store
// for an emplementation
type Store interface {
	// GetSince returns all config changed since the last version
	//
	// Versions start from 1. Passing in a smaller version than
	// that would automatically fetch all config entries
	GetSince(version int) (newVersion int, configs map[string]string)

	// Set updates the connfig entry for the specific key
	Set(key string, val string)

	// History fetches changes for a specific key in reverse
	// chronological order.  The epoch can be used for
	// continuation with an empty string being the initial value
	History(key, epoch string) (newEpoch string, configs []string)
}

// Handler returns a HTTP handler for the config server service
func Handler(s Store) http.Handler {
	m := mux.NewRouter()

	m.Handle("/items", wrap(s, handleGetSince)).Methods("GET")
	m.Handle("/items/{key}", wrap(s, handleSet)).Methods("POST")
	m.Handle("/items/{key}", wrap(s, handleHistory)).Methods("GET")

	return m
}

func handleGetSince(s Store, w http.ResponseWriter, r *http.Request) interface{} {
	ver := -1
	if n, err := strconv.Atoi(r.URL.Query().Get("version")); err == nil {
		ver = n
	}
	ver, config := s.GetSince(ver)
	return map[string]interface{}{"version": ver, "config": config}
}

func handleSet(s Store, w http.ResponseWriter, r *http.Request) interface{} {
	var decoded interface{}
	if err := json.NewDecoder(r.Body).Decode(&decoded); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return map[string]interface{}{"error": err.Error()}
	}
	if err := isValid(decoded); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return map[string]interface{}{"error": err.Error()}
	}
	val, err := json.Marshal(decoded)
	if err != nil {
		panic(err)
	}
	s.Set(mux.Vars(r)["key"], string(val))
	return nil
}

func handleHistory(s Store, w http.ResponseWriter, r *http.Request) interface{} {
	epoch, history := s.History(mux.Vars(r)["key"], r.URL.Query().Get("epoch"))
	return map[string]interface{}{"epoch": epoch, "history": history}
}

func wrap(s Store, fn func(s Store, w http.ResponseWriter, r *http.Request) interface{}) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		if result := fn(s, w, r); result != nil {
			w.Header().Add("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(result)
			if err != nil {
				panic(err)
			}
		}
	}
	return handlers.CombinedLoggingHandler(
		os.Stdout,
		handlers.RecoveryHandler()(handlers.CORS()(http.HandlerFunc(f))),
	)
}

func isValid(v interface{}) error {
	switch v := v.(type) {
	case string, float64, bool:
		return nil
	case []interface{}:
		if len(v) == 0 {
			return errors.New("empty array not allowed")
		}
		for _, vv := range v {
			if err := isValid(vv); err != nil {
				return err
			}
		}
		return nil
	}

	return errors.New("unexpected type")
}
