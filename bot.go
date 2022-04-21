package tgp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	multipartreader "github.com/pikoUsername/MultipartReader"
	"github.com/pikoUsername/tgp/objects"
)

// If TG API method's arguments count less than 3,
// then it will be passed as arguments in the function.

// Bot can be created using Json config,
// Copy-pasted from go-telegram-bot-api
//
// Client field have timeout, and default timeout is 5 second
type Bot struct {
	// Token uses for authonificate using URL
	// Url template {api_url}/bot{bot_token}/{method}?{args}
	Token string `json:"token"`

	// i will recomend to use HTML parse_mode
	// bc, HTML easy to use, and more conforatble
	ParseMode string   `json:"parse_mode"`
	Markdown  Markdown `json:"-"`

	// default server must be here
	// if you wanna create own, just create
	// using this structure instead of NewBot function
	Server ITelegramServer `json:"-"`

	// Using prefix Bot, for avoid names conflict
	// and golang dont love name conflicts
	// by default this values is nil,
	// when you make get_me request, result
	// caches there, and you can take that
	// value in any moment.
	// Using Lazy method, instead of on startup init
	Me *objects.User `json:"me"`

	// Client uses for requests
	Client *http.Client `json:"-"`
}

// NewBot returns a new bot struct which need to interact with Telegram Bot API
// Bot structure should provide only Telegram bot API methods
func NewBot(token string, parseMode string, client *http.Client) (*Bot, error) {
	// Check out for correct token
	if strings.Contains(token, " ") {
		return nil, tgpErr.New("token is invalid! token contains space")
	}
	if client == nil {
		client = &http.Client{}
	}
	var mrkdown Markdown
	if strings.ToLower(parseMode) == "html" {
		mrkdown = HTMLDecoration
	} else {
		mrkdown = MarkdownDecoration
	}
	return &Bot{
		Token:     token,
		ParseMode: parseMode,
		Server:    DefaultTelegramServer,
		Markdown:  mrkdown,
		Client:    client,
	}, nil
}

func (bot *Bot) SetTimeout(dur time.Duration) {
	bot.Client.Timeout = dur
}

// ===================
// sending requests
// ===================

// Request to telegram servers
// and result parses to TelegramResponse
func (bot *Bot) Request(Method string, params url.Values) (*objects.TelegramResponse, error) {
	tgurl := bot.Server.ApiURL(bot.Token, Method)

	request, err := http.NewRequest("POST", tgurl, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := bot.Client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tgresp, err := responseDecode(resp.Body)
	if err != nil {
		return nil, err
	}
	return checkResult(tgresp)
}

// BoolRequest call a Request, and return bool
// in telegram api there are many methods that return the Boolean value
func (bot *Bot) BoolRequest(method string, params url.Values) (bool, error) {
	var ok bool
	resp, err := bot.Request(method, params)
	if err != nil {
		return false, err
	}
	json.Unmarshal(resp.Result, &ok)
	return ok, err
}

// DownloadFile uses for download file from any URL,
func (bot *Bot) DownloadFile(path string, w io.Writer) error {
	request, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	resp, err := bot.Client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(bs)
	if err != nil {
		return err
	}

	return nil
}

// Upload file uploads file to telegram server
func (b *Bot) UploadFile(method string, v map[string]string, data ...*objects.InputFile) (*objects.TelegramResponse, error) {
	var name string
	var err error

	ms := multipartreader.New()
	err = ms.WriteFields(v)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)

	for _, value := range data {
		defer value.Close()
		if value.URL == "" && value.Name == "" && value.File == nil && value.Path == "" {
			return nil, tgpErr.New("err while uploading inputfile, file is empty")
		}
		if value.URL != "" && value.Name != "" {
			values[value.Name] = value.URL

			err = ms.WriteFields(values)
			if err != nil {
				return nil, err
			}
		}
		if value.File != nil && value.Name != "" {
			ms.AddFormReader(value.Name, value.Path, int64(value.Length), value.File)
		}
		if value.Path != "" && value.Name != "" && value.File == nil {
			err = ms.WriteFile(value.Name, value.Path)
			if err != nil {
				return nil, err
			}
		}
	}
	// ms.AddValuesReader(f_v)
	tgurl := b.Server.ApiURL(b.Token, name)

	req, err := http.NewRequest("POST", tgurl+method, nil)
	if err != nil {
		return nil, err
	}
	ms.SetupRequest(req)
	// sending request
	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}

	// closing body
	defer resp.Body.Close()
	tgresp, err := responseDecode(resp.Body)
	if err != nil {
		return nil, err
	}

	// returns response instant
	return checkResult(tgresp)
}

