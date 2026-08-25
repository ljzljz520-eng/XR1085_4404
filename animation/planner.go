package animation

import (
	"memorialcandle/domain"
)

type PlanStep struct {
	Name      string
	From      domain.Stage
	To        domain.Stage
	Particles int
}

func Plan() []PlanStep {
	return []PlanStep{{Name: "ignite", From: domain.StageQuiet, To: domain.StageExpand, Particles: 12}, {Name: "scatter", From: domain.StageExpand, To: domain.StageStars, Particles: 42}, {Name: "settle", From: domain.StageStars, To: domain.StageQuiet, Particles: 1}}
}

func StepFor(stage domain.Stage) PlanStep {
	for _, step := range Plan() {
		if step.From == stage {
			return step
		}
	}
	return PlanStep{Name: "unknown", From: stage, To: stage, Particles: 0}
}

func CanAdvance(state domain.CandleAnimation) bool {
	return state.Stage == domain.StageExpand || state.Stage == domain.StageStars
}

func CompletionGeneration(state domain.CandleAnimation) int {
	if state.Stage == domain.StageQuiet {
		return state.Generation
	}
	return state.Generation + 1
}
