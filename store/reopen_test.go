package store

import (
	"context"
	"testing"
	"time"

	"coursechain/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/records.db"
	database, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.Record{ID: "reopen-record", Course: "course10", StudentID: "student", Title: "reopen", Payload: "persisted", Status: domain.StatusImmediate, Revision: 1, SubmittedAt: time.Unix(10, 0), UpdatedAt: time.Unix(10, 0)}
	if err := database.PutRecord(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	if err := database.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	got, err := reopened.GetRecord(context.Background(), record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Payload != record.Payload || got.Status != record.Status {
		t.Fatalf("reopened record = %#v", got)
	}
}
