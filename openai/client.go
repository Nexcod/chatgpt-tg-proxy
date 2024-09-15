package openai

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const OpenAPIURL = "https://api.openai.com/v1/chat/completions"

type Client struct {
	apiToken string
}

func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
	}
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func (c *Client) CreateChatCompletion(
	chatCompletionRequest ChatCompletionRequest,
) (chatCompletionResponse *ChatCompletionResponse, err error) {
	jsonData, err := json.Marshal(chatCompletionRequest)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", OpenAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&chatCompletionResponse)
	if err != nil {
		return
	}

	return
}
