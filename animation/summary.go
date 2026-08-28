package animation

import (
	"fmt"
	"memorialcandle/domain"
)

func Summary(state domain.CandleAnimation) string {
	return fmt.Sprintf("animation %s stage=%s generation=%d clicks=%d layers=%d", state.ID, state.Stage, state.Generation, state.Clicks, state.Layers)
}

func Clone(state domain.CandleAnimation) domain.CandleAnimation { return state }
