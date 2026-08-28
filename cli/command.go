package cli

import (
	"errors"
	"fmt"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

func Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{Name: "help"}, nil
	}
	name := strings.ToLower(strings.TrimSpace(args[0]))
	switch name {
	case "help", "record", "show", "click", "advance", "complete", "list":
		return Command{Name: name, Args: append([]string(nil), args[1:]...)}, nil
	default:
		return Command{}, fmt.Errorf("unknown command %q", name)
	}
}

func (c Command) Validate() error {
	switch c.Name {
	case "record":
		if len(c.Args) != 4 {
			return errors.New("record requires visitor author text started-at")
		}
	case "show", "click", "advance", "complete":
		if len(c.Args) != 2 {
			return fmt.Errorf("%s requires message-id session-id", c.Name)
		}
	case "help", "list":
		if len(c.Args) != 0 {
			return fmt.Errorf("%s takes no arguments", c.Name)
		}
	default:
		return errors.New("command is empty")
	}
	return nil
}
