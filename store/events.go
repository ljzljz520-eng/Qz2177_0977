package store

import (
	"context"
	"fmt"
	"sort"

	"coursechain/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) AppendEvent(ctx context.Context, event domain.Event) error {
	if err := ensureID(event.ID); err != nil {
		return err
	}
	if event.RecordID == "" || event.Kind == "" {
		return fmt.Errorf("record id and kind are required")
	}
	data, err := encode(event)
	if err != nil {
		return err
	}
	return s.Update(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Event"]).Put(keyFor(event.ID), data)
	})
}

func (s *Store) ListEvents(ctx context.Context, recordID string) ([]domain.Event, error) {
	items := make([]domain.Event, 0)
	err := s.View(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Event"]).ForEach(func(_, value []byte) error {
			var event domain.Event
			if err := decode(value, &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				items = append(items, event)
			}
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return items, nil
}

func (s *Store) CountEvents(ctx context.Context, recordID string) (int, error) {
	count := 0
	err := s.View(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Event"]).ForEach(func(_, value []byte) error {
			var event domain.Event
			if err := decode(value, &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				count++
			}
			return nil
		})
	})
	return count, err
}
