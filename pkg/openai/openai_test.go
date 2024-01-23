package openai

import (
	"log"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestOpen(t *testing.T) {
	gpt := NewChatGptTool("")
	message := []Gpt3Dot5Message{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "怎么聊天",
		},
	}
	res, err := gpt.ChatGPT3Dot5Turbo(message)
	log.Println(res, err)
}
