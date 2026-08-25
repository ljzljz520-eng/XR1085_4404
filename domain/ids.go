package domain

import "fmt"

func MessageID(sequence int) string        { return fmt.Sprintf("memorial-%04d", sequence) }
func SessionID(sequence int) string        { return fmt.Sprintf("visitor-%04d", sequence) }
func AnimationID(messageID string) string  { return messageID + "-candle" }
func SnapshotID(animationID string) string { return animationID + "-snapshot" }

func NormalizeSequence(sequence int) int {
	if sequence < 1 {
		return 1
	}
	if sequence > 9999 {
		return 9999
	}
	return sequence
}
