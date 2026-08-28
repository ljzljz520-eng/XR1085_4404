package validate

import (
	"memorialcandle/domain"
	"testing"
)

func TestValidateMemorialMessage(t *testing.T) {
	if err := Message("Ada", "A steady light"); err != nil {
		t.Fatal(err)
	}
	if Message("", "text") == nil || Message("Ada", "") == nil {
		t.Fatal("blank input accepted")
	}
	if Sequence(0) == nil || Sequence(10000) == nil {
		t.Fatal("bad sequence accepted")
	}
}

func TestValidateAnimation(t *testing.T) {
	a := domain.NewCandleAnimation("a", "m")
	if err := AnimationEntity(a); err != nil {
		t.Fatal(err)
	}
	a.Stage = domain.Stage("broken")
	if AnimationEntity(a) == nil {
		t.Fatal("broken stage accepted")
	}
}

func TestValidateLinksAndQuietState(t *testing.T) {
	m := domain.NewMemorialMessage("m", "A", "Text", "d")
	a := domain.NewCandleAnimation("m-candle", "m")
	d := domain.NewDisplaySnapshot(a.ID, m.Text, a)
	if err := Linked(m, a, d); err != nil {
		t.Fatal(err)
	}
	if err := QuietState(a, d); err != nil {
		t.Fatal(err)
	}
}
