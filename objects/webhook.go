package objects

import "time"

// webhookInfo structure from getWebhookInfo
// https://core.telegram.org/bots/api#webhookinfo
type WebhookInfo struct {
	URL                  string        `json:"url"`
	HasCustomCertificate bool          `json:"has_custom_certificate"`
	PendingUpdateCount   int           `json:"pending_update_count"`
	IpAddress            string        `json:"ip_address"`
	LastErrorDate        time.Duration `json:"last_error_date"`
	LastErrorMessage     string        `json:"last_error_message"`
	MaxConnections       int           `json:"max_connections"`
	AllowedUpdates       []string      `json:"allowed_updates"`
}
