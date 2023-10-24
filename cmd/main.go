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

func main() {
	logFile, err := pkg.InitLogger()
	if err != nil {
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
	logFile.Close()
	client.WAClient.Disconnect()
}
