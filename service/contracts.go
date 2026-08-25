package service

import (
	"errors"
	"memorialcandle/domain"
	"memorialcandle/store"
	"memorialcandle/validate"
)

func ValidateBundle(bundle store.Bundle) error { return store.ValidateBundle(bundle) }

func CheckSession(session domain.VisitorSession) error {
	if err := validate.SessionEntity(session); err != nil {
		return err
	}
	if session.InteractionCount > domain.MaxClicks {
		return errors.New("session interaction count is too high")
	}
	return nil
}
