package render

import (
	"memorialcandle/domain"
	"strings"
)

func Render(message domain.MemorialMessage, animation domain.CandleAnimation) domain.DisplaySnapshot {
	text := strings.TrimSpace(message.Text)
	if animation.Stage != domain.StageQuiet {
		text = strings.TrimSpace(message.Author) + ": " + text
	}
	return domain.NewDisplaySnapshot(animation.ID, text, animation)
}

func RenderLines(message domain.MemorialMessage, animation domain.CandleAnimation) []string {
	snapshot := Render(message, animation)
	return []string{snapshot.VisibleText, string(snapshot.Stage), formatParticles(snapshot.ParticleCount), formatQuiet(snapshot.Quiet)}
}
