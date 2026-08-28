package cli

import (
	"fmt"
	"io"
	"memorialcandle/store"
)

func WriteBundle(out io.Writer, bundle store.Bundle) error {
	_, err := fmt.Fprintf(out, "%s\nstage=%s particles=%d quiet=%t generation=%d clicks=%d layers=%d\n", bundle.Snapshot.VisibleText, bundle.Snapshot.Stage, bundle.Snapshot.ParticleCount, bundle.Snapshot.Quiet, bundle.Snapshot.Generation, bundle.Animation.Clicks, bundle.Animation.Layers)
	return err
}
