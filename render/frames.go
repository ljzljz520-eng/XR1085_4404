package render

import (
	"fmt"
	"memorialcandle/animation"
	"memorialcandle/domain"
)

func RenderTimeline(message domain.MemorialMessage, state domain.CandleAnimation) ([]string, error) {
	frames, err := animation.BuildTimeline(state)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0, len(frames))
	for _, frame := range frames {
		result = append(result, fmt.Sprintf("%s %s particles=%d layers=%d", message.ID, frame.Stage, frame.Particles, frame.Layers))
	}
	return result, nil
}

func AccessibilityText(message domain.MemorialMessage, state domain.CandleAnimation) string {
	if state.Quiet {
		return "Quiet candle: " + message.Text
	}
	return "Candle animation " + string(state.Stage) + ": " + message.Text
}

func ParticleCountFor(stage domain.Stage, layers int) int { return len(ParticlePoints(stage, layers)) }
