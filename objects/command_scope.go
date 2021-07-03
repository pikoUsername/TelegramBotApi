package objects

// BotCommadnsScope is generic interface{}
// GetType is stub method, for filter a other interfaces
// And reporesents a Telegram scope object
// https://core.telegram.org/bots/api#botcommandscope
type BotCommandScope interface {
	GetType() string
}

// BotCommandScopeDefault ...
type BotCommandScopeDefault struct {
	Type_ string `json:"type"`
}

func (bccd *BotCommandScopeDefault) GetType() string {
	return bccd.Type_
}

// BotCommandScopeAllPrivateChats ...
type BotCommandScopeAllPrivateChats struct {
	Type_ string `json:"type"`
}

func (bcsapc *BotCommandScopeAllPrivateChats) GetType() string {
	return bcsapc.Type_
}

// BotCommandScopeAllGroupChats ...
type BotCommandScopeAllGroupChats struct {
	Type_ string `json:"type"`
}

func (bcsagc *BotCommandScopeAllGroupChats) GetType() string {
	return bcsagc.Type_
}

// BotCommandScopeChat ...
type BotCommandScopeChat struct {
	Type_  string `json:"type"`
	ChatID int64  `json:"chat_id"`
}

func (bcsc *BotCommandScopeChat) GetType() string {
	return bcsc.Type_
}

// BotCommandScopeChatAdministrators ...
type BotCommandScopeChatAdministrators struct {
	Type_  string `json:"type"`
	ChatID int64  `json:"chat_id"`
}

func (bcsc *BotCommandScopeChatAdministrators) GetType() string {
	return bcsc.Type_
}

// BotCommandScopeChatMember ...
type BotCommandScopeChatMember struct {
	Type_  string `json:"type"`
	ChatID int64  `json:"chat_id"`
	UserID int64  `json:"user_id"`
}

func (bcscm *BotCommandScopeChatMember) GetType() string {
	return bcscm.Type_
}
