package objects

// From https://core.telegram.org/bots/api#inlinequery:
// id	String	Unique identifier for this query
// from	User	Sender
// location	Location	Optional. Sender location, only for bots that request user location
// query	String	Text of the query (up to 256 characters)
// offset	String	Offset of the results to be returned, can be controlled by the bot
type InlineQuery struct {
	Id       string `json:"id"`
	From     *User  `json:"from"`
	Location `json:"Location"`
	Query    string `json:"query"`
	Offset   string `json:"offset"`
}
