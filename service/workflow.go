package service

import (
	"errors"
	"memorialcandle/animation"
	"memorialcandle/domain"
	"memorialcandle/render"
	"memorialcandle/store"
)

func storeAnimation(s *Service, bundle store.Bundle) error {
	if bundle.Animation.Stage == domain.StageQuiet && !animation.IsIdle(bundle.Animation) {
		return errors.New("quiet animation is not idle")
	}
	bundle.Snapshot = render.Render(bundle.Message, bundle.Animation)
	return s.store.SaveBundle(bundle)
}

func (s *Service) EnsureQuiet(messageID, sessionID string) (store.Bundle, error) {
	bundle, err := s.store.LoadBundle(messageID, sessionID)
	if err != nil {
		return store.Bundle{}, err
	}
	if bundle.Animation.Layers != 1 {
		return store.Bundle{}, errors.New("multiple layers require manual review")
	}
	if !animation.IsIdle(bundle.Animation) {
		if _, err := s.CompleteAnimation(messageID, sessionID); err != nil {
			return store.Bundle{}, err
		}
		bundle, err = s.store.LoadBundle(messageID, sessionID)
	}
	return bundle, err
}

func (s *Service) Replay(messageID, sessionID string, clicks int) (store.Bundle, error) {
	if clicks < 1 || clicks > domain.MaxClicks {
		return store.Bundle{}, errors.New("clicks out of range")
	}
	var bundle store.Bundle
	var err error
	for i := 0; i < clicks; i++ {
		bundle, err = s.ClickCandle(messageID, sessionID)
		if err != nil {
			return store.Bundle{}, err
		}
	}
	return bundle, nil
}

func (s *Service) Close() error {
	if s == nil || s.store == nil {
		return nil
	}
	return s.store.Close()
}
