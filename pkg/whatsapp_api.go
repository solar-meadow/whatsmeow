package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type MyClient struct {
	WAClient       *whatsmeow.Client
	eventHandlerID uint32
}

type typeMess string

const (
	received typeMess = "received"
	posted   typeMess = "posted"
)

type MyMessage struct {
	ChatID  string
	UserID  string
	Text    string
	EvMes   *events.Message
	mesType typeMess
}

func (mycli *MyClient) Register() {
	mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.myEventHandler)
}

func (mycli *MyClient) myEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		if v.Info.Chat.User == os.Getenv("TEST_ID") { // v.Info.Chat.User == os.Getenv("GROUP_ID")
			text := strings.ReplaceAll(v.Message.GetConversation(), "\n", " ")
			if text == "" {
				text = "none"
			}
			number, status := ExtractPhoneNumber(text)
			mycli.loggingMessage(&MyMessage{Text: text, EvMes: v, mesType: received})
			if status {
				message, err := GetRequestSmcs(context.Background(), number)
				if err == fmt.Errorf(ErrNoUserHistory) {
					text = ErrNoUserHistory
				} else if err != nil {
					fmt.Println("worked")
					if err := mycli.sendMessage(&MyMessage{
						UserID: os.Getenv("MY_ID"),
						Text:   err.Error(),
						EvMes:  v,
					}); err != nil {
						log.Println(err)
					}
				}

				if err := mycli.sendMessage(&MyMessage{
					ChatID: v.Info.Chat.User,
					UserID: v.Info.Sender.User,
					Text:   *message,
					EvMes:  v,
				}); err != nil {
					log.Println(err)
				}

			}
		}
	}
}

func (mycli *MyClient) sendMessage(message *MyMessage) error {
	_, err := mycli.WAClient.SendMessage(context.Background(), types.JID{
		User:   message.UserID,
		Server: types.DefaultUserServer,
	}, &proto.Message{
		Conversation: &message.Text,
	})
	if err == nil {
		mycli.loggingMessage(&MyMessage{Text: message.Text, EvMes: message.EvMes, mesType: posted, UserID: message.UserID})
	}
	return err
}

func (mycli *MyClient) loggingMessage(message *MyMessage) {
	switch message.mesType {
	case received:
		fmt.Printf("[received message]: %s - [sender]: %s - [chat_id]: %s\n", message.Text, message.EvMes.Info.Sender.User, message.EvMes.Info.Chat.User)
	case posted:
		fmt.Printf("[posted message]: %s - [to]: %s\n", message.Text, message.UserID)
	}
}

func WAConnect() (*MyClient, error) {
	container, err := sqlstore.New("sqlite3", "file:wapp.db?_foreign_keys=on", waLog.Noop)
	if err != nil {
		return nil, err
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, err
	}
	client := whatsmeow.NewClient(deviceStore, waLog.Noop)
	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			return nil, err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err := client.Connect()
		if err != nil {
			return nil, err
		}
	}
	return &MyClient{
		WAClient: client,
	}, nil
}
