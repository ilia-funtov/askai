package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTextOne(t *testing.T) {
	const text = "One sentence."

	tok := NewTokenizer("")

	tokenLen, err := tok.CalcTokenNum(text)
	assert.NoError(t, err)

	parts, err := tok.SplitText(text, tokenLen)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), 1)
	assert.Equal(t, parts[0], text)
}

func TestSplitTextTwo(t *testing.T) {
	const s1 = "Sentence one."
	const s2 = "Sentence two."
	const text = s1 + s2

	tok := NewTokenizer("")

	tokenLen, err := tok.CalcTokenNum(text)
	assert.NoError(t, err)

	parts, err := tok.SplitText(text, tokenLen/2)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), 2)
	assert.Equal(t, parts[0], s1)
	assert.Equal(t, parts[1], s2)
}

func TestSplitTextThree(t *testing.T) {
	const s1 = "Sentence one."
	const s2 = "Sentence two."
	const s3 = "Sentence three."
	const text = s1 + s2 + s3

	tok := NewTokenizer("")

	tokenLen, err := tok.CalcTokenNum(text)
	assert.NoError(t, err)

	parts, err := tok.SplitText(text, tokenLen/3)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), 3)
	assert.Equal(t, parts[0], s1)
	assert.Equal(t, parts[1], s2)
	assert.Equal(t, parts[2], s3)
}

func TestSplitTextMany(t *testing.T) {
	const s1 = "How do I pass text with size larger than model input?"
	const s2 = "Split it into multiple parts!"
	const s3 = "Then pass each part to model for summarization."
	const s4 = "Then concatenate all parts together."
	const s5 = "Be like map-reduce..."

	const text = s1 + s2 + s3 + s4 + s5

	tok := NewTokenizer("")

	tokenLen, err := tok.CalcTokenNum(text)
	assert.NoError(t, err)

	parts, err := tok.SplitText(text, tokenLen/5)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), 5)
	assert.Equal(t, parts[0], s1)
	assert.Equal(t, parts[1], s2)
	assert.Equal(t, parts[2], s3)
	assert.Equal(t, parts[3], s4)
	assert.Equal(t, parts[4], s5)
}

func TestSplitTextMessyPunctuation(t *testing.T) {
	const text = "Some texts could be terrible with punctuation!!!Like this..........." +
		"But what about that??????Oooo!!!???;;;!!....???So terrible.......!!!!!"

	tok := NewTokenizer("")

	tokenLen, err := tok.CalcTokenNum(text)
	assert.NoError(t, err)

	parts, err := tok.SplitText(text, tokenLen/4)
	assert.NoError(t, err)

	recovered := ""
	for _, part := range parts {
		recovered += part
	}

	assert.Equal(t, text, recovered)
}

func TestSplitTextOnlyPunctuation(t *testing.T) {
	const text = "?????!!!!!"

	tok := NewTokenizer("")

	parts, err := tok.SplitText(text, 1)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), len(text))
}

func TestSplitTextEmpty(t *testing.T) {
	tok := NewTokenizer("")

	{
		parts, err := tok.SplitText("", 10)
		assert.NoError(t, err)
		assert.Equal(t, len(parts), 0)
	}

	{
		parts, err := tok.SplitText("text", 0)
		assert.NoError(t, err)
		assert.Equal(t, len(parts), 0)
	}
}

func TestSplitTextWord(t *testing.T) {
	tok := NewTokenizer("")

	parts, err := tok.SplitText("word", 10)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), 1)
	assert.Equal(t, parts[0], "word")
}

func TestSplitTextEllipsis(t *testing.T) {
	tok := NewTokenizer("")

	parts, err := tok.SplitText("...", 10)
	assert.NoError(t, err)

	assert.Equal(t, len(parts), 1)
	assert.Equal(t, parts[0], "...")
}
