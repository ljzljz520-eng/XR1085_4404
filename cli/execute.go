package cli

import (
	"fmt"
	"io"
	"memorialcandle/render"
	"memorialcandle/service"
)

func Execute(command Command, svc *service.Service, out io.Writer) error {
	if err := command.Validate(); err != nil {
		return err
	}
	if command.Name == "help" {
		_, err := fmt.Fprintln(out, render.Header()+" commands: "+render.Help())
		return err
	}
	if command.Name == "list" {
		messages, err := svc.ListMemorials()
		if err != nil {
			return err
		}
		for _, message := range messages {
			if _, err := fmt.Fprintln(out, message.ID+" "+message.Author+": "+message.Text); err != nil {
				return err
			}
		}
		return nil
	}
	var err error
	switch command.Name {
	case "record":
		bundle, e := svc.RecordMemorial(command.Args[0], command.Args[1], command.Args[2], command.Args[3])
		if e != nil {
			return e
		}
		return WriteBundle(out, bundle)
	case "show":
		bundle, e := svc.QueryDisplay(command.Args[0], command.Args[1])
		if e != nil {
			return e
		}
		return WriteBundle(out, bundle)
	case "click":
		bundle, e := svc.ClickCandle(command.Args[0], command.Args[1])
		if e != nil {
			return e
		}
		return WriteBundle(out, bundle)
	case "advance":
		bundle, e := svc.AdvanceAnimation(command.Args[0], command.Args[1])
		if e != nil {
			return e
		}
		return WriteBundle(out, bundle)
	case "complete":
		bundle, e := svc.CompleteAnimation(command.Args[0], command.Args[1])
		if e != nil {
			return e
		}
		return WriteBundle(out, bundle)
	}
	return err
}
