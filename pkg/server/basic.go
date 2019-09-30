// Package server implements the config server
//
// The server consists of an API handler and a store.
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rameshvk/fig/pkg/fire"
	"github.com/rameshvk/fig/pkg/parse"
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
			parsed, errs := parse.String(v)
			if len(errs) > 0 {
				return unauthorized(r)
			}

			ctx := context.Background()
			scope := fire.Scope(
				ctx,
				fire.Globals(),
				[2]fire.Value{fire.String("key"), fire.String(user)},
				[2]fire.Value{fire.String("secret"), fire.String(pass)},
				[2]fire.Value{fire.String("api"), fire.String(apiName(r))},
			)
			result := fire.Eval(ctx, parsed, scope)
			if b, ok := result.Bool(ctx); b && ok {
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
	setting := fmt.Sprintf(`secret == %s`, encoded)
	s.Set("auth:basic:"+user, setting)
}
