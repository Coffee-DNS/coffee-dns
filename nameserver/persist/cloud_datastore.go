package persist

import (
	"context"
	"fmt"

	"github.com/coffee-dns/coffee-dns/nameserver/record"

	"cloud.google.com/go/datastore"
)

// Datastore uses Google Cloud Datastore to persist dns records
type Datastore struct {
	// The Google Cloud project hosting the datastore
	ProjectID string

	// Cloud Datastore entity Kind
	Kind string

	client *datastore.Client
}

// Init configures the cloud datastore client
func (d *Datastore) Init() error {
	if d.ProjectID == "" {
		return fmt.Errorf("field must be set: ProjectID")
	}

	if d.Kind == "" {
		return fmt.Errorf("field must be set: Kind")
	}

	ctx := context.Background()
	c, err := datastore.NewClient(ctx, d.ProjectID)
	if err != nil {
		return fmt.Errorf("failed to configure cloud datastore in project %s: %s", d.ProjectID, err)
	}

	d.client = c
	return nil
}

// Get returns a value from the domain map
func (d *Datastore) Get(key string) (string, error) {
	ctx := context.Background()
	k := datastore.NameKey(d.Kind, key, nil)
	e := record.Record{}
	if err := d.client.Get(ctx, k, &e); err != nil {
		return "", fmt.Errorf("failed to get key %s from cloud datastore: %s", key, err)
	}
	return e.Value, nil
}

// Set adds a key value pair to the domain map
func (d *Datastore) Set(r record.Record) error {
	ctx := context.Background()
	k := datastore.NameKey(d.Kind, r.Hostname, nil)
	if _, err := d.client.Put(ctx, k, &r); err != nil {
		return fmt.Errorf("failed to create record %s with value %s: %s", r.Hostname, r.Value, err)
	}
	return nil
}

// Delete deletes a key from the domain map
func (d *Datastore) Delete(key string) error {
	ctx := context.Background()
	k := datastore.NameKey(d.Kind, key, nil)
	if err := d.client.Delete(ctx, k); err != nil {
		return fmt.Errorf("failed to delete key %s from cloud datastore: %s", key, err)
	}
	return nil
}

// Close always returns nil
func (d *Datastore) Close() error {
	return d.client.Close()
}

// Dump is not implemented
func (d *Datastore) Dump() ([]byte, error) {
	return nil, fmt.Errorf("datastore does not implement Dump()")
}
