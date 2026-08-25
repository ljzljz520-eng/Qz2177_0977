package store

import (
	"context"
	"fmt"
	"sort"

	"coursechain/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) PutRecords(ctx context.Context, records []domain.Record) error {
	if err := domain.ValidateRecordSet(records); err != nil {
		return err
	}
	encoded := make([][]byte, len(records))
	for index, record := range records {
		data, err := encode(record)
		if err != nil {
			return err
		}
		encoded[index] = data
	}
	return s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Record"])
		for index, record := range records {
			if err := bucket.Put(keyFor(record.ID), encoded[index]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) PutEvents(ctx context.Context, events []domain.Event) error {
	return s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Event"])
		for _, event := range events {
			if err := ensureID(event.ID); err != nil || event.RecordID == "" {
				return fmt.Errorf("invalid event %s", event.ID)
			}
			data, err := encode(event)
			if err != nil {
				return err
			}
			if err := bucket.Put(keyFor(event.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) PutAudits(ctx context.Context, audits []domain.Audit) error {
	return s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Audit"])
		for _, audit := range audits {
			if err := ensureID(audit.ID); err != nil || audit.RecordID == "" {
				return fmt.Errorf("invalid audit %s", audit.ID)
			}
			data, err := encode(audit)
			if err != nil {
				return err
			}
			if err := bucket.Put(keyFor(audit.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) PurgeRejected(ctx context.Context) (int, error) {
	removed := 0
	err := s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Record"])
		keys := make([][]byte, 0)
		if err := bucket.ForEach(func(key, value []byte) error {
			var record domain.Record
			if err := decode(value, &record); err != nil {
				return err
			}
			if record.Status == domain.StatusRejected {
				keys = append(keys, append([]byte(nil), key...))
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
			removed++
		}
		return nil
	})
	return removed, err
}

func (s *Store) StatusCounts(ctx context.Context) (map[domain.Status]int, error) {
	counts := make(map[domain.Status]int)
	items, err := s.ListRecords(ctx)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		counts[item.Status]++
	}
	return counts, nil
}

func SortRecordsByRevision(records []domain.Record) []domain.Record {
	result := append([]domain.Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Revision == result[j].Revision {
			return result[i].ID < result[j].ID
		}
		return result[i].Revision > result[j].Revision
	})
	return result
}
