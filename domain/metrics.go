package domain

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CourseSummary struct {
	Course       string         `json:"course"`
	Total        int            `json:"total"`
	Visible      int            `json:"visible"`
	Immediate    int            `json:"immediate"`
	Delayed      int            `json:"delayed"`
	Archived     int            `json:"archived"`
	Rejected     int            `json:"rejected"`
	LatestUpdate time.Time      `json:"latest_update"`
	ByStudent    map[string]int `json:"by_student"`
}

type StatusCount struct {
	Status Status `json:"status"`
	Count  int    `json:"count"`
}

func BuildCourseSummary(records []Record, course string) CourseSummary {
	summary := CourseSummary{Course: course, ByStudent: make(map[string]int)}
	for _, record := range records {
		if course != "" && !strings.EqualFold(record.Course, course) {
			continue
		}
		summary.Total++
		if IsVisible(record.Status) {
			summary.Visible++
		}
		summary.ByStudent[record.StudentID]++
		switch record.Status {
		case StatusImmediate:
			summary.Immediate++
		case StatusDelayed:
			summary.Delayed++
		case StatusArchived:
			summary.Archived++
		case StatusRejected:
			summary.Rejected++
		}
		if record.UpdatedAt.After(summary.LatestUpdate) {
			summary.LatestUpdate = record.UpdatedAt
		}
	}
	return summary
}

func CountStatuses(records []Record) []StatusCount {
	counts := make(map[Status]int)
	for _, record := range records {
		counts[record.Status]++
	}
	result := make([]StatusCount, 0, len(counts))
	for status, count := range counts {
		result = append(result, StatusCount{Status: status, Count: count})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Status < result[j].Status
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func MergeTags(records []Record) []string {
	seen := make(map[string]struct{})
	for _, record := range records {
		for _, tag := range record.Tags {
			key := strings.ToLower(strings.TrimSpace(tag))
			if key != "" {
				seen[key] = struct{}{}
			}
		}
	}
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}

func ValidateRecordSet(records []Record) error {
	seen := make(map[string]struct{}, len(records))
	for _, record := range records {
		if err := ensureRecord(record); err != nil {
			return err
		}
		if _, exists := seen[record.ID]; exists {
			return fmt.Errorf("duplicate record id %s", record.ID)
		}
		seen[record.ID] = struct{}{}
	}
	return nil
}

func ensureRecord(record Record) error {
	if strings.TrimSpace(record.ID) == "" {
		return fmt.Errorf("record id is required")
	}
	if strings.TrimSpace(record.Course) == "" {
		return fmt.Errorf("record %s has no course", record.ID)
	}
	if err := ValidateStatus(record.Status); err != nil {
		return fmt.Errorf("record %s: %w", record.ID, err)
	}
	if record.Revision < 1 {
		return fmt.Errorf("record %s has invalid revision", record.ID)
	}
	return nil
}

func CourseNames(records []Record) []string {
	seen := make(map[string]struct{})
	for _, record := range records {
		if record.Course != "" {
			seen[strings.ToLower(record.Course)] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for course := range seen {
		result = append(result, course)
	}
	sort.Strings(result)
	return result
}

func IsCourseTen(record Record) bool {
	value := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(record.Course), " ", ""))
	return value == "course10" || value == "课程10" || value == "10"
}
