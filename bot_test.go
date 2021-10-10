package tgp

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/pikoUsername/tgp/objects"
)

var (
	parseMode     = "HTML"
	testChatID, _ = strconv.ParseInt(os.Getenv("test_chat_id"), 10, 64)
	fileDirectory = "./.sandbox"
	saveFile      = path.Join(fileDirectory, "file")
	webhookURL    = ""
	testToken     = os.Getenv("TEST_TOKEN")

	// here could be any image, file, anthing else
	downloadFromURL = "https://random.imagecdn.app/500/150"
)

func failIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func panicIfErr(t *testing.T, err error) {
	if err != nil {
		panic(err)
	}
}

func getBot(t *testing.T) *Bot {
	b, err := NewBot(testToken, parseMode, nil)
	if err != nil {
		t.Fatal(err)
	}
	b.Debug = true
	return b
}

func TestCheckToken(t *testing.T) {
	b, err := NewBot("bla:bla", "HTML", nil)
	if err != nil && b == nil {
		t.Fatal(err)
	}
}

func TestDownloadFile(t *testing.T) {
	b := getBot(t)
	dir, err := os.Open(fileDirectory)
	if err == os.ErrNotExist {
		os.Mkdir(fileDirectory, 0777)
		dir, err = os.Open(fileDirectory)
		if err != nil {
			t.Fatal(err)
		}
	}
	stat, err := dir.Stat()
	if err != nil {
		t.Fatal(err)
	}
	if !stat.IsDir() {
		t.Fatal("Sorry but -"+fileDirectory, "Is file, delete file and try again!")
	}

	f, err := os.OpenFile(saveFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		t.Fatal(err)
	}
	err = b.DownloadFile(downloadFromURL, f, true)
	if err != nil {
		t.Fatal(err)
	}
	stat, err = f.Stat()
	if err != nil {
		t.Fatal(err)
	}
	bs := make([]byte, stat.Size())
	f.Read(bs)
	if bs == nil || strings.Compare(string(bs), "") == -1 && downloadFromURL != "" {
		t.Fatal(
			"Cannot download file from ethernet, debug: file - ", bs,
			", URL: ", downloadFromURL, ", Directory: ", stat.Name(), ", Bot: ", b)
	}
}

func TestGetMe(t *testing.T) {
	b := getBot(t)
	u, err := b.GetMe()
	if err != nil {
		t.Fatal(err)
	}
	if b.Me == nil {
		t.Fatal("Me is empty pointer")
	}
	if b.Me.ID != u.ID {
		t.Fatal("Getted User is defferent from bot user")
	}
}

func TestGetUpdates(t *testing.T) {
	b := getBot(t)
	_, err := b.GetUpdates(&GetUpdatesConfig{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseMode(t *testing.T) {
	b := getBot(t)
	line, err := NewHTMLMarkdown().Link("https://www.google.com", "lol")
	if err != nil {
		t.Fatal(err)
	}
	m := &SendMessageConfig{
		ChatID: int64(testChatID),
		Text:   line,
	}
	_, err = b.SendMessageable(m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetWebhook(t *testing.T) {
	b := getBot(t)
	resp, err := b.SetWebhook(NewSetWebhook(webhookURL))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(resp)
}

func TestSetCommands(t *testing.T) {
	// OK
	b := getBot(t)
	cmd := &objects.BotCommand{Command: "31321", Description: "ALLOO"}
	ok, err := b.SetMyCommands(NewSetMyCommands(cmd))
	if err != nil {
		t.Fatal(err, ok)
	}
	cmds, err := b.GetMyCommands(&GetMyCommandsConfig{})
	if err != nil {
		t.Fatal(err, cmds)
	}
	t.Log(cmds)
	for _, c := range cmds {
		if c.Command == cmd.Command && c.Description == cmd.Description {
			t.Skip("Original: ", cmd, "From tg: ", c)
			return
		}
	}
	t.Fatal("Command which getted from telegram, is not same as original, Original: ", cmd)
}
