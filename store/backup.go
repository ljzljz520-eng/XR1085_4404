package store

import (
	"encoding/json"
	"fmt"
)

func (s *Store) ExportJSON() ([]byte, error) {
	data, err := s.Export()
	if err != nil {
		return nil, err
	}
	return data.JSON()
}
func (s *Store) ImportJSON(raw []byte) error {
	data, err := ImportJSON(raw)
	if err != nil {
		return err
	}
	return s.Import(data)
}
func (s *Store) BackupSummary() (string, error) {
	stats, err := s.Stats()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("messages=%d animations=%d snapshots=%d sessions=%d interactions=%d", stats.Messages, stats.Animations, stats.Snapshots, stats.Sessions, stats.Interactions), nil
}
func (s *Store) EncodeStats() ([]byte, error) {
	stats, err := s.Stats()
	if err != nil {
		return nil, err
	}
	return json.Marshal(stats)
}
