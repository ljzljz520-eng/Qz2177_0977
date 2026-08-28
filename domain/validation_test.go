package domain

import "testing"

func TestValidationAndTransitions(t *testing.T) {
	input := NormalizeSubmission(Submission{Course: " Course10 ", StudentID: " s1 ", Title: " title ", Payload: " body ", Tags: []string{"Go", "go", ""}})
	if err := ValidateSubmission(input); err != nil {
		t.Fatal(err)
	}
	if input.Course != "course10" || len(input.Tags) != 1 {
		t.Fatalf("normalized input: %#v", input)
	}
	if !CanTransition(StatusProcessing, StatusImmediate) {
		t.Fatal("processing should transition to immediate")
	}
	if CanTransition(StatusArchived, StatusImmediate) {
		t.Fatal("archived should not transition to immediate")
	}
}
