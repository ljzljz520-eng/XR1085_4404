package service

import (
	"errors"
	"memorialcandle/domain"
	"memorialcandle/render"
	"memorialcandle/store"
)

func (s *Service) Search(term string, limit int) ([]domain.MemorialMessage, error) {
	return s.store.QueryMessages(store.MessageQuery{Term: term, ActiveOnly: true, Limit: limit})
}

func (s *Service) Report() (render.Report, error) {
	messages, err := s.store.ListMessages()
	if err != nil {
		return render.Report{}, err
	}
	animations, err := s.store.ListAnimations()
	if err != nil {
		return render.Report{}, err
	}
	return render.BuildReport(messages, animations), nil
}

func (s *Service) Deactivate(messageID string) error {
	message, err := s.store.GetMemorialMessage(messageID)
	if err != nil {
		return err
	}
	message.Active = false
	return s.store.PutMemorialMessage(message)
}

func (s *Service) Activate(messageID string) error {
	message, err := s.store.GetMemorialMessage(messageID)
	if err != nil {
		return err
	}
	message.Active = true
	return s.store.PutMemorialMessage(message)
}

func (s *Service) SetQuiet(bundle store.Bundle) (store.Bundle, error) {
	if bundle.Animation.Layers != 1 {
		return store.Bundle{}, errors.New("cannot settle layered animation")
	}
	bundle.Animation.Stage = domain.StageQuiet
	bundle.Animation.ParticleCount = 1
	bundle.Animation.Quiet = true
	bundle.Animation.Clicks = 0
	bundle.Snapshot = render.Render(bundle.Message, bundle.Animation)
	if err := s.store.SaveBundle(bundle); err != nil {
		return store.Bundle{}, err
	}
	return bundle, nil
}

func (s *Service) SnapshotLines(messageID, sessionID string) ([]string, error) {
	bundle, err := s.QueryDisplay(messageID, sessionID)
	if err != nil {
		return nil, err
	}
	return render.RenderLines(bundle.Message, bundle.Animation), nil
}

func (s *Service) Export() (store.ExportData, error) { return s.store.Export() }
