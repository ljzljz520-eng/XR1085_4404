package animation

import (
	"errors"
	"runtime"
	"sync"

	"memorialcandle/domain"
)

var readHookMu sync.RWMutex
var readHook func()
var clickWriteMu sync.Mutex

func SetReadHook(hook func()) { readHookMu.Lock(); readHook = hook; readHookMu.Unlock() }
func ClearReadHook()          { SetReadHook(nil) }
func invokeReadHook() {
	readHookMu.RLock()
	hook := readHook
	readHookMu.RUnlock()
	if hook != nil {
		hook()
	}
}

func Click(state *domain.CandleAnimation) error {
	if state == nil {
		return errors.New("animation is nil")
	}
	if err := state.Validate(); err != nil {
		return err
	}
	observedStage := state.Stage
	observedClicks := state.Clicks
	invokeReadHook()
	runtime.Gosched()
	clickWriteMu.Lock()
	defer clickWriteMu.Unlock()
	if observedClicks >= domain.MaxClicks {
		return errors.New("click limit reached")
	}
	if observedStage == domain.StageQuiet {
		if state.Stage == domain.StageQuiet {
			state.Stage = domain.StageExpand
			state.Quiet = false
			state.ParticleCount = 12
			state.Layers = 1
			state.Clicks = observedClicks + 1
		} else {
			// A concurrent click already moved this layer off quiet.
			// Brighten the single expanding layer instead of spawning a
			// second one — multiple layers can never settle to quiet.
			state.Stage = domain.StageExpand
			state.Quiet = false
			state.ParticleCount = 24
			state.Layers = 1
			state.Clicks = observedClicks + 1
		}
		return nil
	}
	if state.Stage == domain.StageQuiet {
		state.Stage = domain.StageExpand
		state.Quiet = false
		state.ParticleCount = 12
		state.Layers = 1
	}
	state.Clicks = observedClicks + 1
	return nil
}

func Advance(state *domain.CandleAnimation) error {
	if state == nil {
		return errors.New("animation is nil")
	}
	if err := state.Validate(); err != nil {
		return err
	}
	switch state.Stage {
	case domain.StageQuiet:
		return nil
	case domain.StageExpand:
		state.Stage = domain.StageStars
		state.ParticleCount = 36 + state.Layers*6
		state.Quiet = false
	case domain.StageStars:
		if state.Layers != 1 {
			return errors.New("multiple active layers cannot settle")
		}
		state.Stage = domain.StageQuiet
		state.Generation++
		state.ParticleCount = 1
		state.Quiet = true
		state.Clicks = 0
	default:
		return errors.New("cannot advance unknown stage")
	}
	return nil
}

func Quiet(state *domain.CandleAnimation) error {
	if state == nil {
		return errors.New("animation is nil")
	}
	if state.Layers != 1 {
		return errors.New("cannot quiet multiple layers")
	}
	state.Stage = domain.StageQuiet
	state.ParticleCount = 1
	state.Quiet = true
	state.Clicks = 0
	return nil
}

func IsIdle(state domain.CandleAnimation) bool {
	return state.Stage == domain.StageQuiet && state.Quiet && state.Layers == 1 && state.ParticleCount == 1
}
func RemainingSteps(state domain.CandleAnimation) int {
	switch state.Stage {
	case domain.StageQuiet:
		return 0
	case domain.StageExpand:
		return 2
	case domain.StageStars:
		return 1
	default:
		return -1
	}
}
