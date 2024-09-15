package main

import (
	"fmt"
	"github.com/Nexcod/chatgpt-tg-proxy/openai"
	"github.com/Nexcod/chatgpt-tg-proxy/telegram"
	"log"
	"os"
	"slices"
	"time"
)

var allowedTelegramUsers = []string{
	"nexcod",
	"mr_sergeomorello196",
}

func handleStartCommand(update telegram.Update) {
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	bot := telegram.NewBot(tgBotToken)

	messageParams := telegram.SendMessageParams{
		ChatID: update.Message.Chat.Id,
		Text:   "ChatGPT Proxy Bot",
	}
	_, err := bot.SendMessage(&messageParams)
	if err != nil {
		log.Fatalln(err)
	}
}

func handleMessage(update telegram.Update) {
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	bot := telegram.NewBot(tgBotToken)

	if !slices.Contains(allowedTelegramUsers, update.Message.From.Username) {
		messageParams := telegram.SendMessageParams{
			ChatID: update.Message.Chat.Id,
			Text:   "Вы не добавлены в список разрешенных пользователей",
		}
		_, err := bot.SendMessage(&messageParams)
		if err != nil {
			log.Fatalln(err)
		}
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	stopTicker := make(chan bool)
	go sendTypingAction(update.Message.Chat.Id, ticker, stopTicker)

	openAPIToken := os.Getenv("OPENAPI_TOKEN")
	openaiClient := openai.NewClient(openAPIToken)
	chatCompletionRequest := openai.ChatCompletionRequest{
		Model: "gpt-4",
		Messages: []openai.Message{
			{Role: "user", Content: update.Message.Text},
		},
	}
	chatCompletionResponse, err := openaiClient.CreateChatCompletion(chatCompletionRequest)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	ticker.Stop()
	stopTicker <- true

	messageParams := telegram.SendMessageParams{
		ChatID: update.Message.Chat.Id,
		Text:   chatCompletionResponse.Choices[0].Message.Content,
	}
	_, err = bot.SendMessage(&messageParams)
	if err != nil {
		log.Fatalln(err)
	}
}

func sendTypingAction(chatId int64, ticker *time.Ticker, stopTicker chan bool) {
	tgBotToken := os.Getenv("TG_BOT_TOKEN")
	bot := telegram.NewBot(tgBotToken)

	for {
		select {
		case <-stopTicker:
			return
		case <-ticker.C:
			chatActionParams := telegram.ChatActionParams{
				ChatID: chatId,
				Action: "typing",
			}
			err := bot.SendChatAction(&chatActionParams)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func main() {
	openAPIToken := os.Getenv("OPENAPI_TOKEN")
	tgBotToken := os.Getenv("TG_BOT_TOKEN")

	if openAPIToken == "" {
		fmt.Println("OPENAPI_KEY is not set")
		os.Exit(1)
	}

	if tgBotToken == "" {
		fmt.Println("TG_BOT_TOKEN is not set")
		os.Exit(1)
	}

	bot := telegram.NewBot(tgBotToken)

	updates, err := bot.HandleWebhook("/webhook")
	if err != nil {
		log.Fatalln("Error starting webhook: ", err)
	}

	handler := telegram.NewHandler(updates)
	handler.Handle(handleStartCommand, telegram.Command("start"))
	handler.Handle(handleMessage, telegram.AnyTextMessage())
	go handler.Start()

	log.Println("Bot is running...")
	if err := bot.ListenWebhook("127.0.0.1:8080"); err != nil {
		log.Fatalf("Error listening webhook: %v", err)
	}
}
