package Gpt

import (
	"context"
	"fmt"
	openai "github.com/sashabaranov/go-openai"
)

func Initialize(key string) *openai.Client {
	client := openai.NewClient(key)
	return client
}

func GenerateChat(client *openai.Client, ctx context.Context, prompt string) (string, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[GPT]Recovered from a potential panic:", err)
			// Log the error for further debugging
			// Potentially return a default value or an error
		}
	}()

	// fineTuningJob, err := client.RetrieveFineTuningJob(ctx, "ftjob-02lsB2S7eg3d13rMigFrIjo6")
	// if err != nil {
	// 	fmt.Printf("Getting fine tune model error: %v\n", err)
	// 	return "", err
	// }

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	data := resp.Choices[0].Message.Content
	return data, nil
}
