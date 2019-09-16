// Package cache implements a simple cached store
package cache

import (
	"sync"
	"time"
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

// New wraps a store with a synchronized cache.
//
// The cache is updated atmost every refresh interval. Set() and History()
// are not cached. Only GetSince is cached
func New(s Store, refresh time.Duration, now func() time.Time) Store {
	n := time.Now
	if now != nil {
		n = now
	}
	return &cache{s, refresh, -1, nil, time.Time{}, n, sync.Mutex{}}
}

type cache struct {
	Store
	refresh       time.Duration
	ver           int
	config        map[string]string
	lastRefreshed time.Time
	now           func() time.Time
	sync.Mutex
}

func (c *cache) GetSince(version int) (newVersion int, configs map[string]string) {
	c.Lock()
	if version > 0 && version != c.ver {
		defer c.Unlock()
		return c.Store.GetSince(version)
	}
	defer c.Unlock()

	if c.lastRefreshed.Add(c.refresh).After(c.now()) {
		return c.ver, c.config
	}

	ver, next := c.Store.GetSince(c.ver)
	result := map[string]string{}
	for k, v := range c.config {
		result[k] = v
	}

	for k, v := range next {
		result[k] = v
	}

	c.config = result
	c.ver = ver
	c.lastRefreshed = c.now()
	return ver, result
}
