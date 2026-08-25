package domain

import (
	"sort"
	"strings"
)

type MemorialCollection struct{ messages map[string]MemorialMessage }

func NewMemorialCollection() *MemorialCollection {
	return &MemorialCollection{messages: make(map[string]MemorialMessage)}
}

func (c *MemorialCollection) Add(message MemorialMessage) error {
	if err := message.Validate(); err != nil {
		return err
	}
	if c.messages == nil {
		c.messages = make(map[string]MemorialMessage)
	}
	c.messages[message.ID] = message
	return nil
}

func (c *MemorialCollection) Find(id string) (MemorialMessage, bool) {
	if c == nil {
		return MemorialMessage{}, false
	}
	message, ok := c.messages[id]
	return message, ok
}

func (c *MemorialCollection) Active() []MemorialMessage {
	result := make([]MemorialMessage, 0)
	if c == nil {
		return result
	}
	for _, message := range c.messages {
		if message.Active {
			result = append(result, message)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (c *MemorialCollection) Search(term string) []MemorialMessage {
	term = strings.ToLower(strings.TrimSpace(term))
	result := make([]MemorialMessage, 0)
	if c == nil || term == "" {
		return result
	}
	for _, message := range c.messages {
		if strings.Contains(strings.ToLower(message.Author), term) || strings.Contains(strings.ToLower(message.Text), term) {
			result = append(result, message)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result
}

func (c *MemorialCollection) Deactivate(id string) bool {
	if c == nil {
		return false
	}
	message, ok := c.messages[id]
	if !ok {
		return false
	}
	message.Active = false
	c.messages[id] = message
	return true
}

func (c *MemorialCollection) Count() int {
	if c == nil {
		return 0
	}
	return len(c.messages)
}
func (c *MemorialCollection) ActiveCount() int { return len(c.Active()) }

func (c *MemorialCollection) IDs() []string {
	ids := make([]string, 0)
	if c == nil {
		return ids
	}
	for id := range c.messages {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func (c *MemorialCollection) Merge(messages []MemorialMessage) error {
	for _, message := range messages {
		if err := c.Add(message); err != nil {
			return err
		}
	}
	return nil
}

func (c *MemorialCollection) Clone() *MemorialCollection {
	copyCollection := NewMemorialCollection()
	if c != nil {
		for id, message := range c.messages {
			copyCollection.messages[id] = message
		}
	}
	return copyCollection
}
