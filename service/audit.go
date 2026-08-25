package service

import (
	"fmt"
	"memorialcandle/animation"
	"memorialcandle/domain"
	"memorialcandle/store"
)

type InteractionSummary struct {
	MessageID  string
	SessionID  string
	Events     int
	Clicks     int
	Generation int
	Settled    bool
	Layers     int
}

func BuildInteractionSummary(bundle store.Bundle, events []domain.InteractionEvent) InteractionSummary {
	summary := InteractionSummary{MessageID: bundle.Message.ID, SessionID: bundle.Session.ID, Clicks: bundle.Animation.Clicks, Generation: bundle.Animation.Generation, Settled: animation.IsIdle(bundle.Animation), Layers: bundle.Animation.Layers}
	for _, event := range events {
		if event.MessageID == bundle.Message.ID && event.SessionID == bundle.Session.ID {
			summary.Events++
		}
	}
	return summary
}

func EventForRecord(bundle store.Bundle) domain.InteractionEvent {
	return domain.NewInteractionEvent("event-record-"+bundle.Message.ID, domain.EventRecorded, bundle.Message.ID, bundle.Session.ID, bundle.Animation.Generation, "memorial recorded")
}

func EventForClick(bundle store.Bundle) domain.InteractionEvent {
	return domain.NewInteractionEvent(fmt.Sprintf("event-click-%s-%d", bundle.Message.ID, bundle.Session.InteractionCount), domain.EventClicked, bundle.Message.ID, bundle.Session.ID, bundle.Animation.Generation, "candle clicked")
}

func EventForAdvance(bundle store.Bundle) domain.InteractionEvent {
	return domain.NewInteractionEvent(fmt.Sprintf("event-advance-%s-%d", bundle.Message.ID, bundle.Animation.Generation), domain.EventAdvanced, bundle.Message.ID, bundle.Session.ID, bundle.Animation.Generation, string(bundle.Animation.Stage))
}

func EventForSettle(bundle store.Bundle) domain.InteractionEvent {
	return domain.NewInteractionEvent("event-settle-"+bundle.Message.ID, domain.EventSettled, bundle.Message.ID, bundle.Session.ID, bundle.Animation.Generation, "single quiet state")
}

func EventForReopen(bundle store.Bundle) domain.InteractionEvent {
	return domain.NewInteractionEvent("event-reopen-"+bundle.Message.ID, domain.EventReopened, bundle.Message.ID, bundle.Session.ID, bundle.Animation.Generation, "persistence reopened")
}

func (s *Service) LifecycleEvents(messageID, sessionID string) ([]domain.InteractionEvent, error) {
	bundle, err := s.QueryDisplay(messageID, sessionID)
	if err != nil {
		return nil, err
	}
	events := []domain.InteractionEvent{EventForRecord(bundle)}
	if bundle.Session.InteractionCount > 0 {
		events = append(events, EventForClick(bundle))
	}
	if bundle.Animation.Stage == domain.StageQuiet {
		events = append(events, EventForSettle(bundle))
	}
	return events, nil
}

func (s *Service) VerifyLifecycle(messageID, sessionID string) (InteractionSummary, error) {
	bundle, err := s.QueryDisplay(messageID, sessionID)
	if err != nil {
		return InteractionSummary{}, err
	}
	events, err := s.LifecycleEvents(messageID, sessionID)
	if err != nil {
		return InteractionSummary{}, err
	}
	summary := BuildInteractionSummary(bundle, events)
	if !summary.Settled {
		return summary, fmt.Errorf("animation is not settled")
	}
	if summary.Layers != 1 {
		return summary, fmt.Errorf("animation has %d layers", summary.Layers)
	}
	return summary, nil
}
