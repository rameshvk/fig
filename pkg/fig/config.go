// Package fig implements the Golang client for fig
package fig

import (
	"errors"
	"time"

	"github.com/rameshvk/fig/pkg/cache"
	"github.com/rameshvk/fig/pkg/eval"
)

// Config creates a new config getter which can be used to evaluate
// the configuration.
func Config(url, key, secret string, cacheFor time.Duration) Getter {
	return ConfigWithClient(New(url).WithKey(key, secret), cacheFor)
}

// ConfigWithClient returns a getter than can be used to efficiently
// access configuration entries.
func ConfigWithClient(c *Client, cacheFor time.Duration) Getter {
	s := cache.New(c, cacheFor, nil)
	return getter(func(key string, arg interface{}) (interface{}, error) {
		_, cfg := s.GetSince(-1)
		if entry, ok := cfg[key]; ok {
			base := eval.DefaultScope
			return eval.Encoded(entry, eval.Reflect(arg, base, nil))
		}
		return nil, ErrConfigNotFound
	})
}

// Getter allows fetching configuration entries
type Getter interface {
	Get(key string, arg interface{}) (interface{}, error)
}

// ErrConfigNotFound is returned by GetConfig if config is not found
var ErrConfigNotFound = errors.New("config not found")

type getter func(key string, arg interface{}) (interface{}, error)

func (g getter) Get(key string, arg interface{}) (interface{}, error) {
	return g(key, arg)
}
