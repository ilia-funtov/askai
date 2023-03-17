package main

import (
	"context"
	"fmt"

	gogpt "github.com/sashabaranov/go-gpt3"
)

const MaxTokensGPT3dot5Chat = 4096
const MaxTokensGPT3dot5 = 4000

func askOpenAIChatCompletionModel(message UserMessage, model string, apiKey string) ([]string, error) {
	prompt := message.GetFullPrompt()
	maxTokens, err := calcModelMaxResponseSize(prompt, MaxTokensGPT3dot5Chat)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	c := gogpt.NewClient(apiKey)

	completionMessage := gogpt.ChatCompletionMessage{Role: "user", Content: prompt}
	messages := []gogpt.ChatCompletionMessage{completionMessage}

	request := gogpt.ChatCompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Messages:  messages,
	}

	response, err := c.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("openai could not create chat completion: %w", err)
	}

	responses := make([]string, 0, len(response.Choices))
	for _, choice := range response.Choices {
		if choice.Message.Role == "assistant" {
			responses = append(responses, choice.Message.Content)
		}
	}

	return responses, nil
}

func askOpenAICompletionModel(message UserMessage, model string, apiKey string) ([]string, error) {
	prompt := message.GetFullPrompt()
	maxTokens, err := calcModelMaxResponseSize(prompt, MaxTokensGPT3dot5)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	c := gogpt.NewClient(apiKey)

	request := gogpt.CompletionRequest{
		Model:     model,
		MaxTokens: maxTokens,
		Prompt:    prompt,
	}

	response, err := c.CreateCompletion(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("openai could not create text completion: %w", err)
	}

	responses := make([]string, 0, len(response.Choices))
	for _, choice := range response.Choices {
		responses = append(responses, choice.Text)
	}

	return responses, nil
}

func askOpenAI(message UserMessage, model string, apiKey string) ([]string, error) {
	if model == gogpt.GPT3Dot5Turbo || model == gogpt.GPT3Dot5Turbo0301 {
		return askOpenAIChatCompletionModel(message, model, apiKey)
	}

	return askOpenAICompletionModel(message, model, apiKey)
}

type OpenAIEngine struct{}

func (e *OpenAIEngine) AskAI(message UserMessage, model string, apiKey string) ([]string, error) {
	return askOpenAI(message, model, apiKey)
}

func (e *OpenAIEngine) GetMaxTokenLimit(model string) int {
	switch model {
	case gogpt.GPT3Dot5Turbo:
	case gogpt.GPT3Dot5Turbo0301:
		return MaxTokensGPT3dot5Chat
	}
	return MaxTokensGPT3dot5
}
