package query

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"

	"coursechain/domain"
)

func CSVHeader() []string {
	return []string{"id", "course", "student_id", "title", "status", "revision", "submitted_at", "updated_at"}
}

func CSVRow(record domain.Record) []string {
	return []string{record.ID, record.Course, record.StudentID, record.Title, string(record.Status), strconv.Itoa(record.Revision), record.SubmittedAt.UTC().Format("2006-01-02T15:04:05Z07:00"), record.UpdatedAt.UTC().Format("2006-01-02T15:04:05Z07:00")}
}

func WriteCSV(writer io.Writer, records []domain.Record) error {
	if writer == nil {
		return fmt.Errorf("writer is required")
	}
	output := csv.NewWriter(writer)
	if err := output.Write(CSVHeader()); err != nil {
		return err
	}
	for _, record := range records {
		if err := output.Write(CSVRow(record)); err != nil {
			return err
		}
	}
	output.Flush()
	return output.Error()
}

func ParseTagList(raw string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)
	for _, value := range strings.Split(raw, ",") {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func FilterByCourse(records []domain.Record, course string) []domain.Record {
	result := make([]domain.Record, 0, len(records))
	for _, record := range records {
		if course == "" || strings.EqualFold(record.Course, course) {
			result = append(result, record)
		}
	}
	return result
}

func Facets(records []domain.Record) map[string][]string {
	result := map[string][]string{"courses": {}, "students": {}, "tags": {}}
	courses := make(map[string]struct{})
	students := make(map[string]struct{})
	tags := make(map[string]struct{})
	for _, record := range records {
		courses[record.Course] = struct{}{}
		students[record.StudentID] = struct{}{}
		for _, tag := range record.Tags {
			tags[tag] = struct{}{}
		}
	}
	for value := range courses {
		result["courses"] = append(result["courses"], value)
	}
	for value := range students {
		result["students"] = append(result["students"], value)
	}
	for value := range tags {
		result["tags"] = append(result["tags"], value)
	}
	return result
}
