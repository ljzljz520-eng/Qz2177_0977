package store

import (
	"context"
	"fmt"
	"time"

	"coursechain/domain"
	bolt "go.etcd.io/bbolt"
)

type MaintenanceResult struct {
	Scanned      int                   `json:"scanned"`
	Visible      int                   `json:"visible"`
	Expired      int                   `json:"expired"`
	Repaired     int                   `json:"repaired"`
	StatusCounts map[domain.Status]int `json:"status_counts"`
	CompletedAt  time.Time             `json:"completed_at"`
}

func (s *Store) Maintenance(ctx context.Context, now time.Time, maxAge time.Duration) (MaintenanceResult, error) {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if maxAge < 0 {
		return MaintenanceResult{}, fmt.Errorf("max age cannot be negative")
	}
	records, err := s.ListRecords(ctx)
	if err != nil {
		return MaintenanceResult{}, err
	}
	result := MaintenanceResult{StatusCounts: make(map[domain.Status]int), CompletedAt: now}
	for _, record := range records {
		result.Scanned++
		result.StatusCounts[record.Status]++
		if domain.IsVisible(record.Status) {
			result.Visible++
		}
		if maxAge > 0 && now.Sub(record.UpdatedAt) > maxAge && !record.IsTerminal() {
			result.Expired++
			if _, updateErr := s.UpdateRecordStatus(ctx, record.ID, domain.StatusDelayed); updateErr == nil {
				result.Repaired++
			}
		}
	}
	return result, nil
}

func (s *Store) VerifyIntegrity(ctx context.Context) error {
	records, err := s.ListRecords(ctx)
	if err != nil {
		return err
	}
	if err := domain.ValidateRecordSet(records); err != nil {
		return err
	}
	for _, record := range records {
		events, eventErr := s.ListEvents(ctx, record.ID)
		if eventErr != nil {
			return eventErr
		}
		for _, event := range events {
			if event.RecordID != record.ID {
				return fmt.Errorf("event %s points to %s", event.ID, event.RecordID)
			}
		}
		audits, auditErr := s.ListAudits(ctx, record.ID)
		if auditErr != nil {
			return auditErr
		}
		for _, audit := range audits {
			if audit.RecordID != record.ID {
				return fmt.Errorf("audit %s points to %s", audit.ID, audit.RecordID)
			}
		}
	}
	return nil
}

func (s *Store) DeleteEventsForRecord(ctx context.Context, recordID string) (int, error) {
	if recordID == "" {
		return 0, fmt.Errorf("record id is required")
	}
	removed := 0
	err := s.Update(ctx, func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketNames["Event"])
		keys := make([][]byte, 0)
		if err := bucket.ForEach(func(key, value []byte) error {
			var event domain.Event
			if err := decode(value, &event); err != nil {
				return err
			}
			if event.RecordID == recordID {
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
	if err != nil {
		return removed, err
	}
	return removed, nil
}
