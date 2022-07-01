package persist

import (
	"github.com/coffee-dns/coffee-dns/nameserver/record"
)

// Persist provides persistent storage
type Persist interface {
	// Initialize the namespace
	Init() error

	// Get a key's value
	Get(key string) (string, error)

	// Set a key's value
	Set(r record.Record) error

	// Delete a key
	Delete(key string) error

	// Close the connection
	Close() error

	// Dump all entries
	Dump() ([]byte, error)
}
