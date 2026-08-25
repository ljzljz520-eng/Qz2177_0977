package store

import (
	"context"
	"testing"
	"time"

	"coursechain/domain"
)

func TestRecordCRUD(t *testing.T) {
	database, err := Open(t.TempDir() + "/records.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	record := domain.Record{ID: "record-1", Course: "course10", StudentID: "student", Title: "title", Payload: "body", Status: domain.StatusReceived, Revision: 1, SubmittedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	if err := database.PutRecord(context.Background(), record); err != nil {
		t.Fatal(err)
	}
	got, err := database.GetRecord(context.Background(), record.ID)
	if err != nil || got.Title != record.Title {
		t.Fatalf("get = %#v, %v", got, err)
	}
	updated, err := database.UpdateRecordStatus(context.Background(), record.ID, domain.StatusValidated)
	if err != nil || updated.Status != domain.StatusValidated {
		t.Fatalf("update = %#v, %v", updated, err)
	}
	if err := database.DeleteRecord(context.Background(), record.ID); err != nil {
		t.Fatal(err)
	}
}
