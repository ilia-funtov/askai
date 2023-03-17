package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitTextOne(t *testing.T) {
	const text = "One sentence."

	tokenLen := calcTokenNum(text)
	parts := splitText(text, tokenLen)

	assert.Equal(t, len(parts), 1)
	assert.Equal(t, parts[0], text)
}

func TestSplitTextTwo(t *testing.T) {
	const s1 = "Sentence one."
	const s2 = "Sentence two."
	const text = s1 + s2

	tokenLen := calcTokenNum(text)
	parts := splitText(text, tokenLen/2)

	assert.Equal(t, len(parts), 2)
	assert.Equal(t, parts[0], s1)
	assert.Equal(t, parts[1], s2)
}

func TestSplitTextThree(t *testing.T) {
	const s1 = "Sentence one."
	const s2 = "Sentence two."
	const s3 = "Sentence three."
	const text = s1 + s2 + s3

	tokenLen := calcTokenNum(text)
	parts := splitText(text, tokenLen/3)

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

	tokenLen := calcTokenNum(text)
	parts := splitText(text, tokenLen/5)

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

	tokenLen := calcTokenNum(text)
	parts := splitText(text, tokenLen/4)

	recovered := ""
	for _, part := range parts {
		recovered += part
	}

	assert.Equal(t, text, recovered)
}

func TestSplitTextOnlyPunctuation(t *testing.T) {
	const text = "?????!!!!!"

	parts := splitText(text, 1)

	assert.Equal(t, len(parts), len(text))
}

func TestSplitTextEmpty(t *testing.T) {
	{
		parts := splitText("", 10)
		assert.Equal(t, len(parts), 0)
	}

	{
		parts := splitText("text", 0)
		assert.Equal(t, len(parts), 0)
	}
}

func TestSplitTextWord(t *testing.T) {
	parts := splitText("word", 10)
	assert.Equal(t, len(parts), 1)
	assert.Equal(t, parts[0], "word")
}

func TestSplitTextEllipsis(t *testing.T) {
	parts := splitText("...", 10)
	assert.Equal(t, len(parts), 1)
	assert.Equal(t, parts[0], "...")
}
