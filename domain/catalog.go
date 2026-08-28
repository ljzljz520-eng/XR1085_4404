package domain

import (
	"sort"
	"strings"
)

type CatalogEntry struct {
	Message   MemorialMessage
	Animation CandleAnimation
	Snapshot  DisplaySnapshot
}

func NewCatalogEntry(message MemorialMessage, animation CandleAnimation, snapshot DisplaySnapshot) CatalogEntry {
	return CatalogEntry{Message: message, Animation: animation, Snapshot: snapshot}
}
func (entry CatalogEntry) ID() string { return entry.Message.ID }
func (entry CatalogEntry) SearchText() string {
	return strings.ToLower(entry.Message.Author + " " + entry.Message.Text)
}
func (entry CatalogEntry) Quiet() bool { return entry.Animation.Quiet && entry.Snapshot.Quiet }
func (entry CatalogEntry) Valid() bool {
	return entry.Message.Validate() == nil && entry.Animation.Validate() == nil && entry.Snapshot.Validate() == nil && entry.Animation.MessageID == entry.Message.ID && entry.Snapshot.AnimationID == entry.Animation.ID
}

type Catalog struct{ entries map[string]CatalogEntry }

func NewCatalog() *Catalog { return &Catalog{entries: make(map[string]CatalogEntry)} }
func (c *Catalog) Put(entry CatalogEntry) bool {
	if !entry.Valid() {
		return false
	}
	if c.entries == nil {
		c.entries = make(map[string]CatalogEntry)
	}
	c.entries[entry.ID()] = entry
	return true
}
func (c *Catalog) Get(id string) (CatalogEntry, bool) {
	if c == nil {
		return CatalogEntry{}, false
	}
	entry, ok := c.entries[id]
	return entry, ok
}
func (c *Catalog) Remove(id string) bool {
	if c == nil {
		return false
	}
	if _, ok := c.entries[id]; !ok {
		return false
	}
	delete(c.entries, id)
	return true
}
func (c *Catalog) IDs() []string {
	ids := make([]string, 0)
	if c == nil {
		return ids
	}
	for id := range c.entries {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
func (c *Catalog) Search(term string) []CatalogEntry {
	term = strings.ToLower(strings.TrimSpace(term))
	result := make([]CatalogEntry, 0)
	if c == nil || term == "" {
		return result
	}
	for _, entry := range c.entries {
		if strings.Contains(entry.SearchText(), term) {
			result = append(result, entry)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID() < result[j].ID() })
	return result
}
func (c *Catalog) QuietEntries() []CatalogEntry {
	result := make([]CatalogEntry, 0)
	if c == nil {
		return result
	}
	for _, entry := range c.entries {
		if entry.Quiet() {
			result = append(result, entry)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID() < result[j].ID() })
	return result
}
func (c *Catalog) Count() int {
	if c == nil {
		return 0
	}
	return len(c.entries)
}
func (c *Catalog) Merge(entries []CatalogEntry) int {
	accepted := 0
	for _, entry := range entries {
		if c.Put(entry) {
			accepted++
		}
	}
	return accepted
}
