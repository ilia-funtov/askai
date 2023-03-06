package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type AskAIFunction func(string, string, string) ([]string, error)

var engineFuncMap = map[string]AskAIFunction{
	"openai": askOpenAI,
	"cohere": askCohere,
}

type EngineCallResult struct {
	engineKey string
	responses []string
	err       error
}

func askAI(engines []string, prompt string, apiKeys map[string]string) (map[string][]string, error) {
	if len(engines) == 0 {
		return nil, fmt.Errorf("no AI engine found")
	}

	result := make(map[string][]string)

	processEngine := func(engine string, prompt string, apiKeys map[string]string) EngineCallResult {
		aiProvider, aiModel, err := splitEngineName(engine)
		if err != nil {
			return EngineCallResult{"", nil, err}
		}

		apiKey, exists := apiKeys[aiProvider]
		if !exists {
			return EngineCallResult{"", nil, fmt.Errorf("no API key found for %s", aiProvider)}
		}

		return callAIEngine(aiProvider, aiModel, prompt, apiKey)
	}

	if len(engines) == 1 {
		callResult := processEngine(engines[0], prompt, apiKeys)
		if callResult.err != nil {
			return nil, callResult.err
		}

		result[callResult.engineKey] = callResult.responses
		return result, nil
	}

	resultChannel := make(chan EngineCallResult)
	processEngineAsync := func(engine string, prompt string, apiKeys map[string]string) {
		callResult := processEngine(engine, prompt, apiKeys)
		resultChannel <- callResult
	}

	for _, engine := range engines {
		go processEngineAsync(engine, prompt, apiKeys)
	}

	for i := 0; i != len(engines); i++ {
		callResult, ok := <-resultChannel
		if ok {
			result[callResult.engineKey] = callResult.responses
		}
	}

	return result, nil
}

func callAIEngine(aiProvider string, aiModel string, prompt string, apiKey string) EngineCallResult {
	if aiModel == "" {
		var exists bool
		aiModel, exists = defaultProviderModel[aiProvider]
		if !exists {
			return EngineCallResult{"", nil, fmt.Errorf("no provider model found for %s", aiProvider)}
		}
	}

	engineKey := fmt.Sprintf("%s:%s", aiProvider, aiModel)

	engineFunc, exists := engineFuncMap[aiProvider]
	if !exists {
		return EngineCallResult{engineKey, nil, fmt.Errorf("no engine found for %s", aiProvider)}
	}

	log.Infof("Asking %s: %s", engineKey, prompt)
	responses, err := engineFunc(prompt, aiModel, apiKey)
	if err == nil {
		log.Tracef("Engine %s returned response: %v", engineKey, responses)
	} else {
		log.Errorf("Engine %s returned error: %v", engineKey, err)
	}

	return EngineCallResult{engineKey, responses, err}
}
