package animation

import (
	"errors"
	"fmt"
	"memorialcandle/domain"
)

type Frame struct {
	Index      int
	Stage      domain.Stage
	Particles  int
	Quiet      bool
	Generation int
	Layers     int
}

func FrameFromState(index int, state domain.CandleAnimation) Frame {
	return Frame{Index: index, Stage: state.Stage, Particles: state.ParticleCount, Quiet: state.Quiet, Generation: state.Generation, Layers: state.Layers}
}

func BuildTimeline(initial domain.CandleAnimation) ([]Frame, error) {
	if err := initial.Validate(); err != nil {
		return nil, err
	}
	frames := []Frame{FrameFromState(0, initial)}
	state := initial
	for index := 1; index <= 2; index++ {
		if err := Advance(&state); err != nil {
			return nil, err
		}
		frames = append(frames, FrameFromState(index, state))
	}
	return frames, nil
}

func ValidateTimeline(frames []Frame) error {
	if len(frames) < 3 {
		return errors.New("timeline requires three frames")
	}
	if frames[0].Stage != domain.StageQuiet || frames[1].Stage != domain.StageExpand || frames[2].Stage != domain.StageStars {
		return errors.New("timeline stages are not ordered")
	}
	for _, frame := range frames {
		if frame.Index < 0 || frame.Particles < 1 || frame.Layers < 1 {
			return errors.New("timeline frame is invalid")
		}
	}
	return nil
}

func MergeClicks(state domain.CandleAnimation, clicks int) (domain.CandleAnimation, error) {
	if clicks < 1 || clicks > domain.MaxClicks {
		return state, errors.New("click count out of range")
	}
	for i := 0; i < clicks; i++ {
		if err := Click(&state); err != nil {
			return state, err
		}
	}
	return state, nil
}

func FrameDigest(frame Frame) string {
	return fmt.Sprintf("%d:%s:%d:%t:%d:%d", frame.Index, frame.Stage, frame.Particles, frame.Quiet, frame.Generation, frame.Layers)
}

func TimelineDigest(frames []Frame) string {
	digest := ""
	for i, frame := range frames {
		if i > 0 {
			digest += "|"
		}
		digest += FrameDigest(frame)
	}
	return digest
}
