package validate

import (
	"fmt"
	"memorialcandle/domain"
	"strings"
)

type RuleResult struct {
	Name   string
	Passed bool
	Reason string
}

func RunMessageRules(message domain.MemorialMessage) []RuleResult {
	rules := []RuleResult{
		{Name: "identity", Passed: strings.TrimSpace(message.ID) != "", Reason: "message id"},
		{Name: "author", Passed: strings.TrimSpace(message.Author) != "", Reason: "author"},
		{Name: "text", Passed: strings.TrimSpace(message.Text) != "", Reason: "text"},
		{Name: "length", Passed: len([]rune(message.Text)) <= domain.MaxMessageLen, Reason: "bounded text"},
		{Name: "active", Passed: message.Active, Reason: "active record"},
	}
	return rules
}

func ValidateMessageBatch(messages []domain.MemorialMessage) error {
	seen := make(map[string]struct{}, len(messages))
	for _, message := range messages {
		if _, exists := seen[message.ID]; exists {
			return fmt.Errorf("duplicate message %s", message.ID)
		}
		seen[message.ID] = struct{}{}
		if err := MessageEntity(message); err != nil {
			return err
		}
	}
	return nil
}

func ValidateSessionActivity(session domain.VisitorSession, animation domain.CandleAnimation) error {
	if err := SessionEntity(session); err != nil {
		return err
	}
	if session.LastAnimationID != "" && session.LastAnimationID != animation.ID {
		return fmt.Errorf("session points to %s instead of %s", session.LastAnimationID, animation.ID)
	}
	if session.InteractionCount < animation.Clicks {
		return fmt.Errorf("session count %d is behind clicks %d", session.InteractionCount, animation.Clicks)
	}
	return nil
}

func ValidateDisplayText(text string) error {
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("display text is empty")
	}
	if len([]rune(text)) > domain.MaxMessageLen+100 {
		return fmt.Errorf("display text is too long")
	}
	return nil
}

func RulesPassed(results []RuleResult) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return true
}
