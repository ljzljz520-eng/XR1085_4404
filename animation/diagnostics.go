package animation

import (
	"fmt"
	"memorialcandle/domain"
)

type Diagnostics struct {
	Stage     domain.Stage
	Valid     bool
	Idle      bool
	Steps     int
	Particles int
	Layers    int
	Message   string
}

func Diagnose(state domain.CandleAnimation) Diagnostics {
	diagnostics := Diagnostics{Stage: state.Stage, Idle: IsIdle(state), Steps: RemainingSteps(state), Particles: state.ParticleCount, Layers: state.Layers}
	if err := state.Validate(); err != nil {
		diagnostics.Message = err.Error()
		diagnostics.Valid = false
	} else {
		diagnostics.Valid = true
		diagnostics.Message = "animation state is valid"
	}
	return diagnostics
}
func (d Diagnostics) String() string {
	return fmt.Sprintf("stage=%s valid=%t idle=%t steps=%d particles=%d layers=%d message=%s", d.Stage, d.Valid, d.Idle, d.Steps, d.Particles, d.Layers, d.Message)
}
func (d Diagnostics) Recoverable() bool {
	return d.Valid && d.Layers == 1 && (d.Stage == domain.StageExpand || d.Stage == domain.StageStars)
}
func (d Diagnostics) RequiresReset() bool { return !d.Valid || d.Layers != 1 }

func Compare(left, right domain.CandleAnimation) []string {
	differences := make([]string, 0)
	if left.Stage != right.Stage {
		differences = append(differences, "stage")
	}
	if left.Generation != right.Generation {
		differences = append(differences, "generation")
	}
	if left.Clicks != right.Clicks {
		differences = append(differences, "clicks")
	}
	if left.Layers != right.Layers {
		differences = append(differences, "layers")
	}
	if left.ParticleCount != right.ParticleCount {
		differences = append(differences, "particles")
	}
	if left.Quiet != right.Quiet {
		differences = append(differences, "quiet")
	}
	return differences
}

func Normalize(state *domain.CandleAnimation) error {
	if state == nil {
		return fmt.Errorf("animation is nil")
	}
	if state.Layers < 1 {
		state.Layers = 1
	}
	if state.Generation < 1 {
		state.Generation = 1
	}
	if state.Clicks < 0 {
		state.Clicks = 0
	}
	if state.Clicks > domain.MaxClicks {
		state.Clicks = domain.MaxClicks
	}
	if state.Stage == domain.StageQuiet {
		state.Quiet = true
		state.ParticleCount = 1
		state.Layers = 1
	}
	return state.Validate()
}
