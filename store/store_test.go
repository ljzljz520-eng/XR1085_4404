package store

import (
	"memorialcandle/domain"
	"testing"
)

func TestStoreRoundTrip(t *testing.T) {
	st, err := Open(t.TempDir() + "/candle.db")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	m := domain.NewMemorialMessage("m", "Ada", "Light remains", "day")
	a := domain.NewCandleAnimation("m-candle", "m")
	d := domain.NewDisplaySnapshot(a.ID, m.Text, a)
	s := domain.NewVisitorSession("visitor-0001", "A visitor", "day")
	if err := st.SaveBundle(Bundle{Message: m, Animation: a, Snapshot: d, Session: s}); err != nil {
		t.Fatal(err)
	}
	loaded, err := st.LoadBundle("m", s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Message != m || loaded.Animation != a || loaded.Snapshot != d || loaded.Session != s {
		t.Fatalf("round trip changed data: %#v", loaded)
	}
	if _, err := st.GetMemorialMessage("missing"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
