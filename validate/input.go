package validate

import (
	"errors"
	"strings"
	"unicode/utf8"

	"memorialcandle/domain"
)

func Message(author, text string) error {
	if strings.TrimSpace(author) == "" {
		return errors.New("author cannot be blank")
	}
	if strings.TrimSpace(text) == "" {
		return errors.New("text cannot be blank")
	}
	if !utf8.ValidString(text) {
		return errors.New("text must be valid utf-8")
	}
	if len([]rune(strings.TrimSpace(text))) > domain.MaxMessageLen {
		return errors.New("text exceeds 280 characters")
	}
	return nil
}

func Visitor(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("visitor name cannot be blank")
	}
	if len([]rune(strings.TrimSpace(name))) > 80 {
		return errors.New("visitor name is too long")
	}
	return nil
}

func Sequence(sequence int) error {
	if sequence < 1 || sequence > 9999 {
		return errors.New("sequence must be between 1 and 9999")
	}
	return nil
}
