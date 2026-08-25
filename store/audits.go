package store

import (
	"context"
	"fmt"
	"sort"

	"coursechain/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) AppendAudit(ctx context.Context, audit domain.Audit) error {
	if err := ensureID(audit.ID); err != nil {
		return err
	}
	if audit.RecordID == "" || audit.Action == "" {
		return fmt.Errorf("record id and action are required")
	}
	data, err := encode(audit)
	if err != nil {
		return err
	}
	return s.Update(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Audit"]).Put(keyFor(audit.ID), data)
	})
}

func (s *Store) ListAudits(ctx context.Context, recordID string) ([]domain.Audit, error) {
	items := make([]domain.Audit, 0)
	err := s.View(ctx, func(tx *bolt.Tx) error {
		return tx.Bucket(bucketNames["Audit"]).ForEach(func(_, value []byte) error {
			var audit domain.Audit
			if err := decode(value, &audit); err != nil {
				return err
			}
			if recordID == "" || audit.RecordID == recordID {
				items = append(items, audit)
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

func (s *Store) LatestAudit(ctx context.Context, recordID string) (domain.Audit, error) {
	items, err := s.ListAudits(ctx, recordID)
	if err != nil {
		return domain.Audit{}, err
	}
	if len(items) == 0 {
		return domain.Audit{}, ErrNotFound
	}
	return items[len(items)-1], nil
}
