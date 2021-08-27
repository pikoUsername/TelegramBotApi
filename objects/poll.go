package objects

type Poll struct {
	ExplanationEntities   []*MessageEntity `json:"explanation_entities"`
	Options               []*PollOption    `json:"options"`
	AllowsMultipleAnswers bool             `json:"allows_multiple_answers"`
	IsClosed              bool             `json:"is_closed"`
	IsAnonymoust          bool             `json:"is_anonymous"`
	Question              string           `json:"question"`
	Type                  string           `json:"type"`
	Explanation           string           `json:"explanation"`
	TotalVoterCount       int              `json:"total_voter_count"`
	CorrectOptionId       int64            `json:"correct_option_id"`
	ID                    int64            `json:"id"`
	OpenPeriod            int64            `json:"open_period"`
	CloseDate             int64            `json:"close_date"`
}

type PollAnswer struct {
	User      *User   `json:"user"`
	PollID    int64   `json:"poll_id"`
	OptionIDs []int64 `json:"option_ids"`
}

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}
