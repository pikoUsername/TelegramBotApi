package filters

import "github.com/pikoUsername/tgp/objects"

// Unreachable filter, always return false
func Unreachable(upd *objects.Update) bool {
	return false
}
