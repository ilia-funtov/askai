package main

import (
	"context"
	"fmt"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func askOpenAIChatCompletionModel(prompt string, model string, apiKey string) ([]string, error) {
	const MaxTokensGPT3dot5Chat = 4096

	tokensInPrompt := len(prompt)
	maxTokens := int(MaxTokensGPT3dot5Chat - tokensInPrompt)
	if maxTokens <= 0 {
		return nil, fmt.Errorf("too many tokens to process")
	}

	ctx := context.Background()
	c := gogpt.NewClient(apiKey)

	message := gogpt.ChatCompletionMessage{Role: "user", Content: prompt}
	messages := []gogpt.ChatCompletionMessage{message}

	request := gogpt.ChatCompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
	}

	response, err := c.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, err
	}

	var responses []string
	for _, choice := range response.Choices {
		if choice.Message.Role == "assistant" {
			responses = append(responses, choice.Message.Content)
		}
	}

	return responses, nil
}

func askOpenAICompletionModel(prompt string, model string, apiKey string) ([]string, error) {
	const MaxTokensGPT3dot5 = 4000

	tokensInPrompt := len(prompt)
	maxTokens := int(MaxTokensGPT3dot5 - tokensInPrompt)
	if maxTokens <= 0 {
		return nil, fmt.Errorf("too many tokens to process")
	}

	ctx := context.Background()
	c := gogpt.NewClient(apiKey)

	message := gogpt.CompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Prompt:    prompt,
	}

	response, err := c.CreateCompletion(ctx, message)
	if err != nil {
		return nil, err
	}

	var responses []string
	for _, choice := range response.Choices {
		responses = append(responses, choice.Text)
	}

	return responses, nil
}

func askOpenAI(prompt string, model string, apiKey string) ([]string, error) {
	if model == gogpt.GPT3Dot5Turbo || model == gogpt.GPT3Dot5Turbo0301 {
		return askOpenAIChatCompletionModel(prompt, model, apiKey)
	}

	return askOpenAICompletionModel(prompt, model, apiKey)
}
