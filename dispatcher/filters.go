package dispatcher

import "github.com/pikoUsername/tgp/objects"

// Filter uses for filtering comming messsage etc.
// When Filter check is sucess handler will not pass
// If not just continue comming update, and etc. etc.
type Filter interface {
	Check(*objects.Update)
}
