package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

func calcTokenNum(text string) int {
	if len(text) == 0 {
		return 0
	}

	nonLettersNum := 0
	letterTokenNum := 0

	runes := []rune(text)

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

	tokensInText := int(math.Ceil(float64(letterTokenNum)*4.0/3.0)) + nonLettersNum // rough estimation
	return tokensInText
}

func calcModelMaxResponseSize(prompt string, modelMaxTokens int) (int, error) {
	tokensInPrompt := calcTokenNum(prompt)
	maxResponseTokens := modelMaxTokens - tokensInPrompt
	if maxResponseTokens <= 0 {
		return 0, fmt.Errorf("too many tokens to process")
	}

	return maxResponseTokens, nil
}

func splitText(text string, maxTokenLen int) []string {
	if maxTokenLen == 0 {
		return []string{}
	}

	textTokenNum := calcTokenNum(text)
	if textTokenNum == 0 {
		return []string{}
	}

	separators := []string{"...", ".", "!", "?"}

	sentences := make([]string, 0)

	runes := []rune(text)

	i := 0
	for j := 0; j < len(runes); j++ {
		str := runes[j:]
		for _, sep := range separators {
			if strings.HasPrefix(string(str), sep) {
				end := j + len(sep)
				sentences = append(sentences, string(runes[i:end]))
				i = end
				j = i - 1
				break
			}
		}
	}

	if len(sentences) == 0 {
		return []string{text}
	}

	parts := make([]string, 0)

	part := ""
	partSize := 0

	for _, sentence := range sentences {
		tokenNum := calcTokenNum(sentence)
		if (tokenNum + partSize) > maxTokenLen {
			if len(part) > 0 {
				parts = append(parts, part)
			}
			part = sentence
			partSize = tokenNum
		} else {
			part += sentence
			partSize += tokenNum
		}
	}

	if len(part) > 0 {
		parts = append(parts, part)
	}

	return parts
}
