package domain

import "fmt"

func StageLabel(stage Stage) string {
	switch stage {
	case StageQuiet:
		return "quiet"
	case StageExpand:
		return "expanding"
	case StageStars:
		return "stars"
	default:
		return "unknown"
	}
}

func DescribeAnimation(a CandleAnimation) string {
	return fmt.Sprintf("%s generation=%d clicks=%d layers=%d particles=%d", StageLabel(a.Stage), a.Generation, a.Clicks, a.Layers, a.ParticleCount)
}
