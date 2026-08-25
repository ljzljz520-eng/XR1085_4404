package store

import (
	"memorialcandle/domain"
	"memorialcandle/validate"
)

func ValidateBundle(bundle Bundle) error {
	if err := validate.MessageEntity(bundle.Message); err != nil {
		return err
	}
	if err := validate.AnimationEntity(bundle.Animation); err != nil {
		return err
	}
	if err := validate.SnapshotEntity(bundle.Snapshot); err != nil {
		return err
	}
	if err := validate.SessionEntity(bundle.Session); err != nil {
		return err
	}
	if err := validate.Linked(bundle.Message, bundle.Animation, bundle.Snapshot); err != nil {
		return err
	}
	if bundle.Session.LastAnimationID != "" && bundle.Session.LastAnimationID != bundle.Animation.ID {
		return nil
	}
	if bundle.Animation.Stage == domain.StageQuiet {
		return validate.QuietState(bundle.Animation, bundle.Snapshot)
	}
	return nil
}
