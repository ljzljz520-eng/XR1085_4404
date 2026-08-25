package render

import (
	"memorialcandle/domain"
	"testing"
)

func TestRenderSnapshot(t *testing.T) {
	m := domain.NewMemorialMessage("m", "Ada", "Light remains", "d")
	a := domain.NewCandleAnimation("m-candle", "m")
	d := Render(m, a)
	if d.VisibleText != "Light remains" || d.ParticleCount != 1 || !d.Quiet {
		t.Fatalf("unexpected quiet render: %#v", d)
	}
	if got := RenderLines(m, a); len(got) != 4 {
		t.Fatalf("unexpected lines: %#v", got)
	}
	if Header() == "" || Help() == "" {
		t.Fatal("render labels are empty")
	}
}
