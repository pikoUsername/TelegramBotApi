package tgp_test

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pikoUsername/tgp"
	"github.com/pikoUsername/tgp/objects"
)

var (
	ParseMode     = "HTML"
	TestChatID, _ = strconv.ParseInt(os.Getenv("test_chat_id"), 10, 64)
	FileDirectory = "./.sandbox"
	SaveFile      = path.Join(FileDirectory, "file")
	WebhookURL    = ""
	TestToken     = os.Getenv("TEST_TOKEN")

	// here could be any image, file, anthing else
	DownloadFromURL = "https://www.google.com/url?sa=i&url=https%3A%2F%2Fwww.dreamstime.com%2Fphotos-images%2Fimag.html&psig=AOvVaw1T5_yBwBBJzGYLRBvYTgA3&ust=1625481461095000&source=images&cd=vfe&ved=0CAoQjRxqFwoTCPiejL6cyfECFQAAAAAdAAAAABAJ"
	NothingInbytes  = []byte{}
	Timeout         = 2 * time.Second
)

func getBot(t *testing.T) *tgp.Bot {
	b, err := tgp.NewBot(TestToken, ParseMode, Timeout)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	b.Debug = true
	return b
}

func TestCheckToken(t *testing.T) {
	b, err := tgp.NewBot("bla:bla", "HTML", Timeout)
	if err != nil && b == nil {
		t.Error(err)
		t.Fail()
	}
}

func TestDownloadFile(t *testing.T) {
	b := getBot(t)
	dir, err := os.Open(FileDirectory)
	if err == os.ErrNotExist {
		os.Mkdir(FileDirectory, 0777)
		dir, err = os.Open(FileDirectory)
		if err != nil {
			t.Error(err)
		}
	}
	stat, err := dir.Stat()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if !stat.IsDir() {
		t.Error("Sorry but -"+FileDirectory, "Is file, delete file and try again!")
		t.Fail()

	}

	f, err := os.Create(SaveFile)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	err = b.DownloadFile(DownloadFromURL, f, true)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	stat, err = f.Stat()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	bs := make([]byte, stat.Size())
	f.Read(bs)
	if bs == nil || strings.Compare(string(bs), "") == -1 && DownloadFromURL != "" {
		t.Error(
			"Cannot download file from ethernet, debug: file - ", bs,
			", URL: ", DownloadFromURL, ", Directory: ", stat.Name(), ", Bot: ", b)
		t.Fail()
	}
}

func TestGetMe(t *testing.T) {
	b := getBot(t)
	u, err := b.GetMe()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if b.Me == nil {
		t.Error("Me is empty pointer")
		t.Fail()
	}
	if b.Me.ID != u.ID {
		t.Error("Getted User is defferent from bot user")
		t.Fail()
	}
}

func TestGetUpdates(t *testing.T) {
	b := getBot(t)
	_, err := b.GetUpdates(&tgp.GetUpdatesConfig{})
	if err != nil {
		t.Error(err)
	}
}

func TestParseMode(t *testing.T) {
	b := getBot(t)
	line, err := tgp.NewHTMLMarkdown().Link("https://www.google.com", "lol")
	if err != nil {
		t.Error(err)
	}
	m := &tgp.SendMessageConfig{
		ChatID: int64(TestChatID),
		Text:   line,
	}
	_, err = b.SendMessageable(m)
	if err != nil {
		t.Error(err)
	}
}

func TestSetWebhook(t *testing.T) {
	b := getBot(t)
	resp, err := b.SetWebhook(tgp.NewSetWebhook(WebhookURL))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
}

func TestSetCommands(t *testing.T) {
	// OK
	b := getBot(t)
	cmd := &objects.BotCommand{Command: "31321", Description: "ALLOO"}
	ok, err := b.SetMyCommands(tgp.NewSetMyCommands(cmd))
	if err != nil {
		t.Error(err, ok)
	}
	cmds, err := b.GetMyCommands(&tgp.GetMyCommandsConfig{})
	if err != nil {
		t.Error(err, cmds)
	}
	t.Log(cmds)
	for _, c := range cmds {
		if c.Command == cmd.Command && c.Description == cmd.Description {
			t.Skip("Original: ", cmd, "From tg: ", c)
			return
		}
	}
	t.Error("Command which getted from telegram, is not same as original, Original: ", cmd)
}
