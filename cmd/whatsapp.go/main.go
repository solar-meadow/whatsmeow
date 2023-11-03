package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/solar-meadow/getCode/pkg"
)

func main() {
	fmt.Println("1")
	logFile, err := pkg.InitLogger()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("2")
	defer logFile.Close()
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("4")
	client, err := pkg.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("3")
	client.Register()
	if err := client.UpdateAllStaff(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("5")
	duration := time.Hour * 24
	ticker := time.NewTicker(duration)

	go func() {
		for range ticker.C {
			client.UpdateAllStaff()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}
