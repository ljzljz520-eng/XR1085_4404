package main

import (
	"fmt"
	"memorialcandle/cli"
	"memorialcandle/service"
	"memorialcandle/store"
	"os"
	"path/filepath"
)

func main() {
	path := "candle.db"
	args := os.Args[1:]
	if len(args) > 0 && args[0] == "--db" {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "--db requires a path")
			os.Exit(2)
		}
		path, args = args[1], args[2:]
	}
	if dir := filepath.Dir(path); dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
	st, err := store.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer st.Close()
	svc, err := service.New(st)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	command, err := cli.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	if err := cli.Execute(command, svc, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}
