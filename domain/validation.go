package domain

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingCourse  = errors.New("course is required")
	ErrMissingStudent = errors.New("student is required")
	ErrMissingTitle   = errors.New("title is required")
	ErrMissingPayload = errors.New("payload is required")
	ErrInvalidStatus  = errors.New("invalid status")
)

func ValidateSubmission(s Submission) error {
	if strings.TrimSpace(s.Course) == "" {
		return ErrMissingCourse
	}
	if strings.TrimSpace(s.StudentID) == "" {
		return ErrMissingStudent
	}
	if strings.TrimSpace(s.Title) == "" {
		return ErrMissingTitle
	}
	if strings.TrimSpace(s.Payload) == "" {
		return ErrMissingPayload
	}
	if len(s.Payload) > 20000 {
		return fmt.Errorf("payload exceeds 20000 bytes")
	}
	return nil
}

func NormalizeSubmission(s Submission) Submission {
	s.Course = strings.ToLower(strings.TrimSpace(s.Course))
	s.StudentID = strings.TrimSpace(s.StudentID)
	s.Title = strings.TrimSpace(s.Title)
	s.Payload = strings.TrimSpace(s.Payload)
	clean := make([]string, 0, len(s.Tags))
	seen := make(map[string]struct{})
	for _, tag := range s.Tags {
		tag = strings.ToLower(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		clean = append(clean, tag)
	}
	s.Tags = clean
	return s
}

func ValidateStatus(status Status) error {
	switch status {
	case StatusReceived, StatusValidated, StatusProcessing, StatusImmediate, StatusDelayed, StatusArchived, StatusRejected:
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrInvalidStatus, status)
	}
}

func CanTransition(from Status, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusReceived:
		return to == StatusValidated || to == StatusRejected
	case StatusValidated:
		return to == StatusProcessing || to == StatusRejected
	case StatusProcessing:
		return to == StatusImmediate || to == StatusDelayed || to == StatusRejected
	case StatusImmediate, StatusDelayed:
		return to == StatusArchived || to == StatusRejected
	default:
		return false
	}
}

func ValidateTransition(from Status, to Status) error {
	if err := ValidateStatus(from); err != nil {
		return err
	}
	if err := ValidateStatus(to); err != nil {
		return err
	}
	if !CanTransition(from, to) {
		return fmt.Errorf("transition %s -> %s is not allowed", from, to)
	}
	return nil
}

func MatchRecord(r Record, filter QueryFilter) bool {
	if filter.Course != "" && !strings.EqualFold(r.Course, filter.Course) {
		return false
	}
	if filter.StudentID != "" && r.StudentID != filter.StudentID {
		return false
	}
	if filter.Status != "" && r.Status != filter.Status {
		return false
	}
	if filter.Tag != "" {
		found := false
		for _, tag := range r.Tags {
			if strings.EqualFold(tag, filter.Tag) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if filter.Search != "" {
		needle := strings.ToLower(filter.Search)
		if !strings.Contains(strings.ToLower(r.Title), needle) && !strings.Contains(strings.ToLower(r.Payload), needle) {
			return false
		}
	}
	return true
}
