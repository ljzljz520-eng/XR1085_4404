package store

import (
	"memorialcandle/domain"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/candle.db"
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	m := domain.NewMemorialMessage("m", "Ada", "Remembered in light", "day")
	a := domain.NewCandleAnimation("m-candle", "m")
	d := domain.NewDisplaySnapshot(a.ID, m.Text, a)
	s := domain.NewVisitorSession("visitor-0001", "Visitor", "day")
	if err := st.SaveBundle(Bundle{Message: m, Animation: a, Snapshot: d, Session: s}); err != nil {
		t.Fatal(err)
	}
	if err := st.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	bundle, err := reopened.LoadBundle(m.ID, s.ID)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Message.Text != m.Text || bundle.Animation.Stage != domain.StageQuiet || !bundle.Snapshot.Quiet {
		t.Fatalf("persistence did not restore quiet state: %#v", bundle)
	}
}
