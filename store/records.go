package store

import (
	"context"
	"fmt"
	"sort"

	"coursechain/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) PutRecord(ctx context.Context, record domain.Record) error {
	if err := ensureID(record.ID); err != nil {
		return err
	}
	if err := domain.ValidateStatus(record.Status); err != nil {
		return err
	}
	data, err := encode(record)
	if err != nil {
		return err
	}
	return s.Update(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Record"]).Put(keyFor(record.ID), data)
	})
}

func (s *Store) GetRecord(ctx context.Context, id string) (domain.Record, error) {
	if err := ensureID(id); err != nil {
		return domain.Record{}, err
	}
	var record domain.Record
	err := s.View(ctx, func(tx *bolt.Tx) error {
		value := tx.Bucket(bucketNames["Record"]).Get(keyFor(id))
		if value == nil {
			return ErrNotFound
		}
		return decode(value, &record)
	})
	return record, err
}

func (s *Store) DeleteRecord(ctx context.Context, id string) error {
	if err := ensureID(id); err != nil {
		return err
	}
	return s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Record"])
		if bucket.Get(keyFor(id)) == nil {
			return ErrNotFound
		}
		return bucket.Delete(keyFor(id))
	})
}

func (s *Store) ListRecords(ctx context.Context) ([]domain.Record, error) {
	items := make([]domain.Record, 0)
	err := s.View(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Record"]).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			items = append(items, record)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].SubmittedAt.Equal(items[j].SubmittedAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].SubmittedAt.Before(items[j].SubmittedAt)
	})
	return items, nil
}

func (s *Store) UpdateRecordStatus(ctx context.Context, id string, next domain.Status) (domain.Record, error) {
	if err := domain.ValidateStatus(next); err != nil {
		return domain.Record{}, err
	}
	var updated domain.Record
	err := s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Record"])
		value := bucket.Get(keyFor(id))
		if value == nil {
			return ErrNotFound
		}
		if err := decode(value, &updated); err != nil {
			return err
		}
		if err := domain.ValidateTransition(updated.Status, next); err != nil {
			return fmt.Errorf("status update: %w", err)
		}
		updated.Status = next
		updated.Revision++
		data, err := encode(updated)
		if err != nil {
			return err
		}
		return bucket.Put(keyFor(id), data)
	})
	return updated, err
}
