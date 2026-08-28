package query

import (
	"fmt"
	"strings"

	"coursechain/domain"
)

type Plan struct {
	Description string
	Filter      domain.QueryFilter
}

func ParsePlan(course, status, search string, limit, offset int) (Plan, error) {
	course = strings.TrimSpace(course)
	status = strings.TrimSpace(status)
	search = strings.TrimSpace(search)
	filter := domain.QueryFilter{Course: course, Search: search, Limit: limit, Offset: offset}
	if status != "" {
		filter.Status = domain.Status(status)
		if err := domain.ValidateStatus(filter.Status); err != nil {
			return Plan{}, err
		}
	}
	if limit < 0 || offset < 0 {
		return Plan{}, fmt.Errorf("limit and offset cannot be negative")
	}
	parts := []string{"course=" + course}
	if status != "" {
		parts = append(parts, "status="+status)
	}
	if search != "" {
		parts = append(parts, "search="+search)
	}
	return Plan{Description: strings.Join(parts, " "), Filter: filter}, nil
}

func SortByRank(records []domain.Record) []domain.Record {
	copyItems := append([]domain.Record(nil), records...)
	for index := 1; index < len(copyItems); index++ {
		current := copyItems[index]
		position := index - 1
		for position >= 0 && domain.StatusRank(copyItems[position].Status) < domain.StatusRank(current.Status) {
			copyItems[position+1] = copyItems[position]
			position--
		}
		copyItems[position+1] = current
	}
	return copyItems
}
