package query

import (
	"context"
	"sort"
	"strings"

	"coursechain/domain"
	"coursechain/linkedlist"
	"coursechain/store"
)

type Engine struct {
	store *store.Store
}

func New(s *store.Store) *Engine {
	return &Engine{store: s}
}

func (e *Engine) Search(ctx context.Context, filter domain.QueryFilter) (domain.Page, error) {
	records, err := e.store.ListRecords(ctx)
	if err != nil {
		return domain.Page{}, err
	}
	list := linkedlist.New[domain.Record]()
	for _, record := range records {
		if domain.MatchRecord(record, filter) {
			list.Append(record)
		}
	}
	if filter.Search != "" {
		needle := strings.ToLower(filter.Search)
		list = linkedlist.Filter(list, func(record domain.Record) bool {
			return strings.Contains(strings.ToLower(record.Title), needle) || strings.Contains(strings.ToLower(record.Payload), needle)
		})
	}
	items := list.Values()
	sort.SliceStable(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
	return paginate(items, filter), nil
}

func paginate(items []domain.Record, filter domain.QueryFilter) domain.Page {
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	page := domain.Page{Total: len(items), Limit: limit, Offset: offset, Items: []domain.Record{}}
	if offset >= len(items) {
		return page
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	page.Items = append(page.Items, items[offset:end]...)
	return page
}

func GroupByStatus(records []domain.Record) map[domain.Status]int {
	grouped := make(map[domain.Status]int)
	for _, record := range records {
		grouped[record.Status]++
	}
	return grouped
}

func Summarize(records []domain.Record) map[string]any {
	counts := GroupByStatus(records)
	return map[string]any{"total": len(records), "statuses": counts}
}
