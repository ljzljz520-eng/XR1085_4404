package render

import "fmt"

func formatParticles(count int) string { return fmt.Sprintf("particles=%d", count) }
func formatQuiet(quiet bool) string {
	if quiet {
		return "quiet=true"
	}
	return "quiet=false"
}
func Header() string { return "memorial candlelight" }
func Help() string   { return "record, show, click, advance" }
