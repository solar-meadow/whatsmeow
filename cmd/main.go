package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solar-meadow/getCode/pkg"
)

func main() {
	if err := (godotenv.Load(".env")); err != nil {
		log.Fatal(err)
	}

	// groups, err := client.GetJoinedGroups()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, v := range groups {
	// 	fmt.Println(v.GroupName, v.JID)
	// }
	// num, err := pkg.GetRequestSmcs("+77751888517")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(*num)
	client, err := pkg.WAConnect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("work1")
	client.AddEventHandler(pkg.EventHandler)
	fmt.Println("work2")
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	// defer cancel()
	// mes := "hello"
	// fmt.Println(os.Getenv("TEST_ID"))
	// _, err = client.SendMessage(ctx, types.JID{
	// 	User:   os.Getenv("TEST_ID"),
	// 	Server: types.DefaultUserServer,
	// }, &proto.Message{
	// 	Conversation: &mes,
	// })
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
