package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"coursechain/domain"
)

type Backup struct {
	Snapshot Snapshot `json:"snapshot"`
	Version  int      `json:"version"`
}

func (s *Store) ExportBackup(ctx context.Context, path string) error {
	if filepath.Clean(path) == "." || path == "" {
		return fmt.Errorf("backup path is required")
	}
	snapshot, err := s.Snapshot(ctx)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(Backup{Snapshot: snapshot, Version: 1}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(filepath.Clean(path)), 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Clean(path), data, 0o600)
}

func DecodeBackup(data []byte) (Backup, error) {
	var backup Backup
	if err := json.Unmarshal(data, &backup); err != nil {
		return Backup{}, err
	}
	if backup.Version != 1 {
		return Backup{}, fmt.Errorf("unsupported backup version %d", backup.Version)
	}
	if err := backup.Snapshot.Validate(); err != nil {
		return Backup{}, err
	}
	return backup, nil
}

func (s *Store) RestoreRecords(ctx context.Context, records []domain.Record) error {
	if err := domain.ValidateRecordSet(records); err != nil {
		return err
	}
	return s.PutRecords(ctx, records)
}
