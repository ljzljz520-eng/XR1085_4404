package service

import (
	"memorialcandle/domain"
	"memorialcandle/store"
	"testing"
)

func TestStoreServiceConstruction(t *testing.T) {
	st, err := store.Open(t.TempDir() + "/candle.db")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc, err := New(st)
	if err != nil {
		t.Fatal(err)
	}
	svc.SetSequence(0)
	if svc.Store() != st {
		t.Fatal("service did not retain store")
	}
	if err := CheckSession(domain.NewVisitorSession("v", "Visitor", "day")); err != nil {
		t.Fatal(err)
	}
}
