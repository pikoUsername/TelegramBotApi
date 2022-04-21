package objects

type (
	LabeledPrice struct {
		Label  string `json:"label"`
		Amount int    `json:"amount"`
	}

	Invoice struct {
		TItle          string `json:"title"`
		Description    string `json:"description"`
		StartParameter string `json:"start_parameter"`
		Currency       string `json:"currency"`
		TotalAmount    int    `json:"total_amount"`
	}

	OrderInfo struct {
		Name            string          `json:"name"`
		PhoneNumber     string          `json:"phone_number"`
		Email           string          `json:"email"`
		ShippingAddress ShippingAddress `json:"shipping_address"`
	}

	ShippingOption struct {
		ID     string         `json:"id"`
		Title  string         `json:"title"`
		Prices []LabeledPrice `json:"prices"`
	}

	SuccessfulPayment struct {
		Currency                string    `json:"currency"`
		TotalAmount             int       `json:"total_amount"`
		InvoicePayload          string    `json:"invoice_payload"`
		ShippingOptionID        string    `json:"shipping_option_id"`
		OrderInfo               OrderInfo `json:"order_info"`
		TelegramPaymentChargeID string    `json:"telegram_payment_charge_id"`
		ProviderPaymentChargeID string    `json:"provider_payment_charge_id"`
	}

	PreCheckoutQuery struct {
		ID               string    `json:"id"`
		From             *User     `json:"from"`
		Currency         string    `json:"currency"`
		TotalAmount      int       `json:"total_amount"`
		InvoicePayload   string    `json:"invoice_payload"`
		ShippingOptionID string    `json:"shipping_option_id"`
		OrderInfo        OrderInfo `json:"order_info"`
	}
)
