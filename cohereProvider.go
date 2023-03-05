package main

import (
	"fmt"

	cohere "github.com/cohere-ai/cohere-go"
)

func askCohere(prompt string, model string, apiKey string) ([]string, error) {
	const MaxTokensCohere = 2048

	tokensInPrompt := len(prompt)
	maxTokens := int(MaxTokensCohere - tokensInPrompt)
	if maxTokens <= 0 {
		return nil, fmt.Errorf("too many tokens to process")
	}

	client, err := cohere.CreateClient(apiKey)
	if err != nil {
		return nil, err
	}

	response, err := client.Generate(
		cohere.GenerateOptions{
			Prompt:    prompt,
			Model:     model,
			MaxTokens: uint(maxTokens),
		})

	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, generation := range response.Generations {
		result = append(result, generation.Text)
	}

	return result, nil
}
