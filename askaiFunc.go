package main

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UserMessage struct {
	Prompt  string
	Context string
}

type AIEngine interface {
	AskAI(message UserMessage, model string, apiKey string) ([]string, error)
	GetMaxTokenLimit(model string) int
	GetTokenizationEncoding(model string) (string, error)
	CalcTokenNum(model string, text string) (int, error)
	SplitText(model string, text string, maxTokenLen int) ([]string, error)
}

const (
	errorMessageCalcTokenNum = "AIEngine.CalcTokenNum failed: %w"
)

type EngineCallResult struct {
	engineKey string
	responses []string
	err       error
}

var engineMap = map[string]AIEngine{
	"openai": &OpenAIEngine{},
	"cohere": &CohereEngine{},
}

func (message UserMessage) GetFullPrompt() string {
	return makeFullPrompt(message.Prompt, message.Context)
}

func askAI(engines []string, message UserMessage, config ProgramConfig) (map[string][]string, error) {
	if len(engines) == 0 {
		return nil, fmt.Errorf("no AI engine found")
	}

	result := make(map[string][]string)

	processEngine := func(engine string, message UserMessage, apiKeys map[string]string) EngineCallResult {
		aiProvider, aiModel, err := splitEngineName(engine)
		if err != nil {
			return EngineCallResult{"", nil, err}
		}

		return callAIEngine(aiProvider, aiModel, message, config)
	}

	if len(engines) == 1 {
		callResult := processEngine(engines[0], message, config.APIKeys)
		if callResult.err != nil {
			return nil, callResult.err
		}

		result[callResult.engineKey] = callResult.responses
		return result, nil
	}

	resultChannel := make(chan EngineCallResult)
	processEngineAsync := func(engine string, message UserMessage, apiKeys map[string]string) {
		callResult := processEngine(engine, message, apiKeys)
		resultChannel <- callResult
	}

	for _, engine := range engines {
		go processEngineAsync(engine, message, config.APIKeys)
	}

	for i := 0; i != len(engines); i++ {
		callResult, ok := <-resultChannel
		if ok {
			result[callResult.engineKey] = callResult.responses
		}
	}

	return result, nil
}

func callAIEngine(aiProvider string, aiModel string, message UserMessage, config ProgramConfig) EngineCallResult {
	if aiModel == "" {
		var exists bool
		aiModel, exists = config.ProviderModel[aiProvider]
		if !exists {
			return EngineCallResult{"", nil, fmt.Errorf("no provider model found for %s", aiProvider)}
		}
	}

	engineKey := fmt.Sprintf("%s:%s", aiProvider, aiModel)

	engine, exists := engineMap[aiProvider]
	if !exists {
		return EngineCallResult{engineKey, nil, fmt.Errorf("no engine found for %s", aiProvider)}
	}

	apiKey, exists := config.APIKeys[aiProvider]
	if !exists {
		return EngineCallResult{"", nil, fmt.Errorf("no API key found for %s", aiProvider)}
	}

	prompt := message.GetFullPrompt()
	log.Infof("Asking %s: %s", engineKey, prompt)

	tokensInFullPrompt, err := engine.CalcTokenNum(aiModel, prompt)
	if err != nil {
		return EngineCallResult{engineKey, nil, err}
	}

	tokenLimit := engine.GetMaxTokenLimit(aiModel)

	if tokensInFullPrompt > tokenLimit {
		log.Infof("Full prompt is too long, shortening it to %d tokens at max", tokenLimit)

		pMessage, err := shortenMessage(message, tokenLimit, engine, aiModel, apiKey, config.SummarizePrompt)
		if err != nil {
			return EngineCallResult{engineKey, nil, err}
		}

		message = *pMessage
	}

	responses, err := engine.AskAI(message, aiModel, apiKey)
	if err == nil {
		log.Tracef("Engine %s returned response: %v", engineKey, responses)
	} else {
		log.Errorf("Engine %s returned error: %v", engineKey, err)
	}

	return EngineCallResult{engineKey, responses, err}
}

func shortenMessage(message UserMessage, tokenLimit int, engine AIEngine, aiModel string, apiKey string, tldrPrompt string) (*UserMessage, error) {
	tokensInPrompt, err := engine.CalcTokenNum(aiModel, message.Prompt)
	if err != nil {
		return nil, fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	tokensInContext, err := engine.CalcTokenNum(aiModel, message.Context)
	if err != nil {
		return nil, fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	shortenedPrompt, err := shortenText(message.Prompt, tokenLimit-tokensInContext-1, engine, aiModel, apiKey, tldrPrompt)
	if err != nil {
		return nil, err
	}

	shortenedContext, err := shortenText(message.Context, tokenLimit-tokensInPrompt-1, engine, aiModel, apiKey, tldrPrompt)
	if err != nil {
		return nil, err
	}

	message = UserMessage{Prompt: shortenedPrompt, Context: shortenedContext}
	return &message, nil
}

func shortenText(text string, maxTokens int, engine AIEngine, aiModel string, apiKey string, tldrPrompt string) (string, error) {
	if text == "" || maxTokens <= 0 {
		return "", nil
	}

	tokensNum, err := engine.CalcTokenNum(aiModel, text)
	if err != nil {
		return "", fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	if tokensNum <= maxTokens {
		return text, nil
	}

	log.Tracef("Shortening text: %s", text)

	tldrLen, err := engine.CalcTokenNum(aiModel, tldrPrompt)
	if err != nil {
		return "", fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	numBlocks := int(math.Ceil(float64(tokensNum) / float64(maxTokens)))
	blockTokensNum := (tokensNum / numBlocks) - (tldrLen + 1)
	parts, err := engine.SplitText(aiModel, text, blockTokensNum)
	if err != nil {
		return "", fmt.Errorf("AIEngine.SplitText failed: %w", err)
	}

	shortenedText, err := shortenTextParts(parts, engine, aiModel, apiKey, tldrPrompt)
	if err != nil {
		return "", err
	}

	if shortenedText == "" {
		return "", fmt.Errorf("text content was completely lost as a result of shortening")
	}

	shortLen, err := engine.CalcTokenNum(aiModel, shortenedText)
	if err != nil {
		return "", fmt.Errorf(errorMessageCalcTokenNum, err)
	}

	if shortLen > maxTokens {
		return shortenText(shortenedText, maxTokens, engine, aiModel, apiKey, tldrPrompt)
	}

	log.Tracef("Shortened text: %s", shortenedText)

	return shortenedText, nil
}

func shortenTextParts(parts []string, engine AIEngine, aiModel string, apiKey string, tldrPrompt string) (string, error) {
	shortenedText := ""

	for _, part := range parts {
		log.Tracef("Asking to shorten part: %s", part)

		message := UserMessage{Prompt: tldrPrompt, Context: part}
		responses, err := engine.AskAI(message, aiModel, apiKey)
		if err != nil {
			log.Errorf("Engine %s returned error: %v", reflect.TypeOf(engine), err)

			return "", fmt.Errorf("could not shorten text: %w", err)
		}

		for _, response := range responses {
			if len(response) > 0 {
				if len(shortenedText) > 0 {
					shortenedText += " "
				}
				shortenedText += response
			}
		}
	}

	return strings.TrimSpace(shortenedText), nil
}
