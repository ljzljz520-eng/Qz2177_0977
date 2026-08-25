package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	ErrNotFound    = errors.New("entity not found")
	ErrStoreClosed = errors.New("store is closed")
)

var bucketNames = map[string][]byte{
	"Record": []byte("records"),
	"User":   []byte("users"),
	"Event":  []byte("events"),
	"Audit":  []byte("audits"),
}

type Store struct {
	mu     sync.RWMutex
	db     *bolt.DB
	path   string
	closed bool
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}
	db, err := bolt.Open(filepath.Clean(path), 0o600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bbolt: %w", err)
	}
	s := &Store{db: db, path: filepath.Clean(path)}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.path
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}

func (s *Store) IsClosed() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.closed
}

func (s *Store) checkOpen() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed || s.db == nil {
		return ErrStoreClosed
	}
	return nil
}

func (s *Store) View(ctx context.Context, fn func(*bolt.Tx) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.View(func(tx *bolt.Tx) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return fn(tx)
	})
}

func (s *Store) Update(ctx context.Context, fn func(*bolt.Tx) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return fn(tx)
	})
}

func encode(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("encode entity: %w", err)
	}
	return data, nil
}

func decode(data []byte, target any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode entity: %w", err)
	}
	return nil
}

func keyFor(id string) []byte {
	return []byte(id)
}

func ensureID(id string) error {
	if id == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}
