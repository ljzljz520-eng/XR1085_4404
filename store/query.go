package store

import (
	"memorialcandle/domain"
	"sort"
	"strings"
)

type MessageQuery struct {
	Term       string
	ActiveOnly bool
	Limit      int
}

func (s *Store) QueryMessages(query MessageQuery) ([]domain.MemorialMessage, error) {
	messages, err := s.ListMessages()
	if err != nil {
		return nil, err
	}
	term := strings.ToLower(strings.TrimSpace(query.Term))
	result := make([]domain.MemorialMessage, 0)
	for _, message := range messages {
		if query.ActiveOnly && !message.Active {
			continue
		}
		if term != "" && !strings.Contains(strings.ToLower(message.Text), term) && !strings.Contains(strings.ToLower(message.Author), term) {
			continue
		}
		result = append(result, message)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	if query.Limit > 0 && len(result) > query.Limit {
		result = result[:query.Limit]
	}
	return result, nil
}

func (s *Store) CountAll() (int, int, int, int, error) {
	messages, err := s.ListMessages()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	animations, err := s.ListAnimations()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	snapshots, err := s.ListSnapshots()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	sessions, err := s.ListSessions()
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return len(messages), len(animations), len(snapshots), len(sessions), nil
}

func (s *Store) MessageIDs() ([]string, error) {
	messages, err := s.ListMessages()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(messages))
	for _, m := range messages {
		ids = append(ids, m.ID)
	}
	sort.Strings(ids)
	return ids, nil
}
