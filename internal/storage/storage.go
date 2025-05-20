package storage

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrKeyNotFound is returned when a key is not found in storage
	ErrKeyNotFound = errors.New("key not found")
	
	// ErrStorageUnavailable is returned when the storage backend is unavailable
	ErrStorageUnavailable = errors.New("storage backend unavailable")
)

// Storage defines the interface for rate limit data storage
type Storage interface {
	// Increment atomically increments the counter for the given key by the given amount
	// and returns the new value and any error encountered
	Increment(ctx context.Context, key string, value int64, expiry time.Duration) (int64, error)
	
	// Get retrieves the current value for the given key
	// Returns ErrKeyNotFound if the key does not exist
	Get(ctx context.Context, key string) (int64, error)
	
	// Set sets the value for the given key with an expiration time
	Set(ctx context.Context, key string, value int64, expiry time.Duration) error
	
	// Delete removes the given key from storage
	Delete(ctx context.Context, key string) error
	
	// Close closes the storage backend and releases any resources
	Close() error
}
