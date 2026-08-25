package store

import (
	"context"
	"fmt"
	"strings"

	"coursechain/domain"
)

type ReconcileItem struct {
	RecordID string        `json:"record_id"`
	From     domain.Status `json:"from"`
	To       domain.Status `json:"to"`
	Reason   string        `json:"reason"`
}

func (s *Store) ReconcileStatuses(ctx context.Context, actor string) ([]ReconcileItem, error) {
	if strings.TrimSpace(actor) == "" {
		return nil, fmt.Errorf("actor is required")
	}
	records, err := s.ListRecords(ctx)
	if err != nil {
		return nil, err
	}
	changes := make([]ReconcileItem, 0)
	for _, record := range records {
		if record.Status != domain.StatusProcessing {
			continue
		}
		if _, updateErr := s.UpdateRecordStatus(ctx, record.ID, domain.StatusImmediate); updateErr != nil {
			changes = append(changes, ReconcileItem{RecordID: record.ID, From: record.Status, To: domain.StatusDelayed, Reason: updateErr.Error()})
			continue
		}
		changes = append(changes, ReconcileItem{RecordID: record.ID, From: record.Status, To: domain.StatusImmediate, Reason: "reconciled by " + actor})
	}
	return changes, nil
}

func (s *Store) CountByCourse(ctx context.Context) (map[string]int, error) {
	records, err := s.ListRecords(ctx)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, record := range records {
		course := strings.ToLower(strings.TrimSpace(record.Course))
		if course != "" {
			counts[course]++
		}
	}
	return counts, nil
}

func (s *Store) RemoveUserData(ctx context.Context, userID string) (int, error) {
	if strings.TrimSpace(userID) == "" {
		return 0, fmt.Errorf("user id is required")
	}
	records, err := s.ListRecords(ctx)
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, record := range records {
		if record.StudentID == userID {
			if deleteErr := s.DeleteRecord(ctx, record.ID); deleteErr == nil {
				removed++
			}
		}
	}
	if err := s.SetUserActive(ctx, userID, false); err != nil && err != ErrNotFound {
		return removed, err
	}
	return removed, nil
}
