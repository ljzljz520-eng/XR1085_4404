package domain

import (
	"errors"
	"fmt"
)

type EventKind string

const (
	EventRecorded EventKind = "recorded"
	EventClicked  EventKind = "clicked"
	EventAdvanced EventKind = "advanced"
	EventSettled  EventKind = "settled"
	EventReopened EventKind = "reopened"
)

type InteractionEvent struct {
	ID         string    `json:"id"`
	Kind       EventKind `json:"kind"`
	MessageID  string    `json:"message_id"`
	SessionID  string    `json:"session_id"`
	Generation int       `json:"generation"`
	Detail     string    `json:"detail"`
}

func NewInteractionEvent(id string, kind EventKind, messageID, sessionID string, generation int, detail string) InteractionEvent {
	return InteractionEvent{ID: id, Kind: kind, MessageID: messageID, SessionID: sessionID, Generation: generation, Detail: detail}
}

func (e InteractionEvent) Validate() error {
	if e.ID == "" || e.MessageID == "" || e.SessionID == "" {
		return errors.New("event identity is required")
	}
	if e.Generation < 1 {
		return errors.New("event generation is invalid")
	}
	switch e.Kind {
	case EventRecorded, EventClicked, EventAdvanced, EventSettled, EventReopened:
	default:
		return errors.New("event kind is invalid")
	}
	return nil
}

func (e InteractionEvent) Summary() string {
	return fmt.Sprintf("%s %s generation=%d %s", e.ID, e.Kind, e.Generation, e.Detail)
}

type EventLog struct{ events []InteractionEvent }

func NewEventLog() *EventLog { return &EventLog{events: make([]InteractionEvent, 0)} }
func (l *EventLog) Append(event InteractionEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	l.events = append(l.events, event)
	return nil
}
func (l *EventLog) Len() int {
	if l == nil {
		return 0
	}
	return len(l.events)
}
func (l *EventLog) All() []InteractionEvent {
	if l == nil {
		return nil
	}
	return append([]InteractionEvent(nil), l.events...)
}
func (l *EventLog) ForMessage(id string) []InteractionEvent {
	result := make([]InteractionEvent, 0)
	if l == nil {
		return result
	}
	for _, e := range l.events {
		if e.MessageID == id {
			result = append(result, e)
		}
	}
	return result
}
func (l *EventLog) Last() (InteractionEvent, bool) {
	if l == nil || len(l.events) == 0 {
		return InteractionEvent{}, false
	}
	return l.events[len(l.events)-1], true
}
