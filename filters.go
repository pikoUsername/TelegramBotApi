package tgp

import "github.com/pikoUsername/tgp/objects"

// Filter uses for filtering comming messsage etc.
// When Filter check is sucess handler will not pass
// If not just continue comming update, and etc. etc.
// if you need aleardy ready filters check filters/
// (looks good)
type Filter interface {
	Check(update *objects.Update) bool
}
