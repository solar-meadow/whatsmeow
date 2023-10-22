package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solar-meadow/getCode/pkg"
)

// type MyClient struct {
//     WAClient *whatsmeow.Client
//     eventHandlerID uint32
// }

// func (mycli *MyClient) register() {
//     mycli.eventHandlerID = mycli.WAClient.AddEventHandler(mycli.myEventHandler)
// }

// func (mycli *MyClient) myEventHandler(evt interface{}) {
//     // Handle event and access mycli.WAClient
// }

func main() {
	if err := pkg.InitLogger(); err != nil {
		log.Fatal(err)
	}
	if err := (godotenv.Load(".env")); err != nil {
		log.Fatal(err)
	}
	client, err := pkg.WAConnect()
	if err != nil {
		log.Fatal(err)
	}
	client.Register()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.WAClient.Disconnect()
}
