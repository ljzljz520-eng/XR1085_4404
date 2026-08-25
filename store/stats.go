package store

import (
	"errors"
	"memorialcandle/domain"
)

type Stats struct {
	Messages        int
	ActiveMessages  int
	Animations      int
	QuietAnimations int
	Snapshots       int
	QuietSnapshots  int
	Sessions        int
	Interactions    int
}

func (s *Store) Stats() (Stats, error) {
	messages, err := s.ListMessages()
	if err != nil {
		return Stats{}, err
	}
	animations, err := s.ListAnimations()
	if err != nil {
		return Stats{}, err
	}
	snapshots, err := s.ListSnapshots()
	if err != nil {
		return Stats{}, err
	}
	sessions, err := s.ListSessions()
	if err != nil {
		return Stats{}, err
	}
	stats := Stats{Messages: len(messages), Animations: len(animations), Snapshots: len(snapshots), Sessions: len(sessions)}
	for _, message := range messages {
		if message.Active {
			stats.ActiveMessages++
		}
	}
	for _, animation := range animations {
		if animation.Quiet {
			stats.QuietAnimations++
		}
	}
	for _, snapshot := range snapshots {
		if snapshot.Quiet {
			stats.QuietSnapshots++
		}
	}
	for _, session := range sessions {
		stats.Interactions += session.InteractionCount
	}
	return stats, nil
}

func (stats Stats) Valid() bool {
	return stats.Messages >= stats.ActiveMessages && stats.Animations >= stats.QuietAnimations && stats.Snapshots >= stats.QuietSnapshots && stats.Sessions >= 0 && stats.Interactions >= 0
}

func (s *Store) ValidateAll() error {
	data, err := s.Export()
	if err != nil {
		return err
	}
	if err := data.Validate(); err != nil {
		return err
	}
	for _, animation := range data.Animations {
		if animation.Stage == domain.StageQuiet && animation.Layers != 1 {
			return errors.New("stored quiet animation has multiple layers")
		}
	}
	return nil
}

func (s *Store) UpdateAnimation(id string, update func(*domain.CandleAnimation) error) error {
	animation, err := s.GetCandleAnimation(id)
	if err != nil {
		return err
	}
	if err := update(&animation); err != nil {
		return err
	}
	return s.PutCandleAnimation(animation)
}

func (s *Store) UpdateSession(id string, update func(*domain.VisitorSession) error) error {
	session, err := s.GetVisitorSession(id)
	if err != nil {
		return err
	}
	if err := update(&session); err != nil {
		return err
	}
	return s.PutVisitorSession(session)
}

func (s *Store) ResetSnapshot(animationID string) error {
	snapshot, err := s.GetDisplaySnapshot(animationID)
	if err != nil {
		return err
	}
	snapshot.Stage = domain.StageQuiet
	snapshot.ParticleCount = 1
	snapshot.Quiet = true
	return s.PutDisplaySnapshot(snapshot)
}
