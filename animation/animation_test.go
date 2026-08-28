package animation

import (
	"memorialcandle/domain"
	"testing"
)

func TestAnimationLifecycle(t *testing.T) {
	ClearReadHook()
	a := domain.NewCandleAnimation("a", "m")
	if err := Click(&a); err != nil {
		t.Fatal(err)
	}
	if a.Stage != domain.StageExpand || a.Quiet {
		t.Fatalf("unexpected expanding state: %#v", a)
	}
	if err := Advance(&a); err != nil {
		t.Fatal(err)
	}
	if a.Stage != domain.StageStars {
		t.Fatalf("unexpected stars state: %#v", a)
	}
	if err := Advance(&a); err != nil {
		t.Fatal(err)
	}
	if !IsIdle(a) || a.Generation != 2 {
		t.Fatalf("unexpected quiet state: %#v", a)
	}
}

func TestSessionInteractionCount(t *testing.T) {
	a := domain.NewCandleAnimation("a", "m")
	if err := Click(&a); err != nil {
		t.Fatal(err)
	}
	if err := Click(&a); err != nil {
		t.Fatal(err)
	}
	if a.Clicks != 2 || a.Layers != 1 {
		t.Fatalf("clicks were not merged: %#v", a)
	}
}

func TestAnimationMetrics(t *testing.T) {
	if ParticleProfile(domain.StageQuiet, 4) != 1 || ParticleProfile(domain.StageStars, 2) != 84 {
		t.Fatal("particle profile mismatch")
	}
	if RemainingSteps(domain.CandleAnimation{Stage: domain.StageExpand}) != 2 {
		t.Fatal("remaining step count mismatch")
	}
}
