package store

import (
	"context"
	"fmt"

	"coursechain/domain"
)

type StatusWriter struct {
	store        *Store
	pending      []domain.StatusChange
	closed       bool
	flushOnClose bool
}

func (s *Store) NewStatusWriter(flushOnClose bool) *StatusWriter {
	return &StatusWriter{store: s, pending: make([]domain.StatusChange, 0), flushOnClose: flushOnClose}
}

func (w *StatusWriter) Queue(change domain.StatusChange) error {
	if w == nil || w.closed {
		return fmt.Errorf("status writer is closed")
	}
	if change.RecordID == "" {
		return fmt.Errorf("record id is required")
	}
	w.pending = append(w.pending, change)
	return nil
}

func (w *StatusWriter) Flush(ctx context.Context) error {
	if w == nil {
		return fmt.Errorf("status writer is nil")
	}
	if w.closed {
		return fmt.Errorf("status writer is closed")
	}
	for _, change := range w.pending {
		if _, err := w.store.UpdateRecordStatus(ctx, change.RecordID, change.After); err != nil {
			return err
		}
	}
	w.pending = w.pending[:0]
	return nil
}

func (w *StatusWriter) Pending() int {
	if w == nil {
		return 0
	}
	return len(w.pending)
}

func (w *StatusWriter) Close(ctx context.Context) error {
	if w == nil || w.closed {
		return nil
	}
	w.closed = true
	if w.flushOnClose {
		return w.Flush(ctx)
	}
	w.pending = nil
	return nil
}
