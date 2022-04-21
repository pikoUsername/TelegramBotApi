package objects

type InputMessageContent interface {
	getContent() // stub
}

type InputTextMessageContent struct {
	MessageText           string          `json:"message_text"`
	ParseMode             string          `json:"parse_mode"`
	Entities              []MessageEntity `json:"entities"`
	DisableWebPagePreview bool            `json:"disable_web_page_preview"`
}

type InputLocationMessageContent struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proxomity_alert_radius"`
}

type InputVenueMessageContent struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Title           string  `json:"title"`
	Address         string  `json:"address"`
	FourSquareId    string  `json:"four_square_id"`
	FourSquareType  string  `json:"four_square_type"`
	GooglePlaceID   string  `json:"google_place_id"`
	GooglePlaceType string  `json:"google_place_type"`
}

type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Vcard       string `json:"vcard"`
}

type InputInvoiceMessageContent struct {
	Title                     string         `json:"title"`
	Description               string         `json:"description"`
	Payload                   string         `json:"payload"`
	ProviderToken             string         `json:"provider_token"`
	Prices                    []LabeledPrice `json:"prices"`
	MaxTipAmount              int            `json:"max_tip_amount"`
	SuggestedTipAmounts       []int          `json:"suggested_tip_amounts"`
	ProviderData              string         `json:"provider_data"`
	PhotoSize                 int            `json:"photo_size"`
	PhotoWidth                int            `json:"photo_width"`
	PhotoHeight               int            `json:"photo_height"`
	NeedName                  bool           `json:"need_name"`
	NeedPhoneNumber           bool           `json:"need_phone_number"`
	NeedEmail                 bool           `json:"need_email"`
	NeedShippingAddress       bool           `json:"need_shipping_address"`
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider"`
	SendEmailToProvider       bool           `json:"send_email_to_provider"`
	IsFlexible                bool           `json:"is_flexible"`
}

func (InputInvoiceMessageContent) getContent()  {}
func (InputContactMessageContent) getContent()  {}
func (InputVenueMessageContent) getContent()    {}
func (InputLocationMessageContent) getContent() {}
func (InputTextMessageContent) getContent()     {}
