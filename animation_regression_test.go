package memorialcandle

import (
	"memorialcandle/animation"
	"memorialcandle/domain"
	"sync"
	"testing"
)

func TestCandleAnimationReturnsToSingleIdleState(t *testing.T) {
	state := domain.NewCandleAnimation("memorial-candle", "memorial")
	ready := make(chan struct{}, 2)
	continueReads := make(chan struct{})
	animation.SetReadHook(func() { ready <- struct{}{}; <-continueReads })
	defer animation.ClearReadHook()
	var group sync.WaitGroup
	for i := 0; i < 2; i++ {
		group.Add(1)
		go func() { defer group.Done(); _ = animation.Click(&state) }()
	}
	<-ready
	<-ready
	close(continueReads)
	group.Wait()
	if state.Layers != 1 {
		t.Fatalf("expected merged click layer, got %d", state.Layers)
	}
	if err := animation.Advance(&state); err != nil {
		t.Fatal(err)
	}
	if err := animation.Advance(&state); err != nil {
		t.Fatal(err)
	}
	if !animation.IsIdle(state) {
		t.Fatalf("expected a single quiet state: %#v", state)
	}
}
