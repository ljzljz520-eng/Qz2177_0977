package workflow

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"coursechain/domain"
	"coursechain/linkedlist"
	"coursechain/store"
)

type Service struct {
	store  *store.Store
	clock  func() time.Time
	seq    atomic.Uint64
	course string
}

func NewService(s *store.Store) *Service {
	return &Service{store: s, clock: func() time.Time { return time.Now().UTC() }, course: "course10"}
}

func (s *Service) SetClock(clock func() time.Time) {
	if clock != nil {
		s.clock = clock
	}
}

func (s *Service) Store() *store.Store {
	return s.store
}

func (s *Service) nextID(prefix string) string {
	value := s.seq.Add(1)
	return fmt.Sprintf("%s-%d-%d", prefix, s.clock().UnixNano(), value)
}

func (s *Service) Submit(ctx context.Context, input domain.Submission, actor string) (domain.Record, error) {
	input = domain.NormalizeSubmission(input)
	if err := domain.ValidateSubmission(input); err != nil {
		return domain.Record{}, err
	}
	if actor == "" {
		return domain.Record{}, fmt.Errorf("actor is required")
	}
	now := s.clock()
	record := input.ToRecord(s.nextID("record"), now)
	if err := s.store.PutRecord(ctx, record); err != nil {
		return domain.Record{}, err
	}
	event := domain.Event{ID: s.nextID("event"), RecordID: record.ID, Kind: "received", ActorID: actor, Detail: "assignment received", CreatedAt: now}
	if err := s.store.AppendEvent(ctx, event); err != nil {
		return domain.Record{}, err
	}
	change, err := domain.NewStatusChange(record, domain.StatusValidated, actor, "input accepted", now)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := s.store.UpdateRecordStatus(ctx, record.ID, change.After); err != nil {
		return domain.Record{}, err
	}
	record.Status = domain.StatusValidated
	change, err = domain.NewStatusChange(record, domain.StatusProcessing, actor, "processing started", now)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := s.store.UpdateRecordStatus(ctx, record.ID, change.After); err != nil {
		return domain.Record{}, err
	}
	record.Status = domain.StatusProcessing
	writer := s.store.NewStatusWriter(true)
	change, err = domain.NewStatusChange(record, domain.StatusImmediate, actor, "submission processed", now)
	if err != nil {
		return domain.Record{}, err
	}
	if err := writer.Queue(change); err != nil {
		return domain.Record{}, err
	}
	closeErr := writer.Close(ctx)
	if closeErr != nil {
		record.Status = domain.StatusDelayed
	} else {
		record.Status = domain.StatusImmediate
	}
	record.Revision++
	record.UpdatedAt = now
	audit := domain.Audit{ID: s.nextID("audit"), RecordID: record.ID, Action: "submit", ActorID: actor, Before: domain.StatusProcessing, After: record.Status, Reason: "submission lifecycle", CreatedAt: now}
	if err := s.store.AppendAudit(ctx, audit); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) Query(ctx context.Context, filter domain.QueryFilter) (domain.Page, error) {
	records, err := s.store.ListRecords(ctx)
	if err != nil {
		return domain.Page{}, err
	}
	list := linkedlist.New[domain.Record]()
	for _, record := range records {
		if domain.MatchRecord(record, filter) && domain.IsVisible(record.Status) {
			list.Append(record)
		}
	}
	values := list.Values()
	page := domain.Page{Limit: filter.Limit, Offset: filter.Offset, Total: len(values)}
	if page.Limit <= 0 {
		page.Limit = 50
	}
	if page.Offset < 0 {
		page.Offset = 0
	}
	if page.Offset >= len(values) {
		page.Items = []domain.Record{}
		return page, nil
	}
	end := page.Offset + page.Limit
	if end > len(values) {
		end = len(values)
	}
	page.Items = append([]domain.Record(nil), values[page.Offset:end]...)
	return page, nil
}
