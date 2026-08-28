package domain

import (
	"fmt"
	"time"
)

type StatusChange struct {
	RecordID string
	Before   Status
	After    Status
	ActorID  string
	Reason   string
	At       time.Time
}

func NewStatusChange(record Record, next Status, actor string, reason string, now time.Time) (StatusChange, error) {
	if err := ValidateTransition(record.Status, next); err != nil {
		return StatusChange{}, err
	}
	if actor == "" {
		return StatusChange{}, fmt.Errorf("actor is required")
	}
	return StatusChange{RecordID: record.ID, Before: record.Status, After: next, ActorID: actor, Reason: reason, At: now}, nil
}

func ApplyStatus(record *Record, change StatusChange) error {
	if record == nil {
		return fmt.Errorf("record is nil")
	}
	if record.ID != change.RecordID || record.Status != change.Before {
		return fmt.Errorf("stale status change")
	}
	record.Status = change.After
	record.Revision++
	record.UpdatedAt = change.At
	return nil
}

func StatusLabel(status Status) string {
	labels := map[Status]string{StatusReceived: "received", StatusValidated: "validated", StatusProcessing: "processing", StatusImmediate: "immediate", StatusDelayed: "delayed", StatusArchived: "archived", StatusRejected: "rejected"}
	if label, ok := labels[status]; ok {
		return label
	}
	return "unknown"
}

func IsVisible(status Status) bool {
	return status != StatusRejected
}

func StatusRank(status Status) int {
	switch status {
	case StatusImmediate:
		return 5
	case StatusValidated:
		return 4
	case StatusProcessing:
		return 3
	case StatusDelayed:
		return 2
	case StatusReceived:
		return 1
	default:
		return 0
	}
}
