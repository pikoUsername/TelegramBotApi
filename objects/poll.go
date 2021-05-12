package objects

import "time"

type Poll struct {
	ID                    int64            `json:"id"`
	Question              string           `json:"question"`
	Options               []*PollOption    `json:"options"`
	TotalVoterCount       int              `json:"total_voter_count"`
	IsClosed              bool             `json:"is_closed"`
	IsAnonymoust          bool             `json:"is_anonymous"`
	Type                  string           `json:"type"`
	AllowsMultipleAnswers bool             `json:"allows_multiple_answers"`
	CorrectOptionId       int64            `json:"correct_option_id"`
	Explanation           string           `json:"explanation"`
	ExplanationEntities   []*MessageEntity `json:"explanation_entities"`
	OpenPeriod            time.Duration    `json:"open_period"`
	CloseDate             time.Duration    `json:"close_date"`
}

type PollAnswer struct {
	PollID    int64 `json:"poll_id"`
	*User     `json:"user"`
	OptionIDs []int64 `json:"option_ids"`
}

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}
