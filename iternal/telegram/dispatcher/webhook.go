package dispatcher

import "net/url"

// WebhookConfig uses for Using as arguemnt
// You may not fill all fields in struct
type WebhookConfig struct {
	URL            *url.URL
	Certificate    interface{}
	Offset         int
	MaxConnections int
}
