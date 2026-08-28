package store

import (
	"errors"
	bolt "go.etcd.io/bbolt"
	"memorialcandle/domain"
)

type Bundle struct {
	Message   domain.MemorialMessage
	Animation domain.CandleAnimation
	Snapshot  domain.DisplaySnapshot
	Session   domain.VisitorSession
}

func (s *Store) SaveBundle(bundle Bundle) error {
	if err := bundle.Message.Validate(); err != nil {
		return err
	}
	if err := bundle.Animation.Validate(); err != nil {
		return err
	}
	if err := bundle.Snapshot.Validate(); err != nil {
		return err
	}
	if err := bundle.Session.Validate(); err != nil {
		return err
	}
	if bundle.Animation.MessageID != bundle.Message.ID || bundle.Snapshot.AnimationID != bundle.Animation.ID {
		return errors.New("bundle links are invalid")
	}
	if err := s.ensureOpen(); err != nil {
		return err
	}
	message, err := marshal(bundle.Message)
	if err != nil {
		return err
	}
	animation, err := marshal(bundle.Animation)
	if err != nil {
		return err
	}
	snapshot, err := marshal(bundle.Snapshot)
	if err != nil {
		return err
	}
	session, err := marshal(bundle.Session)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := tx.Bucket(BucketMessages).Put([]byte(bundle.Message.ID), message); err != nil {
			return err
		}
		if err := tx.Bucket(BucketAnimations).Put([]byte(bundle.Animation.ID), animation); err != nil {
			return err
		}
		if err := tx.Bucket(BucketSnapshots).Put([]byte(bundle.Snapshot.AnimationID), snapshot); err != nil {
			return err
		}
		return tx.Bucket(BucketSessions).Put([]byte(bundle.Session.ID), session)
	})
}

func (s *Store) LoadBundle(messageID, sessionID string) (Bundle, error) {
	message, err := s.GetMemorialMessage(messageID)
	if err != nil {
		return Bundle{}, err
	}
	animation, err := s.GetCandleAnimation(domain.AnimationID(message.ID))
	if err != nil {
		return Bundle{}, err
	}
	snapshot, err := s.GetDisplaySnapshot(animation.ID)
	if err != nil {
		return Bundle{}, err
	}
	session, err := s.GetVisitorSession(sessionID)
	if err != nil {
		return Bundle{}, err
	}
	return Bundle{Message: message, Animation: animation, Snapshot: snapshot, Session: session}, nil
}
