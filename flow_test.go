package coursechain

import (
	"context"
	"testing"
	"time"

	"coursechain/domain"
	"coursechain/store"
	"coursechain/workflow"
)

func TestRecordFlow10(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/flow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := workflow.NewService(database)
	service.SetClock(func() time.Time { return time.Unix(1700000000, 0).UTC() })
	record, err := service.Submit(context.Background(), domain.Submission{Course: "课程10", StudentID: "stu-10", Title: "链表作业", Payload: "nodes=4", Tags: []string{"链表", "课程10"}}, "stu-10")
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != domain.StatusImmediate {
		t.Fatalf("submission status = %s, want %s", record.Status, domain.StatusImmediate)
	}
}
