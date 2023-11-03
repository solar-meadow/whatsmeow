package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type MyClient struct {
	client         *whatsmeow.Client
	mu             sync.Mutex
	blockList      map[string]struct{}
	eventHandlerID uint32
}

type typeMess string

const (
	received typeMess = "received"
	posted   typeMess = "posted"
	reported typeMess = "reported"
)

type MyMessage struct {
	ChatID  string
	UserID  string
	Text    string
	EvMes   *events.Message
	mesType typeMess
}

func NewClient() (*MyClient, error) {
	client, err := WAConnect()
	if err != nil {
		return nil, err
	}

	return &MyClient{
		client:    client,
		blockList: make(map[string]struct{}),
	}, nil

}

func (mycli *MyClient) Disconnect() {
	mycli.client.Disconnect()
}

func (mycli *MyClient) Register() {
	mycli.eventHandlerID = mycli.client.AddEventHandler(mycli.myEventHandler)
}

func (mycli *MyClient) myEventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		//		fmt.Println(v)
		if mycli.shouldProcessMessage(v) {
			if err := mycli.processMessage(v); err != nil {
				log.Println(err)
			}
		}
	}
}

func (mycli *MyClient) shouldProcessMessage(v *events.Message) bool {
	return v.Info.Chat.User == os.Getenv("TEST_ID") || v.Info.Chat.User == os.Getenv("GROUP_ID")
}

func (mycli *MyClient) processMessage(v *events.Message) error {
	text := strings.ReplaceAll(v.Message.GetConversation(), "\n", " ")
	if text == "" {
		text = "none"
	}
	number, status := ExtractNumber(text)
	mycli.loggingMessage(&MyMessage{Text: text, EvMes: v, mesType: received})
	_, exist := mycli.blockList[number]

	if number == v.Info.Sender.User {
		exist = false
	}
	if status && !exist {
		message, err := GetRequestSmcs(number)
		if err == fmt.Errorf(ErrNoUserHistory) {
			message = &ErrNoUserHistory
			if sendErr := mycli.sendMessage(&MyMessage{
				ChatID: v.Info.Chat.User,
				UserID: v.Info.Sender.User,
				Text:   *message,
				EvMes:  v,
			}); sendErr != nil {
				return sendErr
			}
		} else if err != nil {
			fmt.Println("worked error")
			if sendErr := mycli.sendReport(&MyMessage{
				UserID: v.Info.Sender.User,
				Text:   err.Error(),
				EvMes:  v,
			}); sendErr != nil {
				return sendErr
			}
		} else if err == nil && !exist {
			if sendErr := mycli.sendReport(&MyMessage{
				ChatID: v.Info.Chat.User,
				UserID: v.Info.Sender.User,
				Text:   *message,
				EvMes:  v,
			}); sendErr != nil {
				return sendErr
			}
		}
	} else if status && exist {
		if sendErr := mycli.sendReport(&MyMessage{
			ChatID: v.Info.Chat.User,
			UserID: v.Info.Sender.User,
			Text:   ErrForbidden + " " + number,
			EvMes:  v,
		}); sendErr != nil {
			return sendErr
		}
	}
	return nil
}

func (mycli *MyClient) sendReport(message *MyMessage) error {
	msg := &proto.Message{
		Conversation: &message.Text,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := mycli.client.SendMessage(ctx, types.JID{
		User:   message.UserID,
		Server: types.DefaultUserServer,
	}, msg)

	if err == nil {
		mycli.loggingMessage(&MyMessage{Text: message.Text, EvMes: message.EvMes, mesType: reported, UserID: message.UserID})
	}
	return err
}

func (mycli *MyClient) sendMessage(message *MyMessage) error {
	var msg = &proto.Message{
		ExtendedTextMessage: &proto.ExtendedTextMessage{
			Text: &message.Text,
			ContextInfo: &proto.ContextInfo{
				StanzaId:      &message.EvMes.Info.ID,
				Participant:   &message.EvMes.Info.Sender.User,
				QuotedMessage: message.EvMes.Message,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, err := mycli.client.SendMessage(ctx, message.EvMes.Info.Chat, msg)
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
		log.Printf("[posted message]: %s - [to]: %s\n", message.Text, message.UserID)
	case reported:
		log.Printf("[reported message]: %s - [to]: %s", message.Text, message.UserID)
	}
}

func WAConnect() (*whatsmeow.Client, error) {
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
	return client, nil
}

func (mycli *MyClient) UpdateAllStaff() error {
	info, err := mycli.client.GetGroupInfo(types.JID{
		User:   os.Getenv("GROUP_ID"),
		Server: types.GroupServer,
	})
	if err != nil {
		return err
	}
	mycli.mu.Lock()
	mycli.blockList = make(map[string]struct{})
	mycli.mu.Unlock()
	staffListToFile := ""
	for _, v := range info.Participants {
		mycli.mu.Lock()
		mycli.blockList[v.JID.User] = struct{}{}
		staffListToFile += v.JID.User + "\n"
		mycli.mu.Unlock()
	}

	WriteToFile("block.txt", staffListToFile)

	fmt.Println("STAFF LIST UPDATED")
	return nil
}
