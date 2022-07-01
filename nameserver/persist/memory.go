package persist

import (
	"encoding/json"
	"fmt"

	"github.com/coffee-dns/coffee-dns/nameserver/record"
)

// Memory is a Memory type that satisfies the Persist interface
type Memory struct {
	domains map[string]record.Record
}

// Init creates the initial domain map
func (m *Memory) Init() error {
	m.domains = make(map[string]record.Record)
	return nil
}

// Get returns a value from the domain map
func (m *Memory) Get(key string) (string, error) {
	v, ok := m.domains[key]
	if !ok {
		return "", fmt.Errorf("no record for key %s", key)
	}
	return v.Value, nil
}

// Set adds a key value pair to the domain map
func (m *Memory) Set(r record.Record) error {
	m.domains[r.Hostname] = r
	return nil
}

// Delete deletes a key from the domain map
func (m *Memory) Delete(key string) error {
	delete(m.domains, key)
	return nil
}

// Close always returns nil
func (m *Memory) Close() error {
	return nil
}

// Dump returns all entires json encoded
func (m *Memory) Dump() ([]byte, error) {
	return json.Marshal(m.domains)
}
