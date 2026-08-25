package store

import (
	"context"
	"testing"
	"time"

	"coursechain/domain"
)

func TestUserAndEvents(t *testing.T) {
	database, err := Open(t.TempDir() + "/users.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	ctx := context.Background()
	if err := database.PutUser(ctx, domain.User{ID: "u1", Name: "Student", Email: "s@example.com", Role: "student", Active: true, CreatedAt: time.Now().UTC()}); err != nil {
		t.Fatal(err)
	}
	if err := database.AppendEvent(ctx, domain.Event{ID: "e1", RecordID: "r1", Kind: "received", ActorID: "u1", Detail: "ok", CreatedAt: time.Now().UTC()}); err != nil {
		t.Fatal(err)
	}
	if err := database.AppendAudit(ctx, domain.Audit{ID: "a1", RecordID: "r1", Action: "submit", ActorID: "u1", Before: domain.StatusReceived, After: domain.StatusValidated, CreatedAt: time.Now().UTC()}); err != nil {
		t.Fatal(err)
	}
	if count, err := database.CountEvents(ctx, "r1"); err != nil || count != 1 {
		t.Fatalf("event count = %d, %v", count, err)
	}
}
