package objects

// from telegram docs https://core.telegram.org/bots/api#location:
// longitude	Float	Longitude as defined by sender
// latitude	Float	Latitude as defined by sender
// horizontal_accuracy	Float number	Optional. The radius of uncertainty for the location, measured in meters; 0-1500
// live_period	Integer	Optional. Time relative to the message sending date, during which the location can be updated, in seconds. For active live locations only.
// heading	Integer	Optional. The direction in which user is moving, in degrees; 1-360. For active live locations only.
// proximity_alert_radius	Integer	Optional. Maximum distance for proximity alerts about approaching another chat member, in meters. For sent live locations only.
type Location struct {
	Longitude            float32 `json:"longitude"`
	Latitude             float32 `json:"latitude"`
	HorizontalAccuracy   float32 `json:"horizontal_accuracy"`
	LifePeriod           int     `json:"live_period"`
	Heading              int     `json:"heading"`
	ProximityAlertRadius int     `json:"proximity_alert_radius"`
}

// ChatLocation same as Location, But
// With string addres, it s important too
// Ofc, if you are in FBI, FSB, or OMON(riot police)
type ChatLocation struct {
	// По айпи вычеслю
	Location *Location `json:"location"`
	Address  string    `json:"address"`
}
