package model

import "fmt"

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

func (update *Update) Info() string {
	return fmt.Sprintf("UpdateID: [%d] | ChatID: [%d] | Message: [%s]\n", update.UpdateId, update.Message.Chat.ChatId, update.Message.Text)
}
