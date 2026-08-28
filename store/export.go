package store

import (
	"encoding/json"
	"fmt"
	"memorialcandle/domain"
)

type ExportData struct {
	Messages   []domain.MemorialMessage `json:"messages"`
	Animations []domain.CandleAnimation `json:"animations"`
	Snapshots  []domain.DisplaySnapshot `json:"snapshots"`
	Sessions   []domain.VisitorSession  `json:"sessions"`
}

func (s *Store) Export() (ExportData, error) {
	messages, err := s.ListMessages()
	if err != nil {
		return ExportData{}, err
	}
	animations, err := s.ListAnimations()
	if err != nil {
		return ExportData{}, err
	}
	snapshots, err := s.ListSnapshots()
	if err != nil {
		return ExportData{}, err
	}
	sessions, err := s.ListSessions()
	if err != nil {
		return ExportData{}, err
	}
	return ExportData{Messages: messages, Animations: animations, Snapshots: snapshots, Sessions: sessions}, nil
}

func (data ExportData) Validate() error {
	if err := validateExportMessages(data.Messages); err != nil {
		return err
	}
	for _, a := range data.Animations {
		if err := a.Validate(); err != nil {
			return err
		}
	}
	for _, d := range data.Snapshots {
		if err := d.Validate(); err != nil {
			return err
		}
	}
	for _, s := range data.Sessions {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func validateExportMessages(messages []domain.MemorialMessage) error {
	seen := map[string]bool{}
	for _, m := range messages {
		if err := m.Validate(); err != nil {
			return err
		}
		if seen[m.ID] {
			return fmt.Errorf("duplicate export message %s", m.ID)
		}
		seen[m.ID] = true
	}
	return nil
}

func (data ExportData) JSON() ([]byte, error) {
	if err := data.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(data, "", "  ")
}

func ImportJSON(raw []byte) (ExportData, error) {
	var data ExportData
	if err := json.Unmarshal(raw, &data); err != nil {
		return ExportData{}, err
	}
	if err := data.Validate(); err != nil {
		return ExportData{}, err
	}
	return data, nil
}

func (s *Store) Import(data ExportData) error {
	if err := data.Validate(); err != nil {
		return err
	}
	for _, m := range data.Messages {
		if err := s.PutMemorialMessage(m); err != nil {
			return err
		}
	}
	for _, a := range data.Animations {
		if err := s.PutCandleAnimation(a); err != nil {
			return err
		}
	}
	for _, d := range data.Snapshots {
		if err := s.PutDisplaySnapshot(d); err != nil {
			return err
		}
	}
	for _, v := range data.Sessions {
		if err := s.PutVisitorSession(v); err != nil {
			return err
		}
	}
	return nil
}
