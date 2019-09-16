// Package server implements the config server
//
// The server consists of an API handler and a store.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rameshvk/fig/pkg/eval"
)

// BasicAuth is the basic auth middleware that checks if a request
// is authorized by looking up the store for the key `auth:basic:user`.
// If allowed, it uses the authoried store, else the unauthorized store
func BasicAuth(s Store, authorized, unauthorized func(r *http.Request) Store) func(r *http.Request) Store {
	return func(r *http.Request) Store {
		user, pass, ok := r.BasicAuth()
		if !ok {
			return unauthorized(r)
		}
		_, configs := s.GetSince(-1)
		if v, ok := configs["auth:basic:"+user]; ok {
			scope := eval.ExtendScope(
				map[string]interface{}{
					"key":    user,
					"secret": pass,
					"api":    apiName(r),
				},
				eval.DefaultScope,
			)
			if v, err := eval.Encoded(v, scope); err == nil && v == true {
				return authorized(r)
			}
		}

		// not authorized
		return unauthorized(r)
	}
}

// SetBasicAuthInfo sets the basic auth password for the provided user
func SetBasicAuthInfo(s Store, user, password string) {
	encoded, err := json.Marshal(password)
	if err != nil {
		panic(err)
	}
	setting := fmt.Sprintf(`["==",["ref", "secret"],%s]`, encoded)
	s.Set("auth:basic:"+user, setting)
}
