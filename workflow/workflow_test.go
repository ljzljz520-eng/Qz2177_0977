package workflow

import (
	"context"
	"testing"
	"time"

	"coursechain/domain"
	"coursechain/store"
)

func TestWorkflowOne(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/workflow.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := NewService(database)
	service.SetClock(func() time.Time { return time.Unix(100, 0).UTC() })
	if err := service.RegisterUser(context.Background(), domain.User{ID: "teacher", Name: "Teacher", Email: "teacher@example.com", Role: "teacher"}); err != nil {
		t.Fatal(err)
	}
	record, err := service.Submit(context.Background(), domain.Submission{Course: "course10", StudentID: "student", Title: "linked list", Payload: "head->tail"}, "teacher")
	if err != nil {
		t.Fatal(err)
	}
	if record.ID == "" || record.Status != domain.StatusDelayed {
		t.Fatalf("record = %#v", record)
	}
}
