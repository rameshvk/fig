// Package server implements the config server
//
// The server consists of an API handler and a store.
package server

// Store is the storage interface for the server.  See NewReids/Store
// for an emplementation
type Store interface {
	// GetSince returns all config changed since the last version
	//
	// Versions start from 1. Passing in a smaller version than
	// that would automatically fetch all config entries
	GetSince(version int)  (newVersion int, configs map[string]string)

	// Set updates the connfig entry for the specific key
	Set(key string, val string)

	// History fetches changes for a specific key in reverse
	// chronological order.  The epoch can be used for
	// continuation with an empty string being the initial value
	History(key, epoch string) (newEpoch string, configs []string)
}
