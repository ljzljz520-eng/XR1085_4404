package animation

import "memorialcandle/domain"

func ParticleProfile(stage domain.Stage, layers int) int {
	if layers < 1 {
		layers = 1
	}
	switch stage {
	case domain.StageQuiet:
		return 1
	case domain.StageExpand:
		return 12 * layers
	case domain.StageStars:
		return 42 * layers
	default:
		return 0
	}
}

func StageSequence() []domain.Stage {
	return []domain.Stage{domain.StageQuiet, domain.StageExpand, domain.StageStars}
}

func ExpectedQuietGeneration(initial int, completed int) int {
	if initial < 1 {
		initial = 1
	}
	if completed < 0 {
		completed = 0
	}
	return initial + completed
}
