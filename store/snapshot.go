package store

import (
	"context"
	"fmt"
	"time"

	"coursechain/domain"
)

type Snapshot struct {
	Records []domain.Record
	Users   []domain.User
	Events  []domain.Event
	Audits  []domain.Audit
	At      time.Time
}

func (s *Store) Snapshot(ctx context.Context) (Snapshot, error) {
	records, err := s.ListRecords(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	users, err := s.ListUsers(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	events, err := s.ListEvents(ctx, "")
	if err != nil {
		return Snapshot{}, err
	}
	audits, err := s.ListAudits(ctx, "")
	if err != nil {
		return Snapshot{}, err
	}
	return Snapshot{Records: records, Users: users, Events: events, Audits: audits, At: time.Now().UTC()}, nil
}

func (s Snapshot) Validate() error {
	if s.At.IsZero() {
		return fmt.Errorf("snapshot timestamp is required")
	}
	for _, record := range s.Records {
		if err := domain.ValidateStatus(record.Status); err != nil {
			return err
		}
	}
	return nil
}
