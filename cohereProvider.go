package main

import (
	"fmt"

	cohere "github.com/cohere-ai/cohere-go"
)

const MaxTokensCohere = 2048

func askCohere(message UserMessage, model string, apiKey string) ([]string, error) {
	prompt := message.GetFullPrompt()

	tok := NewTokenizer("")
	maxTokens, err := tok.CalcModelMaxResponseSize(prompt, MaxTokensCohere)
	if err != nil {
		return nil, err
	}

	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, fmt.Errorf("could not create cohere client: %w", err)
	}

	response, err := client.Generate(
		cohere.GenerateOptions{
			Prompt:    prompt,
			Model:     model,
			MaxTokens: uint(maxTokens),
		})

	if err != nil {
		return nil, fmt.Errorf("cohere could not generate text completion: %w", err)
	}

	result := make([]string, 0, len(response.Generations))
	for _, generation := range response.Generations {
		result = append(result, generation.Text)
	}

	return result, nil
}

type CohereEngine struct{}

func (e *CohereEngine) AskAI(message UserMessage, model string, apiKey string) ([]string, error) {
	return askCohere(message, model, apiKey)
}

func (e *CohereEngine) GetMaxTokenLimit(model string) int {
	return MaxTokensCohere
}

func (e *CohereEngine) GetTokenizationEncoding(model string) (string, error) {
	return "", nil
}

func (e *CohereEngine) CalcTokenNum(model string, text string) (int, error) {
	tok := NewTokenizer("")
	return tok.CalcTokenNum(text)
}

func (e *CohereEngine) SplitText(model string, text string, maxTokenLen int) ([]string, error) {
	tok := NewTokenizer("")
	return tok.SplitText(text, maxTokenLen)
}
