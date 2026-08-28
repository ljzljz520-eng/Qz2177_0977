package query

import (
	"testing"
	"time"

	"coursechain/domain"
)

func TestPlanAndRank(t *testing.T) {
	plan, err := ParsePlan("course10", "immediate", "graphs", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Filter.Status != domain.StatusImmediate {
		t.Fatalf("status = %s", plan.Filter.Status)
	}
	items := []domain.Record{{ID: "a", Status: domain.StatusDelayed, UpdatedAt: time.Unix(1, 0)}, {ID: "b", Status: domain.StatusImmediate, UpdatedAt: time.Unix(2, 0)}}
	ranked := SortByRank(items)
	if ranked[0].ID != "b" {
		t.Fatalf("ranked = %#v", ranked)
	}
}
