package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"main/functions"
)

var (
	// Menu texts

	bot *tgbotapi.BotAPI

	config struct {
		Token    string `json:"token"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}
)

func main() {

	config = functions.CreateConfig()

	var err error
	bot, err = tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)
	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	switch {
	// Handle messages
	case update.Message != nil:
		handleMessage(update.Message)
		break
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	var msg tgbotapi.MessageConfig
	if text == "" {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Пустое сообщение")
		_, _ = bot.Send(msg)
		return
	}

	if !strings.HasPrefix(text, "https://music.yandex.ru/album") {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Url должен начинаться с https://music.yandex.ru/album")
		_, _ = bot.Send(msg)
		return
	}
	myUrl, err := url.Parse(text)
	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Ошибка парсинга урла")
		_, _ = bot.Send(msg)
		fmt.Println(err)
		return
	}

	c := exec.Command("curl", fmt.Sprintf("https://api.music.yandex.net/tracks/%s/similar", path.Base(myUrl.Path)))
	out, err := c.Output()
	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Ошибка запроса")
		_, _ = bot.Send(msg)
		fmt.Println(string(out))
		fmt.Println(err)
		return
	}

	var trackInfo TrackInfo
	err = json.Unmarshal(out, &trackInfo)
	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Ошибка при парсинге")
		_, _ = bot.Send(msg)
		fmt.Println(err)
		return
	}
	track := trackInfo.Result.Track

	reply := `
Track  - <a href="%s">%s</a>
Album  - <a href="%s">%s</a> - %d, %s
Artist - <a href="%s">%s</a>`
	msg = tgbotapi.NewMessage(message.Chat.ID,
		fmt.Sprintf(reply,
			text,
			track.Title,
			"https://music.yandex.ru/album/"+strconv.Itoa(track.Albums[0].Id),
			track.Albums[0].Title,
			track.Albums[0].Year,
			track.Albums[0].Genre,
			"https://music.yandex.ru/artist/"+strconv.Itoa(track.Artists[0].Id),
			track.Artists[0].Name,
		))
	msg.ParseMode = tgbotapi.ModeHTML
	_, err = bot.Send(msg)
	if err != nil {
		msg = tgbotapi.NewMessage(message.Chat.ID, "Ошибка при отправке сообщения")
		_, _ = bot.Send(msg)
		fmt.Println(err)
	}

}
