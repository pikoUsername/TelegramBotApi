package objects

type ShippingQuery struct {
	ID               string `json:"id"`
	From             *User  `json:"from"`
	InvoicePayload   string `json:"invoice_payload"`
	*ShippingAddress `json:"shipping_address"`
}
