package service

import (
	"memorialcandle/domain"
	"memorialcandle/store"
	"testing"
)

func newWorkflowService(t *testing.T) *Service {
	t.Helper()
	st, err := store.Open(t.TempDir() + "/candle.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	svc, err := New(st)
	if err != nil {
		t.Fatal(err)
	}
	return svc
}

func TestWorkflowRecordMemorial(t *testing.T) {
	svc := newWorkflowService(t)
	bundle, err := svc.RecordMemorial("Visitor One", "Ada", "A quiet remembrance", "day-1")
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Message.ID != "memorial-0001" || bundle.Animation.Stage != domain.StageQuiet || bundle.Snapshot.ParticleCount != 1 {
		t.Fatalf("record workflow failed: %#v", bundle)
	}
}

func TestWorkflowQueryDisplay(t *testing.T) {
	svc := newWorkflowService(t)
	created, err := svc.RecordMemorial("Visitor One", "Ada", "A quiet remembrance", "day-1")
	if err != nil {
		t.Fatal(err)
	}
	queried, err := svc.QueryDisplay(created.Message.ID, created.Session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if queried.Snapshot.VisibleText != created.Message.Text || queried.Snapshot.Stage != domain.StageQuiet {
		t.Fatalf("query workflow failed: %#v", queried)
	}
}

func TestWorkflowInteractiveCandle(t *testing.T) {
	svc := newWorkflowService(t)
	created, err := svc.RecordMemorial("Visitor One", "Ada", "A quiet remembrance", "day-1")
	if err != nil {
		t.Fatal(err)
	}
	clicked, err := svc.ClickCandle(created.Message.ID, created.Session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if clicked.Animation.Stage != domain.StageExpand || clicked.Snapshot.Quiet {
		t.Fatalf("click did not expand: %#v", clicked)
	}
	completed, err := svc.CompleteAnimation(created.Message.ID, created.Session.ID)
	if err != nil {
		t.Fatal(err)
	}
	if completed.Animation.Stage != domain.StageQuiet || completed.Animation.Layers != 1 || !completed.Snapshot.Quiet {
		t.Fatalf("animation did not settle: %#v", completed)
	}
}
