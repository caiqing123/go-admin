package openai

import (
	"context"
	"log"

	"github.com/sashabaranov/go-openai"
)

type ChatGptTool struct {
	Secret string
	Client *openai.Client
}
type Gpt3Dot5Message openai.ChatCompletionMessage

func NewChatGptTool(secret string) *ChatGptTool {
	config := openai.DefaultConfig(secret)
	//config.BaseURL = "https://openai.geekr.cool/v1"
	config.BaseURL = "https://cheap-weasel-97.deno.dev/v1"
	client := openai.NewClientWithConfig(config)
	//client := openai.NewClient(secret)
	return &ChatGptTool{
		Secret: secret,
		Client: client,
	}
}

/**
调用gpt3.5接口
*/
func (this *ChatGptTool) ChatGPT3Dot5Turbo(messages []Gpt3Dot5Message) (string, error) {
	reqMessages := make([]openai.ChatCompletionMessage, 0)
	for _, row := range messages {
		reqMessage := openai.ChatCompletionMessage{
			Role:    row.Role,
			Content: row.Content,
			Name:    row.Name,
		}
		reqMessages = append(reqMessages, reqMessage)
	}
	resp, err := this.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: reqMessages,
		},
	)

	if err != nil {
		log.Println("ChatGPT3Dot5Turbo error: ", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
