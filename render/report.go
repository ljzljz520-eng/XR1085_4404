package render

import (
	"fmt"
	"memorialcandle/domain"
	"sort"
	"strings"
)

type Report struct {
	Title          string
	MessageCount   int
	ActiveCount    int
	QuietCount     int
	TotalParticles int
	Lines          []string
}

func BuildReport(messages []domain.MemorialMessage, animations []domain.CandleAnimation) Report {
	report := Report{Title: "Memorial candlelight report", MessageCount: len(messages), Lines: make([]string, 0, len(messages))}
	for _, message := range messages {
		if message.Active {
			report.ActiveCount++
		}
		report.Lines = append(report.Lines, message.ID+" "+message.Author+": "+message.Text)
	}
	for _, animation := range animations {
		report.TotalParticles += animation.ParticleCount
		if animation.Quiet {
			report.QuietCount++
		}
	}
	sort.Strings(report.Lines)
	return report
}

func (r Report) Text() string {
	return fmt.Sprintf("%s messages=%d active=%d quiet=%d particles=%d", r.Title, r.MessageCount, r.ActiveCount, r.QuietCount, r.TotalParticles)
}

func RenderTable(report Report) string {
	lines := []string{report.Text()}
	lines = append(lines, report.Lines...)
	return strings.Join(lines, "\n")
}

func EscapeText(text string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(text)
}

func StageClass(stage domain.Stage) string {
	switch stage {
	case domain.StageQuiet:
		return "stage-quiet"
	case domain.StageExpand:
		return "stage-expanding"
	case domain.StageStars:
		return "stage-stars"
	default:
		return "stage-unknown"
	}
}

func ParticlePoints(stage domain.Stage, layers int) []string {
	count := 1
	if stage == domain.StageExpand {
		count = 12
	}
	if stage == domain.StageStars {
		count = 42
	}
	if layers < 1 {
		layers = 1
	}
	points := make([]string, 0, count*layers)
	for layer := 0; layer < layers; layer++ {
		for i := 0; i < count; i++ {
			points = append(points, fmt.Sprintf("%d:%d", layer, i))
		}
	}
	return points
}
