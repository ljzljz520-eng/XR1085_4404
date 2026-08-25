package service

import (
	"errors"
	"fmt"
	"memorialcandle/animation"
	"memorialcandle/domain"
	"memorialcandle/render"
	"memorialcandle/store"
	"memorialcandle/validate"
)

type Service struct {
	store    *store.Store
	sequence int
}

func New(st *store.Store) (*Service, error) {
	if st == nil {
		return nil, errors.New("store is required")
	}
	if err := st.Healthy(); err != nil {
		return nil, err
	}
	return &Service{store: st, sequence: 1}, nil
}
func (s *Service) Store() *store.Store      { return s.store }
func (s *Service) SetSequence(sequence int) { s.sequence = domain.NormalizeSequence(sequence) }

func (s *Service) RecordMemorial(visitor, author, text, startedAt string) (store.Bundle, error) {
	if err := validate.Visitor(visitor); err != nil {
		return store.Bundle{}, err
	}
	if err := validate.Message(author, text); err != nil {
		return store.Bundle{}, err
	}
	sequence := domain.NormalizeSequence(s.sequence)
	message := domain.NewMemorialMessage(domain.MessageID(sequence), author, text, startedAt)
	session := domain.NewVisitorSession(domain.SessionID(sequence), visitor, startedAt)
	animationState := domain.NewCandleAnimation(domain.AnimationID(message.ID), message.ID)
	snapshot := render.Render(message, animationState)
	bundle := store.Bundle{Message: message, Animation: animationState, Snapshot: snapshot, Session: session}
	if err := store.ValidateBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	if err := s.store.SaveBundle(bundle); err != nil {
		return store.Bundle{}, fmt.Errorf("record memorial: %w", err)
	}
	s.sequence = sequence + 1
	return bundle, nil
}

func (s *Service) QueryDisplay(messageID, sessionID string) (store.Bundle, error) {
	bundle, err := s.store.LoadBundle(messageID, sessionID)
	if err != nil {
		return store.Bundle{}, fmt.Errorf("query display: %w", err)
	}
	if err := store.ValidateBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	if bundle.Snapshot != render.Render(bundle.Message, bundle.Animation) {
		return store.Bundle{}, errors.New("stored snapshot differs from render")
	}
	return bundle, nil
}

func (s *Service) ClickCandle(messageID, sessionID string) (store.Bundle, error) {
	bundle, err := s.store.LoadBundle(messageID, sessionID)
	if err != nil {
		return store.Bundle{}, err
	}
	if err := animation.Click(&bundle.Animation); err != nil {
		return store.Bundle{}, err
	}
	bundle.Session.InteractionCount++
	bundle.Session.LastAnimationID = bundle.Animation.ID
	bundle.Snapshot = render.Render(bundle.Message, bundle.Animation)
	if err := store.ValidateBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	if err := s.store.SaveBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	return bundle, nil
}

func (s *Service) AdvanceAnimation(messageID, sessionID string) (store.Bundle, error) {
	bundle, err := s.store.LoadBundle(messageID, sessionID)
	if err != nil {
		return store.Bundle{}, err
	}
	if err := animation.Advance(&bundle.Animation); err != nil {
		return store.Bundle{}, err
	}
	bundle.Snapshot = render.Render(bundle.Message, bundle.Animation)
	if err := store.ValidateBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	if err := s.store.SaveBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	return bundle, nil
}

func (s *Service) CompleteAnimation(messageID, sessionID string) (store.Bundle, error) {
	bundle, err := s.store.LoadBundle(messageID, sessionID)
	if err != nil {
		return store.Bundle{}, err
	}
	for steps := animation.RemainingSteps(bundle.Animation); steps > 0; steps-- {
		if err := animation.Advance(&bundle.Animation); err != nil {
			return store.Bundle{}, err
		}
	}
	bundle.Snapshot = render.Render(bundle.Message, bundle.Animation)
	if err := store.ValidateBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	if err := s.store.SaveBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	return bundle, nil
}

func (s *Service) ListMemorials() ([]domain.MemorialMessage, error) { return s.store.ListMessages() }
