package domain

import (
	"time"
)

type Status string

const (
	StatusReceived   Status = "received"
	StatusValidated  Status = "validated"
	StatusProcessing Status = "processing"
	StatusImmediate  Status = "immediate"
	StatusDelayed    Status = "delayed"
	StatusArchived   Status = "archived"
	StatusRejected   Status = "rejected"
)

type Record struct {
	ID          string    `json:"id"`
	Course      string    `json:"course"`
	StudentID   string    `json:"student_id"`
	Title       string    `json:"title"`
	Payload     string    `json:"payload"`
	Status      Status    `json:"status"`
	Revision    int       `json:"revision"`
	SubmittedAt time.Time `json:"submitted_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

type Event struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Kind      string    `json:"kind"`
	ActorID   string    `json:"actor_id"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

type Audit struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	Action    string    `json:"action"`
	ActorID   string    `json:"actor_id"`
	Before    Status    `json:"before"`
	After     Status    `json:"after"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

type Submission struct {
	Course    string
	StudentID string
	Title     string
	Payload   string
	Tags      []string
}

type QueryFilter struct {
	Course    string
	StudentID string
	Status    Status
	Tag       string
	Search    string
	Limit     int
	Offset    int
}

type Page struct {
	Items  []Record `json:"items"`
	Total  int      `json:"total"`
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
}

type Review struct {
	RecordID string `json:"record_id"`
	Reviewer string `json:"reviewer"`
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

type Tracking struct {
	RecordID string  `json:"record_id"`
	Events   []Event `json:"events"`
	Audits   []Audit `json:"audits"`
}

func (r Record) IsTerminal() bool {
	return r.Status == StatusArchived || r.Status == StatusRejected
}

func (r Record) Clone() Record {
	copyRecord := r
	copyRecord.Tags = append([]string(nil), r.Tags...)
	return copyRecord
}

func (s Submission) ToRecord(id string, now time.Time) Record {
	return Record{ID: id, Course: s.Course, StudentID: s.StudentID, Title: s.Title, Payload: s.Payload, Status: StatusReceived, Revision: 1, SubmittedAt: now, UpdatedAt: now, Tags: append([]string(nil), s.Tags...)}
}
