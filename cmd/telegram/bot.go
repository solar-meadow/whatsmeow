package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/solar-meadow/getCode/model"
	"github.com/solar-meadow/getCode/pkg"
)

const BotApi = "https://api.telegram.org/bot"

func main() {
	botToken, err := pkg.CheckEnvExist("BOT_TOKEN")
	if err != nil {
		log.Println(err)
		return
	}
	botUrl := BotApi + botToken
	client := NewHttpClient(botUrl)
	botName, err := pkg.CheckEnvExist("BOT_NAME")
	if err != nil {
		botName = "Bot_Name"
	}
	log.Println("Bot running... name in telegram:", botName)
	for {
		updates, err := client.GetUpdates()
		if err != nil {
			log.Println(err)

			return
		}
		for _, update := range updates {

			err = client.Respond(update)
			if err != nil {
				log.Println(err)
				return
			}
			client.Offset = update.UpdateId + 1
		}

	}
}

type MessageService struct {
	client    *http.Client
	botUrl    string
	Offset    int
	mu        *sync.Mutex
	blockList map[string]struct{}
}

func NewHttpClient(botUrl string) *MessageService {

	return &MessageService{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				fmt.Println(req.Response.Status)
				fmt.Println("[REDIRECT]")
				return nil
			},
			Transport: http.DefaultTransport,
			Timeout:   time.Second * 30,
		},
		botUrl:    botUrl,
		Offset:    0,
		blockList: make(map[string]struct{}),
		mu:        &sync.Mutex{},
	}
}

func (m *MessageService) GetUpdates() ([]model.Update, error) {
	resp, err := http.Get(m.botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(m.Offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println("bot updates", body)
	var restResponse model.RestResponse
	if err = json.Unmarshal(body, &restResponse); err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func (m *MessageService) Respond(update model.Update) error {
	fmt.Println(update.Info())

	var botMessage model.BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	pkg.ExtractNumber(update.Message.Text)
	switch update.Message.Text {
	case "/start":
		botMessage.Text = "Write book name: "
	default:
		//botMessage.Text = flibusta.GetBookLinks(update.Message.Text)
	}

	//	botMessage.Text = flibusta.GetBookLinks(strings.ReplaceAll(update.Message.Text, " ", "+"))

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	if _, err = m.client.Post(m.botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf)); err != nil {
		return err
	}
	return nil
}
