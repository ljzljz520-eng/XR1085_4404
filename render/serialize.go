package render

import (
	"encoding/json"
	"fmt"
	"memorialcandle/domain"
)

type WireSnapshot struct {
	ID         string `json:"id"`
	Text       string `json:"text"`
	Stage      string `json:"stage"`
	Particles  int    `json:"particles"`
	Quiet      bool   `json:"quiet"`
	Generation int    `json:"generation"`
	Layers     int    `json:"layers"`
}

func ToWire(message domain.MemorialMessage, animation domain.CandleAnimation) WireSnapshot {
	return WireSnapshot{ID: animation.ID, Text: EscapeText(message.Text), Stage: string(animation.Stage), Particles: animation.ParticleCount, Quiet: animation.Quiet, Generation: animation.Generation, Layers: animation.Layers}
}

func EncodeWire(message domain.MemorialMessage, animation domain.CandleAnimation) ([]byte, error) {
	return json.Marshal(ToWire(message, animation))
}

func DecodeWire(data []byte) (WireSnapshot, error) {
	var snapshot WireSnapshot
	err := json.Unmarshal(data, &snapshot)
	return snapshot, err
}

func (w WireSnapshot) String() string {
	return fmt.Sprintf("%s [%s] particles=%d quiet=%t generation=%d layers=%d", w.Text, w.Stage, w.Particles, w.Quiet, w.Generation, w.Layers)
}

func StageColor(stage domain.Stage) string {
	switch stage {
	case domain.StageQuiet:
		return "amber"
	case domain.StageExpand:
		return "gold"
	case domain.StageStars:
		return "silver"
	default:
		return "neutral"
	}
}

func FrameLabel(index int, stage domain.Stage) string {
	return fmt.Sprintf("frame-%02d-%s", index, stage)
}
