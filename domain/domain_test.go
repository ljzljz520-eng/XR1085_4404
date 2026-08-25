package domain

import "testing"

func TestEntityConstructors(t *testing.T) {
	m := NewMemorialMessage(MessageID(2), " Ada ", " Light remains ", "day-2")
	if m.Author != "Ada" || m.Text != "Light remains" || !m.Active {
		t.Fatalf("unexpected message: %#v", m)
	}
	s := NewVisitorSession(SessionID(2), " Visitor ", "day-2")
	a := NewCandleAnimation(AnimationID(m.ID), m.ID)
	d := NewDisplaySnapshot(a.ID, m.Text, a)
	for name, err := range map[string]error{"message": m.Validate(), "session": s.Validate(), "animation": a.Validate(), "snapshot": d.Validate()} {
		if err != nil {
			t.Fatalf("%s invalid: %v", name, err)
		}
	}
}

func TestIdentifiersAndLabels(t *testing.T) {
	if MessageID(NormalizeSequence(0)) != "memorial-0001" {
		t.Fatal("sequence normalization failed")
	}
	if SessionID(7) != "visitor-0007" || AnimationID("memorial-0007") != "memorial-0007-candle" {
		t.Fatal("identifier format failed")
	}
	if StageLabel(StageStars) != "stars" || StageLabel(Stage("bad")) != "unknown" {
		t.Fatal("stage label failed")
	}
}
