package store

import (
	"errors"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) Healthy() error {
	if err := s.ensureOpen(); err != nil {
		return err
	}
	return s.db.View(func(tx *bolt.Tx) error {
		for _, bucket := range [][]byte{BucketMessages, BucketAnimations, BucketSnapshots, BucketSessions} {
			if tx.Bucket(bucket) == nil {
				return errors.New("required bucket missing")
			}
		}
		return nil
	})
}
