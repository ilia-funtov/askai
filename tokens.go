package main

import (
	"fmt"
	"math"
	"unicode"
)

func calcModelMaxResponseSize(prompt string, modelMaxTokens int) (int, error) {
	nonLettersNum := 0
	letterTokenNum := 0

	if len(prompt) > 0 {
		runes := []rune(prompt)

		charLast := runes[0]
		if unicode.IsLetter(charLast) {
			letterTokenNum++
		} else {
			nonLettersNum++
		}

		for _, c := range runes[1:] {
			if unicode.IsLetter(c) {
				if !unicode.IsLetter(charLast) {
					letterTokenNum++
				}
			} else {
				nonLettersNum++
			}
			charLast = c
		}
	}

	tokensInPrompt := int(math.Ceil(float64(letterTokenNum)*4.0/3.0)) + nonLettersNum + 1 // rough estimation
	maxResponseTokens := int(modelMaxTokens - tokensInPrompt)
	if maxResponseTokens <= 0 {
		return 0, fmt.Errorf("too many tokens to process")
	}

	return maxResponseTokens, nil
}
