package workflow

import (
	"context"
	"testing"

	"coursechain/domain"
	"coursechain/store"
)

func TestWorkflowTwo(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/registration.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	service := NewService(database)
	if err := service.RegisterUser(context.Background(), domain.User{ID: "student", Name: "Student", Email: "student@example.com"}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.EnsureActor(context.Background(), "student"); err != nil {
		t.Fatal(err)
	}
	if err := service.DeactivateUser(context.Background(), "student"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.EnsureActor(context.Background(), "student"); err == nil {
		t.Fatal("inactive actor was accepted")
	}
}
