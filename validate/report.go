package validate

import (
	"fmt"
	"memorialcandle/domain"
)

type Report struct {
	Checked  int
	Passed   int
	Failures []string
}

func NewReport() Report { return Report{Failures: make([]string, 0)} }
func (r *Report) add(name string, err error) {
	r.Checked++
	if err == nil {
		r.Passed++
	} else {
		r.Failures = append(r.Failures, name+": "+err.Error())
	}
}

func ValidateEntities(message domain.MemorialMessage, session domain.VisitorSession, animation domain.CandleAnimation, snapshot domain.DisplaySnapshot) Report {
	report := NewReport()
	report.add("message", MessageEntity(message))
	report.add("session", SessionEntity(session))
	report.add("animation", AnimationEntity(animation))
	report.add("snapshot", SnapshotEntity(snapshot))
	report.add("links", Linked(message, animation, snapshot))
	return report
}

func (r Report) OK() bool { return r.Checked == r.Passed && len(r.Failures) == 0 }
func (r Report) Summary() string {
	return fmt.Sprintf("checked=%d passed=%d failures=%d", r.Checked, r.Passed, len(r.Failures))
}

func EnsureQuietBundle(message domain.MemorialMessage, animation domain.CandleAnimation, snapshot domain.DisplaySnapshot) error {
	if err := Linked(message, animation, snapshot); err != nil {
		return err
	}
	return QuietState(animation, snapshot)
}

func CounterBounds(clicks, generation, layers int) error {
	if clicks < 0 || clicks > domain.MaxClicks {
		return fmt.Errorf("clicks out of bounds")
	}
	if generation < 1 {
		return fmt.Errorf("generation out of bounds")
	}
	if layers < 1 {
		return fmt.Errorf("layers out of bounds")
	}
	return nil
}
