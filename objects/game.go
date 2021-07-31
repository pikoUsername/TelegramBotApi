package objects

// Game respresents a Game telegram object
// https://core.telegram.org/bots/api#game
type Game struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Photo        []PhotoSize      `json:"photo_size"`
	Text         string           `json:"text"`
	TextEntities []*MessageEntity `json:"text_entities"`
	Animation    *Animation       `json:"animation"`
}