// GetMe represents telegram "getMe" method
// https://core.telegram.org/bots/api#getme
func (bot *Bot) GetMe() (*objects.User, error) {
	if bot.Me != nil {
		return bot.Me, nil
	}
	resp, err := bot.Request("getMe", url.Values{})
	if err != nil {
		return nil, err
	}
	var user objects.User
	json.Unmarshal(resp.Result, &user)
	bot.Me = &user // caching result
	return &user, nil
}

// Logout your bot from telegram
// https://core.telegram.org/bots/api#logout
func (bot *Bot) Logout() (*objects.TelegramResponse, error) {
	return bot.Request("logout", url.Values{})
}

// ===============================
// No returning value Telegram api methods
// ===============================

// DeleteChatPhoto represents deleteChatPhoto method
// https://core.telegram.org/bots/api#deletechatphoto
func (bot *Bot) DeleteChatPhoto(ChatId int64) (*objects.TelegramResponse, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatId, 10))

	resp, err := bot.Request("deleteChatPhoto", v)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetchatTitle respresents setChatTitle method
// https://core.telegram.org/bots/api#setChatTitle
func (bot *Bot) SetChatTitle(ChatId int64, Title string) (*objects.TelegramResponse, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("title", Title)

	resp, err := bot.Request("setChatTitle", v)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// SetChatDescription respresents setChatDescription method
// https://core.telegram.org/bots/api#setChatDescription
func (bot *Bot) SetChatDescription(ChatId int64, Description string) (*objects.TelegramResponse, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("description", Description)
	resp, err := bot.Request("setChatDescription", v)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// PinChatMessage respresents pinChatMessage method
// https://core.telegram.org/bots/api#pinchatmessage
func (bot *Bot) PinChatMessage(
	ChatId int64,
	MessageId int64,
	DisableNotifiaction bool,
) (*objects.TelegramResponse, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	v.Add("message_id", strconv.FormatInt(MessageId, 10))
	v.Add("disable_notifications", strconv.FormatBool(DisableNotifiaction))
	resp, err := bot.Request("pinChatMessage", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UnpinAllChatMessage respresents unpinAllChatMessages method
// https://core.telegram.org/bots/api#unpinAllChatMessages
func (bot *Bot) UnpinAllChatMessages(ChatId int64) (*objects.TelegramResponse, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatId, 10))
	resp, err := bot.Request("unpinAllChatMessages", v)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// =============================
// Message sending
// =============================

// Send uses as sender for almost all stuff
func (bot *Bot) SendMessageable(c Configurable) (*objects.Message, error) {
	v, err := c.values()
	if err != nil {
		return nil, err
	}
	// Check out for parse_mode and set bot.ParseMode if config parse_mode is empty
	if v.Get("parse_mode") == "" {
		v.Set("parse_mode", bot.ParseMode)
	}
	resp, err := bot.Request(c.method(), v)

	if err != nil {
		return nil, err
	}
	var msg objects.Message
	json.Unmarshal(resp.Result, &msg)
	return &msg, nil
}

// uploadAndSend will send a Message with a new file to Telegram.
func (bot *Bot) UploadAndSend(config FileableConf) (*objects.Message, error) {
	params, err := config.params()
	if err != nil {
		return nil, err
	}

	method := config.method()
	resp, err := bot.UploadFile(method, params, config.getFiles()...)
	if err != nil {
		return nil, err
	}

	var message *objects.Message
	json.Unmarshal(resp.Result, &message)
	return message, nil
}

// Send ...
func (bot *Bot) Send(config Configurable) (*objects.Message, error) {
	switch c := config.(type) {
	case FileableConf:
		return bot.UploadAndSend(c)
	case Configurable:
		return bot.SendMessageable(c)
	}
	return nil, tgpErr.New("config is not correct")
}

// CopyMessage copies message
// https://core.telegram.org/bots/api#copymessage
func (bot *Bot) CopyMessage(config *CopyMessageConfig) (*objects.MessageID, error) {
	v, err := config.values()

	if err != nil {
		return &objects.MessageID{}, err
	}
	resp, err := bot.Request(config.method(), v)
	if err != nil {
		return &objects.MessageID{}, err
	}
	var msg objects.MessageID

	json.Unmarshal(resp.Result, &msg)
	return &msg, nil
}

// SendPhoto ...
func (bot *Bot) SendPhoto(config *SendPhotoConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendAudio ...
func (bot *Bot) SendAudio(config *SendAudioConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendDocument ...
func (bot *Bot) SendDocument(config *SendDocumentConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVideo ...
func (bot *Bot) SendVideo(config *SendVideoConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendAnimation ...
func (bot *Bot) SendAnimation(config *SendAnimationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVoice ...
func (bot *Bot) SendVoice(config *SendVoiceConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendVideoName ...
func (bot *Bot) SendVideoNote(config *SendVideoNoteConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendMediaGroup ...
func (bot *Bot) SendMediaGroup(config *SendMediaGroupConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendLocation ...
func (bot *Bot) SendLocation(config *SendLocationConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// editMessageLiveLocation ...
func (bot *Bot) EditMessageLiveLocation(config *EditMessageLLConf) (*objects.Message, error) {
	return bot.Send(config)
}

// SendMessage sends message using ChatID
// see: https://core.telegram.org/bots/api#sendmessage
func (bot *Bot) SendMessage(config *SendMessageConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// SendPoll Use this method to send a native poll
// https://core.telegram.org/bots/api#sendpoll
func (bot *Bot) SendPoll(config *SendPollConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendDice(config *SendDiceConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendContact(config *SendContactConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendVenue(config *SendVenueConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendGame(config *SendGameConfig) (*objects.Message, error) {
	return bot.Send(config)
}

func (bot *Bot) SendSticker(config *SendStickerConfig) (*objects.Message, error) {
	return bot.Send(config)
}

// =========================
// Commands Methods
// =========================

// SetMyCommands Setup command to Telegram bot
// https://core.telegram.org/bots/api#setmycommands
func (bot *Bot) SetMyCommands(conf *SetMyCommandsConfig) (bool, error) {
	v, _ := conf.values()
	return bot.BoolRequest(conf.method(), v)
}

// GetMyCommands get from bot commands command
// https://core.telegram.org/bots/api#getmycommands
func (bot *Bot) GetMyCommands(c *GetMyCommandsConfig) ([]objects.BotCommand, error) {
	v, _ := c.values()
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return []objects.BotCommand{}, err
	}
	var cmds []objects.BotCommand
	json.Unmarshal(resp.Result, &cmds)
	return cmds, nil
}

// ======================
// Getting Updates
// ======================

// DeleteWebhook if result is True, will be nil, if not so err
// https://core.telegram.org/bots/api#deletewebhook
func (bot *Bot) DeleteWebhook(c *DeleteWebhookConfig) (*objects.TelegramResponse, error) {
	v, err := c.values()
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	return resp, nil
}

// GetUpdates uses for long polling
// https://core.telegram.org/bots/api#getupdates
// Telegram will hold updates only on 24 hours
func (bot *Bot) GetUpdates(c *GetUpdatesConfig) ([]*objects.Update, error) {
	var updates []*objects.Update
	v, err := c.values()
	if err != nil {
		return updates, err
	}
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return updates, err
	}
	var upd []*objects.Update
	json.Unmarshal(resp.Result, &upd)
	return upd, nil
}

// SetWebhook make subscribe to telegram events
// or sends to telegram a request for make
// Subscribe to specific IP, and when user
// sends a message to your bot, Telegram know
// Your bot IP and sends to your bot a Update
// https://core.telegram.org/bots/api#setwebhook
func (bot *Bot) SetWebhook(c *SetWebhookConfig) (*objects.TelegramResponse, error) {
	v, err := c.values()
	meth := c.method()

	// checkout for certificate, webhook may use without cert
	if c.Certificate == nil { // you don't have to send your certificate to telegram
		return bot.Request(meth, v)
	}
	if err != nil {
		return &objects.TelegramResponse{}, err
	}
	params := make(map[string]string)
	urlValuesToMapString(v, params)

	// uploads a certificate file, with other parametrs
	resp, err := bot.UploadFile(meth, params, c.Certificate)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

// GetWebhookInfo not require parametrs
// https://core.telegram.org/bots/api#getwebhookinfo
func (bot *Bot) GetWebhookInfo() (*objects.WebhookInfo, error) {
	resp, err := bot.Request("getWebhookInfo", url.Values{})
	if err != nil {
		return &objects.WebhookInfo{}, err
	}
	var wi objects.WebhookInfo
	json.Unmarshal(resp.Result, &wi)
	return &wi, nil
}

// =====================
// Chat methods
// =====================

// SendChatAction Resrpesents sendChatAction method
// https://core.telegram.org/bots/api#sendchataction
// Use this method when you need to tell the user that something is happening on the bot's side.
// The status is set for 5 seconds or less
// (when a message arrives from your bot, Telegram clients clear its typing status).
// Returns True on success.
func (bot *Bot) SendChatAction(c SendChatActionConf) (bool, error) {
	v, err := c.values()
	if err != nil {
		return false, err
	}
	return bot.BoolRequest(c.method(), v)
}

// DeleteChatStickerSet represents deleteChatStickerSet method
// https://core.telegram.org/bots/api#deletechatstickerset
func (bot *Bot) DeleteChatStickerSet(chat_id int64) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	return bot.BoolRequest("deleteChatStickerSet", v)
}

// GetChat ...
func (bot *Bot) GetChat(chat_id int64) (*objects.Chat, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))

	resp, err := bot.Request("getChat", v)

	if err != nil {
		return &objects.Chat{}, err
	}

	var chat objects.Chat
	json.Unmarshal(resp.Result, &chat)
	return &chat, nil
}

// BanChatMember ...
func (bot *Bot) BanChatMember(c *BanChatMemberConfig) (bool, error) {
	v, _ := c.values()
	return bot.BoolRequest(c.method(), v)
}

func (bot *Bot) GetChatMemberCount(chat_id int64) (int, error) {
	var count int
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	resp, err := bot.Request("getChatMemberCount", v)
	if err != nil {
		return count, err
	}
	err = json.Unmarshal(resp.Result, &count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// UnbanChatMember ...
func (bot *Bot) UnbanChatMember(chat_id int64, user_id int64, only_if_banned bool) (bool, error) {
	v := make(url.Values)
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	v.Add("user_id", strconv.FormatInt(user_id, 10))
	v.Add("only_if_banned", strconv.FormatBool(only_if_banned))
	return bot.BoolRequest("unbanChatMember", v)
}

// RestrictChatMember ...
func (bot *Bot) RestrictChatMember(c *RestrictChatMemberConfig) (bool, error) {
	v, _ := c.values()
	return bot.BoolRequest(c.method(), v)
}

// SetChatPermissions ...
func (bot *Bot) SetChatPermissions(chat_id int64, perms objects.ChatMemberPermissions) (bool, error) {
	v := make(url.Values)

	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	v.Add("permissions", ObjectToJson(perms))

	return bot.BoolRequest("setChatPermissions", v)
}

// SCACT too long to read, represents a telegram method
// https://core.telegram.org/bots/api#setchatadministratorcustomtitle
func (bot *Bot) SetChatAdministratorCustomTitle(chat_id, user_id int64, title string) (bool, error) {
	v := make(url.Values)

	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	v.Add("user_id", strconv.FormatInt(user_id, 10))
	v.Add("title", title)

	return bot.BoolRequest("setChatAdministratorCustomTitle", v)
}

// ExportchatInviteLink ...
func (bot *Bot) ExportChatInviteLink(chat_id int64) (string, error) {
	v := make(url.Values)

	v.Add("chat_id", strconv.FormatInt(chat_id, 10))

	resp, err := bot.Request("exportChatInviteLink", v)
	if err != nil {
		return "", err
	}

	var val string
	json.Unmarshal(resp.Result, &val)
	return val, err
}

// EditInviteLink ...
func (bot *Bot) EditInviteLink(c *EditChatInviteLinkConf) (cil *objects.ChatInviteLink, err error) {
	v, _ := c.values()
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp.Result, cil)
	return
}

func (bot *Bot) SetChatPhoto(chat_id int64, file *objects.InputFile) (bool, error) {
	v := make(map[string]string)

	v["chat_id"] = strconv.FormatInt(chat_id, 10)

	bot.UploadFile("setChatPhoto", v, file)

	return false, nil
}

func (bot *Bot) RevokeChatInviteLink(chat_id int64, invoke string) (val *objects.ChatInviteLink, err error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	v.Add("invoke", invoke)

	resp, err := bot.Request("revokeChatInviteLink", v)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(resp.Result, val)
	return
}

func (bot *Bot) ApproveChatJoinRequest(chat_id int64, user_id int64) (bool, error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(chat_id, 10))
	v.Add("user_id", strconv.FormatInt(user_id, 10))
	return bot.BoolRequest("approveChatJoinRequest", v)
}

func (bot *Bot) SetMyDefaultAdministratorRights(rights *objects.ChatAdministratorRights, for_channels bool) (bool, error) {
	v := url.Values{}

	if rights != nil {
		bs, err := json.Marshal(rights)
		if err != nil {
			return false, err
		}
		v.Add("rights", BytesToString(bs))
	}
	v.Add("for_channels", strconv.FormatBool(for_channels))
	return bot.BoolRequest("setMyDefaultAdministratorRights", v)
}

func (bot *Bot) GetMyDefaultAdministratorRights(for_channels bool) (*objects.ChatAdministratorRights, error) {
	v := url.Values{}

	v.Add("for_channels", strconv.FormatBool(for_channels))
	resp, err := bot.Request("getMyDefaultAdministratorRights", v)

	if err != nil {
		return nil, err
	}

	var rights *objects.ChatAdministratorRights
	err = json.Unmarshal(resp.Result, rights)
	if err != nil {
		return nil, err
	}

	return rights, nil
}

// ================
// User methods
// ================

// GetUserProfilePhotos resresents getUserProfilePhotos method
// https://core.telegram.org/bots/api#getuserprofilephotos
func (bot *Bot) GetUserProfilePhotos(c GetUserProfilePhotosConf) (*objects.UserProfilePhotos, error) {
	v, _ := c.values()
	resp, err := bot.Request(c.method(), v)

	if err != nil {
		return nil, err
	}

	var photos objects.UserProfilePhotos
	err = json.Unmarshal(resp.Result, &photos)
	if err != nil {
		return nil, err
	}

	return &photos, nil
}

// ====================
// Sticker methods
// ====================

func (bot *Bot) DeleteStickerFromSet(sticker string) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)
	return bot.BoolRequest("setStickerPositionInSet", v)
}

func (bot *Bot) SetStickerPositionInSet(sticker, position string) (bool, error) {
	v := url.Values{}
	v.Add("sticker", sticker)
	v.Add("position", position)
	return bot.BoolRequest("setStickerPositionInSet", v)
}

func (bot *Bot) GetStickerSet(name string) (ss *objects.StickerSet, err error) {
	v := url.Values{}
	v.Add("name", name)
	resp, err := bot.Request("getStickerSet", v)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp.Result, ss)
	return ss, err
}

func (bot *Bot) UploadStickerFile(user_id int64, png_sticker *objects.InputFile) (*objects.File, error) {
	v := make(map[string]string)
	v["user_id"] = strconv.FormatInt(user_id, 10)
	resp, err := bot.UploadFile("uploadStickerFile	", v, png_sticker)
	if err != nil {
		return nil, err
	}

	var file *objects.File
	err = json.Unmarshal(resp.Result, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (bot *Bot) SetStickerSetThumb(c *SetStickerSetThumbConf) (bool, error) {
	v, _ := c.params()
	resp, err := bot.UploadFile("uploadStickerFile	", v, c.Thumb)
	if err != nil {
		return false, err
	}

	var ok bool
	err = json.Unmarshal(resp.Result, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (bot *Bot) CreateNewStickerSet(c *CreateNewStickerSetConf) (bool, error) {
	v, _ := c.values()
	return bot.BoolRequest(c.method(), v)
}

func (bot *Bot) AddStickerToSet(c *AddStickerToSetConf) (bool, error) {
	v, _ := c.values()
	return bot.BoolRequest(c.method(), v)
}

// ====================
// other methods
// ====================

func (bot *Bot) GetFile(file_id string) (*objects.File, error) {
	v := url.Values{}
	v.Add("file_id", file_id)
	resp, err := bot.Request("getFile", v)

	if err != nil {
		return nil, err
	}

	var file *objects.File
	json.Unmarshal(resp.Result, &file)
	return file, nil
}

func (bot *Bot) PromoteChatMember(config PromoteChatMemberConfig) (bool, error) {
	var ok bool
	v, err := config.values()
	if err != nil {
		return false, err
	}
	resp, err := bot.Request(config.method(), v)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(resp.Result, &ok)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// =====================
// Web
// =====================
func (bot *Bot) AnswerWebAppQuery(c AnswerWebAppQueryConf) (*objects.SentWebAppMessage, error) {
	v, err := c.values()
	if err != nil {
		return nil, err
	}
	resp, err := bot.Request(c.method(), v)
	if err != nil {
		return nil, err
	}
	var result *objects.SentWebAppMessage
	err = json.Unmarshal(resp.Result, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (bot *Bot) SetChatMenuButton(ChatID int64, MenuButton *objects.MenuButton) (ok bool, err error) {
	v := url.Values{}

	v.Add("chat_id", strconv.FormatInt(ChatID, 10))
	bs, err := json.Marshal(MenuButton)

	if err != nil {
		return false, err
	}
	v.Add("menu_button", BytesToString(bs))

	return bot.BoolRequest("setChatMenuButton", v)
}

func (bot *Bot) GetChatMenuButton(ChatID int64) (menu *objects.MenuButton, err error) {
	v := url.Values{}
	v.Add("chat_id", strconv.FormatInt(ChatID, 10))
	resp, err := bot.Request("getChatMenuButton", v)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Result, menu)
	if err != nil {
		return nil, err
	}
	return menu, nil
}
