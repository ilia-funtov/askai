package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/pkoukk/tiktoken-go"
)

type Tokenizer struct {
	encoding string
}

func NewTokenizer(ecoding string) *Tokenizer {
	return &Tokenizer{
		encoding: ecoding,
	}
}

func (t *Tokenizer) CalcTokenNum(text string) (int, error) {
	if t.encoding == "" {
		return calcTokenNumRoughly(text), nil
	}

	return calcTokenNumExact(text, t.encoding)
}

func (t *Tokenizer) CalcModelMaxResponseSize(prompt string, modelMaxTokens int) (int, error) {
	tokensInPrompt, err := t.CalcTokenNum(prompt)
	if err != nil {
		return 0, fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	maxResponseTokens := modelMaxTokens - tokensInPrompt
	if maxResponseTokens <= 0 {
		return 0, fmt.Errorf("too many tokens to process")
	}

	return maxResponseTokens, nil
}

func (t *Tokenizer) SplitText(text string, maxTokenLen int) ([]string, error) {
	if maxTokenLen == 0 {
		return []string{}, nil
	}

	textTokenNum, err := t.CalcTokenNum(text)
	if err != nil {
		return []string{}, fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	if textTokenNum == 0 {
		return []string{}, nil
	}

	sentences := splitTextIntoSentences(text)

	if len(sentences) == 0 {
		return []string{text}, nil
	}

	parts := make([]string, 0)

	part := ""
	partSize := 0

	for _, sentence := range sentences {
		tokenNum, err := t.CalcTokenNum(sentence)
		if err != nil {
			return []string{}, fmt.Errorf(errorMessageCalcTokenNum, err)
		}

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

	return parts, nil
}

func calcTokenNumExact(text string, encoding string) (int, error) {
	if len(text) == 0 {
		return 0, nil
	}

	tke, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return 0, fmt.Errorf("tiktoken.GetEncoding: %w", err)
	}

	tokens := tke.Encode(text, nil, nil)
	return len(tokens), nil
}

func calcTokenNumRoughly(text string) int {
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

	for _, r := range runes[1:] {
		if unicode.IsLetter(r) {
			if !unicode.IsLetter(charLast) {
				letterTokenNum++
			}
		} else {
			nonLettersNum++
		}
		charLast = r
	}

	const wordsToTokensRatio = 4.0 / 3.0
	tokensInText := int(math.Ceil(float64(letterTokenNum)*wordsToTokensRatio)) + nonLettersNum // rough estimation
	return tokensInText
}

func splitTextIntoSentences(text string) []string {
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

	return sentences
}
