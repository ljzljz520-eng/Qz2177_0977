package api

import (
	"net/http/httptest"
	"strings"
	"testing"

	"coursechain/store"
	"coursechain/workflow"
)

func TestWorkflowThree(t *testing.T) {
	database, err := store.Open(t.TempDir() + "/api.db")
	if err != nil {
		t.Fatal(err)
	}
	defer database.Close()
	request := httptest.NewRequest("GET", "/healthz", nil)
	response := httptest.NewRecorder()
	NewRouter(workflow.NewService(database)).ServeHTTP(response, request)
	if response.Code != 200 || !strings.Contains(response.Body.String(), "course10") {
		t.Fatalf("health response: %d %s", response.Code, response.Body.String())
	}
}
