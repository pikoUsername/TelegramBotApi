package objects

// Chat type
type Chat struct {
	ID        int32  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type ChatMember struct {
}

type MyChatMember struct {
}
