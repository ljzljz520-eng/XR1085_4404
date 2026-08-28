package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	bolt "go.etcd.io/bbolt"
	"memorialcandle/domain"
)

var (
	BucketMessages   = []byte("memorial_messages")
	BucketAnimations = []byte("candle_animations")
	BucketSnapshots  = []byte("display_snapshots")
	BucketSessions   = []byte("visitor_sessions")
)

var ErrNotFound = errors.New("entity not found")

type Store struct {
	db   *bolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("open bbolt: %w", err)
	}
	s := &Store{db: db, path: path}
	if err := s.createBuckets(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) createBuckets() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range [][]byte{BucketMessages, BucketAnimations, BucketSnapshots, BucketSessions} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}
func (s *Store) ensureOpen() error {
	if s == nil || s.db == nil {
		return errors.New("store is closed")
	}
	return nil
}

func marshal(value any) ([]byte, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("marshal entity: %w", err)
	}
	return data, nil
}
func unmarshal(data []byte, target any) error {
	if len(data) == 0 {
		return ErrNotFound
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("decode entity: %w", err)
	}
	return nil
}

func (s *Store) put(bucket []byte, key string, value any) error {
	if key == "" {
		return errors.New("entity key is required")
	}
	if err := s.ensureOpen(); err != nil {
		return err
	}
	data, err := marshal(value)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucket).Put([]byte(key), data) })
}

func (s *Store) get(bucket []byte, key string, target any) error {
	if key == "" {
		return errors.New("entity key is required")
	}
	if err := s.ensureOpen(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errors.New("bucket missing")
		}
		return unmarshal(b.Get([]byte(key)), target)
	})
}

func (s *Store) delete(bucket []byte, key string) error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bolt.Tx) error { return tx.Bucket(bucket).Delete([]byte(key)) })
}

func Exists(path string) bool { _, err := os.Stat(path); return err == nil }

func (s *Store) PutMemorialMessage(value domain.MemorialMessage) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return s.put(BucketMessages, value.ID, value)
}
func (s *Store) GetMemorialMessage(id string) (domain.MemorialMessage, error) {
	var value domain.MemorialMessage
	err := s.get(BucketMessages, id, &value)
	return value, err
}
func (s *Store) DeleteMemorialMessage(id string) error { return s.delete(BucketMessages, id) }

func (s *Store) PutCandleAnimation(value domain.CandleAnimation) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return s.put(BucketAnimations, value.ID, value)
}
func (s *Store) GetCandleAnimation(id string) (domain.CandleAnimation, error) {
	var value domain.CandleAnimation
	err := s.get(BucketAnimations, id, &value)
	return value, err
}
func (s *Store) DeleteCandleAnimation(id string) error { return s.delete(BucketAnimations, id) }

func (s *Store) PutDisplaySnapshot(value domain.DisplaySnapshot) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return s.put(BucketSnapshots, value.AnimationID, value)
}
func (s *Store) GetDisplaySnapshot(id string) (domain.DisplaySnapshot, error) {
	var value domain.DisplaySnapshot
	err := s.get(BucketSnapshots, id, &value)
	return value, err
}
func (s *Store) DeleteDisplaySnapshot(id string) error { return s.delete(BucketSnapshots, id) }

func (s *Store) PutVisitorSession(value domain.VisitorSession) error {
	if err := value.Validate(); err != nil {
		return err
	}
	return s.put(BucketSessions, value.ID, value)
}
func (s *Store) GetVisitorSession(id string) (domain.VisitorSession, error) {
	var value domain.VisitorSession
	err := s.get(BucketSessions, id, &value)
	return value, err
}
func (s *Store) DeleteVisitorSession(id string) error { return s.delete(BucketSessions, id) }
