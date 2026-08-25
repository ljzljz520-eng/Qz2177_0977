package workflow

import (
	"context"
	"fmt"

	"coursechain/domain"
)

func (s *Service) Review(ctx context.Context, review domain.Review) (domain.Record, error) {
	if review.RecordID == "" || review.Reviewer == "" {
		return domain.Record{}, fmt.Errorf("record and reviewer are required")
	}
	if _, err := s.EnsureActor(ctx, review.Reviewer); err != nil {
		return domain.Record{}, err
	}
	record, err := s.store.GetRecord(ctx, review.RecordID)
	if err != nil {
		return domain.Record{}, err
	}
	next := domain.StatusRejected
	if review.Approved {
		next = domain.StatusImmediate
	}
	if record.Status == domain.StatusDelayed && review.Approved {
		next = domain.StatusImmediate
	}
	change, err := domain.NewStatusChange(record, next, review.Reviewer, review.Reason, s.clock())
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := s.store.UpdateRecordStatus(ctx, record.ID, change.After)
	if err != nil {
		return domain.Record{}, err
	}
	audit := domain.Audit{ID: s.nextID("audit"), RecordID: record.ID, Action: "review", ActorID: review.Reviewer, Before: change.Before, After: change.After, Reason: review.Reason, CreatedAt: s.clock()}
	if err := s.store.AppendAudit(ctx, audit); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) Archive(ctx context.Context, id, actor string) (domain.Record, error) {
	if id == "" || actor == "" {
		return domain.Record{}, fmt.Errorf("record and actor are required")
	}
	if _, err := s.EnsureActor(ctx, actor); err != nil {
		return domain.Record{}, err
	}
	record, err := s.store.GetRecord(ctx, id)
	if err != nil {
		return domain.Record{}, err
	}
	change, err := domain.NewStatusChange(record, domain.StatusArchived, actor, "archive requested", s.clock())
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := s.store.UpdateRecordStatus(ctx, id, change.After)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.store.AppendAudit(ctx, domain.Audit{ID: s.nextID("audit"), RecordID: id, Action: "archive", ActorID: actor, Before: change.Before, After: change.After, Reason: change.Reason, CreatedAt: s.clock()}); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) Reject(ctx context.Context, id, actor, reason string) (domain.Record, error) {
	record, err := s.store.GetRecord(ctx, id)
	if err != nil {
		return domain.Record{}, err
	}
	change, err := domain.NewStatusChange(record, domain.StatusRejected, actor, reason, s.clock())
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := s.store.UpdateRecordStatus(ctx, id, change.After)
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.store.AppendAudit(ctx, domain.Audit{ID: s.nextID("audit"), RecordID: id, Action: "reject", ActorID: actor, Before: change.Before, After: change.After, Reason: reason, CreatedAt: s.clock()}); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}
