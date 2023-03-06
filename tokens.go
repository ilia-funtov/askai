package main

import (
	"fmt"
	"strings"
)

func calcModelMaxResponseSize(prompt string, modelMaxTokens int) (int, error) {
	splitPrompt := strings.Fields(prompt)
	tokensInPrompt := int(len(splitPrompt) * 4.0 / 3.0) // rough estimation
	maxResponseTokens := int(modelMaxTokens - tokensInPrompt)
	if maxResponseTokens <= 0 {
		return 0, fmt.Errorf("too many tokens to process")
	}

	return maxResponseTokens, nil
}
