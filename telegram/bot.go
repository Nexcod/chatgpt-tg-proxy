package telegram

import (
	"encoding/json"
	"net/http"
)

type Bot struct {
	token string
}

func NewBot(token string) *Bot {
	return &Bot{token: token}
}

func (b *Bot) HandleWebhook(path string) (<-chan Update, error) {
	updatesChannel := make(chan Update)

	http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
		var update Update
		if err := json.NewDecoder(request.Body).Decode(&update); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}

		updatesChannel <- update
	})

	return updatesChannel, nil
}

func (b *Bot) ListenWebhook(address string) error {
	if err := http.ListenAndServe(address, nil); err != nil {
		return err
	}

	return nil
}
