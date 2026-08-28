package cli

import (
	"bytes"
	"memorialcandle/service"
	"memorialcandle/store"
	"testing"
)

func TestCLIParsing(t *testing.T) {
	command, err := Parse([]string{"record", "visitor", "author", "text", "day"})
	if err != nil {
		t.Fatal(err)
	}
	if command.Name != "record" || command.Validate() != nil {
		t.Fatalf("record parse failed: %#v", command)
	}
	if _, err := Parse([]string{"unknown"}); err == nil {
		t.Fatal("unknown command accepted")
	}
	help, err := Parse(nil)
	if err != nil || help.Name != "help" {
		t.Fatal("empty args should select help")
	}
}

func TestCLIHelpAndRecord(t *testing.T) {
	st, err := store.Open(t.TempDir() + "/candle.db")
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	svc, err := service.New(st)
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := Execute(Command{Name: "help"}, svc, &out); err != nil {
		t.Fatal(err)
	}
	if out.Len() == 0 {
		t.Fatal("help output empty")
	}
	out.Reset()
	if err := Execute(Command{Name: "record", Args: []string{"visitor", "author", "text", "day"}}, svc, &out); err != nil {
		t.Fatal(err)
	}
	if out.Len() == 0 {
		t.Fatal("record output empty")
	}
}
