package domain

import (
	"errors"
	"strings"
)

type Stage string

const (
	StageQuiet    Stage = "quiet"
	StageExpand   Stage = "expanding"
	StageStars    Stage = "stars"
	MaxMessageLen       = 280
	MaxClicks           = 64
)

type MemorialMessage struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	Active    bool   `json:"active"`
}

type CandleAnimation struct {
	ID            string `json:"id"`
	MessageID     string `json:"message_id"`
	Stage         Stage  `json:"stage"`
	Generation    int    `json:"generation"`
	Clicks        int    `json:"clicks"`
	Layers        int    `json:"layers"`
	ParticleCount int    `json:"particle_count"`
	Quiet         bool   `json:"quiet"`
}

type DisplaySnapshot struct {
	AnimationID   string `json:"animation_id"`
	VisibleText   string `json:"visible_text"`
	ParticleCount int    `json:"particle_count"`
	Stage         Stage  `json:"stage"`
	Quiet         bool   `json:"quiet"`
	Generation    int    `json:"generation"`
}

type VisitorSession struct {
	ID               string `json:"id"`
	VisitorName      string `json:"visitor_name"`
	StartedAt        string `json:"started_at"`
	InteractionCount int    `json:"interaction_count"`
	LastAnimationID  string `json:"last_animation_id"`
}

func NewMemorialMessage(id, author, text, createdAt string) MemorialMessage {
	return MemorialMessage{ID: id, Author: strings.TrimSpace(author), Text: strings.TrimSpace(text), CreatedAt: createdAt, Active: true}
}

func NewVisitorSession(id, visitor, startedAt string) VisitorSession {
	return VisitorSession{ID: id, VisitorName: strings.TrimSpace(visitor), StartedAt: startedAt}
}

func NewCandleAnimation(id, messageID string) CandleAnimation {
	return CandleAnimation{ID: id, MessageID: messageID, Stage: StageQuiet, Generation: 1, Layers: 1, ParticleCount: 1, Quiet: true}
}

func NewDisplaySnapshot(animationID, text string, animation CandleAnimation) DisplaySnapshot {
	return DisplaySnapshot{AnimationID: animationID, VisibleText: text, ParticleCount: animation.ParticleCount, Stage: animation.Stage, Quiet: animation.Quiet, Generation: animation.Generation}
}

func (m MemorialMessage) Validate() error {
	if strings.TrimSpace(m.ID) == "" {
		return errors.New("message id is required")
	}
	if strings.TrimSpace(m.Author) == "" {
		return errors.New("author is required")
	}
	if strings.TrimSpace(m.Text) == "" {
		return errors.New("message text is required")
	}
	if len([]rune(m.Text)) > MaxMessageLen {
		return errors.New("message text is too long")
	}
	return nil
}

func (s VisitorSession) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.VisitorName) == "" {
		return errors.New("session identity is required")
	}
	if s.InteractionCount < 0 {
		return errors.New("interaction count cannot be negative")
	}
	return nil
}

func (a CandleAnimation) Validate() error {
	if strings.TrimSpace(a.ID) == "" || strings.TrimSpace(a.MessageID) == "" {
		return errors.New("animation links are required")
	}
	if a.Stage != StageQuiet && a.Stage != StageExpand && a.Stage != StageStars {
		return errors.New("unknown animation stage")
	}
	if a.Generation < 1 || a.Clicks < 0 || a.Clicks > MaxClicks {
		return errors.New("animation counters are invalid")
	}
	if a.Layers < 1 {
		return errors.New("animation must have a layer")
	}
	if a.Stage == StageQuiet && (!a.Quiet || a.Layers != 1 || a.ParticleCount != 1) {
		return errors.New("quiet animation must have one layer")
	}
	return nil
}

func (d DisplaySnapshot) Validate() error {
	if strings.TrimSpace(d.AnimationID) == "" {
		return errors.New("snapshot animation id is required")
	}
	if d.ParticleCount < 1 || d.Generation < 1 {
		return errors.New("snapshot counters are invalid")
	}
	if d.Stage == StageQuiet && (!d.Quiet || d.ParticleCount != 1) {
		return errors.New("quiet snapshot must have one particle")
	}
	return nil
}
