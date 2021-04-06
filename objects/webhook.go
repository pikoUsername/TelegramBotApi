package objects

// webhookInfo structure from getWebhookInfo
// https://core.telegram.org/bots/api#webhookinfo
type WebhookInfo struct {
	URL                  string   `json:"url"`
	HasCustomCertificate bool     `json:"has_custom_certificate"`
	PendingUpdateCount   int      `json:"pending_update_count"`
	IpAddress            string   `json:"ip_address"`
	LastErrorDate        int      `json:"last_error_date"`
	MaxConnections       int      `json:"max_connections"`
	AllowedUpdates       []string `json:"allowed_updates"`
}
