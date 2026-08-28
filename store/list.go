package store

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"memorialcandle/domain"
)

func (s *Store) list(bucket []byte, factory func() any) ([]any, error) {
	if err := s.ensureOpen(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]any, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s missing", bucket)
		}
		return b.ForEach(func(_, data []byte) error {
			item := factory()
			if err := json.Unmarshal(data, item); err != nil {
				return err
			}
			items = append(items, item)
			return nil
		})
	})
	return items, err
}

func (s *Store) ListMessages() ([]domain.MemorialMessage, error) {
	items, err := s.list(BucketMessages, func() any { return &domain.MemorialMessage{} })
	if err != nil {
		return nil, err
	}
	result := make([]domain.MemorialMessage, 0, len(items))
	for _, item := range items {
		result = append(result, *(item.(*domain.MemorialMessage)))
	}
	return result, nil
}

func (s *Store) ListAnimations() ([]domain.CandleAnimation, error) {
	items, err := s.list(BucketAnimations, func() any { return &domain.CandleAnimation{} })
	if err != nil {
		return nil, err
	}
	result := make([]domain.CandleAnimation, 0, len(items))
	for _, item := range items {
		result = append(result, *(item.(*domain.CandleAnimation)))
	}
	return result, nil
}

func (s *Store) ListSnapshots() ([]domain.DisplaySnapshot, error) {
	items, err := s.list(BucketSnapshots, func() any { return &domain.DisplaySnapshot{} })
	if err != nil {
		return nil, err
	}
	result := make([]domain.DisplaySnapshot, 0, len(items))
	for _, item := range items {
		result = append(result, *(item.(*domain.DisplaySnapshot)))
	}
	return result, nil
}

func (s *Store) ListSessions() ([]domain.VisitorSession, error) {
	items, err := s.list(BucketSessions, func() any { return &domain.VisitorSession{} })
	if err != nil {
		return nil, err
	}
	result := make([]domain.VisitorSession, 0, len(items))
	for _, item := range items {
		result = append(result, *(item.(*domain.VisitorSession)))
	}
	return result, nil
}
