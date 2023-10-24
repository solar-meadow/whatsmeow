package main

import (
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
	logFile, err := pkg.InitLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	client, err := pkg.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	client.Register()
	if err := client.UpdateAllStaff(); err != nil {
		log.Fatal(err)
	}
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
