package workflow

import (
	"context"
	"fmt"
	"sort"

	"coursechain/domain"
)

type Report struct {
	Course   domain.CourseSummary `json:"course"`
	Statuses []domain.StatusCount `json:"statuses"`
	Tags     []string             `json:"tags"`
	Records  int                  `json:"records"`
}

func (s *Service) Report(ctx context.Context, course string) (Report, error) {
	records, err := s.store.ListRecords(ctx)
	if err != nil {
		return Report{}, err
	}
	selected := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if course == "" || record.Course == course {
			selected = append(selected, record)
		}
	}
	return Report{Course: domain.BuildCourseSummary(selected, course), Statuses: domain.CountStatuses(selected), Tags: domain.MergeTags(selected), Records: len(selected)}, nil
}

func (s *Service) BulkRegister(ctx context.Context, users []domain.User) (int, error) {
	registered := 0
	for _, user := range users {
		if err := s.RegisterUser(ctx, user); err != nil {
			return registered, err
		}
		registered++
	}
	return registered, nil
}

func (s *Service) ProcessQueue(ctx context.Context, ids []string, actor string) ([]domain.Record, error) {
	result := make([]domain.Record, 0, len(ids))
	for _, id := range ids {
		record, err := s.Reconcile(ctx, id, actor)
		if err != nil {
			return result, fmt.Errorf("process %s: %w", id, err)
		}
		result = append(result, record)
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].UpdatedAt.Before(result[j].UpdatedAt) })
	return result, nil
}

func (s *Service) CourseReady(ctx context.Context, course string) (bool, error) {
	records, err := s.store.ListRecords(ctx)
	if err != nil {
		return false, err
	}
	seen := 0
	for _, record := range records {
		if record.Course != course {
			continue
		}
		seen++
		if record.Status != domain.StatusImmediate && record.Status != domain.StatusArchived {
			return false, nil
		}
	}
	return seen > 0, nil
}

func (s *Service) Timeline(ctx context.Context, id string) ([]string, error) {
	tracking, err := s.Track(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(tracking.Events)+len(tracking.Audits))
	for _, event := range tracking.Events {
		result = append(result, event.CreatedAt.Format("2006-01-02T15:04:05Z07:00")+" "+event.Kind)
	}
	for _, audit := range tracking.Audits {
		result = append(result, audit.CreatedAt.Format("2006-01-02T15:04:05Z07:00")+" "+audit.Action)
	}
	sort.Strings(result)
	return result, nil
}
