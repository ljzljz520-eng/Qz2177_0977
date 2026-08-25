package workflow

import (
	"context"
	"fmt"

	"coursechain/domain"
)

func (s *Service) Notify(ctx context.Context, id, actor, detail string) error {
	if id == "" || actor == "" {
		return fmt.Errorf("record and actor are required")
	}
	if detail == "" {
		detail = "status notification"
	}
	return s.store.AppendEvent(ctx, domain.Event{ID: s.nextID("event"), RecordID: id, Kind: "notification", ActorID: actor, Detail: detail, CreatedAt: s.clock()})
}

func (s *Service) Track(ctx context.Context, id string) (domain.Tracking, error) {
	if id == "" {
		return domain.Tracking{}, fmt.Errorf("record id is required")
	}
	if _, err := s.store.GetRecord(ctx, id); err != nil {
		return domain.Tracking{}, err
	}
	events, err := s.store.ListEvents(ctx, id)
	if err != nil {
		return domain.Tracking{}, err
	}
	audits, err := s.store.ListAudits(ctx, id)
	if err != nil {
		return domain.Tracking{}, err
	}
	return domain.Tracking{RecordID: id, Events: events, Audits: audits}, nil
}

func (s *Service) Reconcile(ctx context.Context, id, actor string) (domain.Record, error) {
	record, err := s.store.GetRecord(ctx, id)
	if err != nil {
		return domain.Record{}, err
	}
	if record.Status != domain.StatusDelayed {
		return record, nil
	}
	return s.Review(ctx, domain.Review{RecordID: id, Reviewer: actor, Approved: true, Reason: "reconciled delayed status"})
}
