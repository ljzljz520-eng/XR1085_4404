package validate

import (
	"errors"
	"memorialcandle/domain"
)

func MessageEntity(message domain.MemorialMessage) error     { return message.Validate() }
func SessionEntity(session domain.VisitorSession) error      { return session.Validate() }
func AnimationEntity(animation domain.CandleAnimation) error { return animation.Validate() }
func SnapshotEntity(snapshot domain.DisplaySnapshot) error   { return snapshot.Validate() }

func Linked(message domain.MemorialMessage, animation domain.CandleAnimation, snapshot domain.DisplaySnapshot) error {
	if message.ID == "" || animation.MessageID != message.ID {
		return errors.New("message and animation are not linked")
	}
	if snapshot.AnimationID != animation.ID {
		return errors.New("animation and snapshot are not linked")
	}
	if snapshot.Stage != animation.Stage || snapshot.Generation != animation.Generation {
		return errors.New("snapshot is stale")
	}
	return nil
}

func QuietState(animation domain.CandleAnimation, snapshot domain.DisplaySnapshot) error {
	if animation.Stage != domain.StageQuiet || !animation.Quiet || animation.Layers != 1 {
		return errors.New("animation is not single quiet state")
	}
	if snapshot.Stage != domain.StageQuiet || !snapshot.Quiet || snapshot.ParticleCount != 1 {
		return errors.New("snapshot is not single quiet state")
	}
	return nil
}
