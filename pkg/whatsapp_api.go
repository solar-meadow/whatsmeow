package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func WAConnect() (*whatsmeow.Client, error) {
	container, err := sqlstore.New("sqlite3", "file:wapp.db?_foreign_keys=on", waLog.Noop)
	if err != nil {
		return nil, err
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
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

func EventHandler(evt interface{}) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("error open file", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	switch v := evt.(type) {
	case *events.Message:
		ctx := context.Background()
		if v.Info.Chat.User == os.Getenv("GROUP_ID") || v.Info.Sender.User == os.Getenv("MY_ID") {
			text := strings.ReplaceAll(v.Message.GetConversation(), "\n", " ")
			if text == "" {
				text = "none"
			}
			fmt.Printf("received message: %s from sender: %s  from chat: %s\n", text, v.Info.Sender.User, v.Info.Chat.User)
			log.Printf("received message: %s from sender: %s  from chat: %s\n", text, v.Info.Sender.User, v.Info.Chat.User)
			fmt.Println(1)
			number, status := ExtractPhoneNumber(text)
			if status {
				fmt.Println(2)
				ctxSmc, cancel := context.WithTimeout(context.Background(), time.Second*4)
				defer cancel()
				code, err := GetRequestSmcs(ctxSmc, number)

				fmt.Println(*code)
				fmt.Println(3)
				if err == nil {
					fmt.Println(4)
					client, err := WAConnect()
					if err != nil {
						fmt.Println(err)
						return
					}
					fmt.Println(5)

					fmt.Println(6)
					message := fmt.Sprintf("Запрос: %s, код: %s", text, *code)
					_, err = client.SendMessage(ctx, types.JID{
						User:   v.Info.Sender.User,
						Server: types.DefaultUserServer,
					}, &proto.Message{
						Conversation: &message,
					})
					if err != nil {
						fmt.Println(err)
					}

				} else {
					client, err := WAConnect()
					if err != nil {
						fmt.Println(err)
						return
					}
					ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
					defer cancel()
					message := fmt.Sprintf("Запрос: %s, код: %s ошибка: %v", text, *code, err)
					_, err = client.SendMessage(ctx, types.JID{
						User:   os.Getenv("MY_ID"),
						Server: types.DefaultUserServer,
					}, &proto.Message{
						Conversation: &message,
					})
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(err)
				}
			} else {
				fmt.Println("invalid format")
			}
		}
	}
}
