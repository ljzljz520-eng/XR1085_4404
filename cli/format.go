package cli

import (
	"fmt"
	"strings"
)

type Option struct {
	Name        string
	Value       string
	Required    bool
	Description string
}

func Options() []Option {
	return []Option{{Name: "--db", Value: "path", Required: false, Description: "bbolt database path"}, {Name: "record", Value: "visitor author text started-at", Required: true, Description: "store a memorial"}, {Name: "show", Value: "message-id session-id", Required: true, Description: "display a memorial"}, {Name: "click", Value: "message-id session-id", Required: true, Description: "trigger candle animation"}, {Name: "advance", Value: "message-id session-id", Required: true, Description: "advance one animation stage"}, {Name: "complete", Value: "message-id session-id", Required: true, Description: "settle the candle"}, {Name: "list", Value: "", Required: false, Description: "list memorials"}}
}

func Usage() string {
	lines := make([]string, 0)
	for _, option := range Options() {
		suffix := ""
		if option.Value != "" {
			suffix = " " + option.Value
		}
		lines = append(lines, fmt.Sprintf("%-8s %-28s %s", option.Name, suffix, option.Description))
	}
	return strings.Join(lines, "\n")
}

func IsMutating(command Command) bool {
	switch command.Name {
	case "record", "click", "advance", "complete":
		return true
	default:
		return false
	}
}

func ParseDatabase(args []string) (string, []string, error) {
	path := "candle.db"
	remaining := append([]string(nil), args...)
	if len(remaining) >= 2 && remaining[0] == "--db" {
		path = remaining[1]
		remaining = remaining[2:]
	}
	if path == "" {
		return "", nil, fmt.Errorf("database path is empty")
	}
	return path, remaining, nil
}
