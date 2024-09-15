package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SendMessageParams struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

func (b *Bot) SendMessage(params *SendMessageParams) (message *Message, err error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return
	}

	resp, err := http.Post("https://api.telegram.org/bot"+b.token+"/sendMessage", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&message)
	return
}

type ChatActionParams struct {
	ChatID int64  `json:"chat_id"`
	Action string `json:"action"`
}

func (b *Bot) SendChatAction(params *ChatActionParams) (err error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return
	}

	resp, err := http.Post("https://api.telegram.org/bot"+b.token+"/sendChatAction", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return nil
}
