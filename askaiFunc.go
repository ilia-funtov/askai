package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type AskAIFunction func(string, string, string) ([]string, error)

var engineMap = map[string]AskAIFunction{
	"openai": askOpenAI,
	"cohere": askCohere,
}

func askAI(engines []string, prompt string, apiKeys map[string]string) (map[string][]string, error) {
	if len(engines) == 0 {
		return nil, fmt.Errorf("no AI engine found")
	}

	result := make(map[string][]string)

	for _, engine := range engines {
		aiProvider, aiModel, err := splitEngineName(engine)
		if err != nil {
			log.Warnln(err)
		}

		if aiModel == "" {
			var exists bool
			aiModel, exists = defaultProviderModel[aiProvider]
			if !exists {
				return nil, fmt.Errorf("no provider model found for %s", aiProvider)
			}
		}

		engineFunc, exists := engineMap[aiProvider]
		if !exists {
			return nil, fmt.Errorf("no engine found for %s", aiProvider)
		}

		apiKey, exists := apiKeys[aiProvider]
		if !exists {
			return nil, fmt.Errorf("no API key found for %s", aiProvider)
		}

		engineKey := fmt.Sprintf("%s:%s", aiProvider, aiModel)
		log.Infof("Asking %s: %s", engineKey, prompt)
		responses, err := engineFunc(prompt, aiModel, apiKey)
		if err == nil {
			result[engineKey] = responses
		} else {
			log.Errorf("Engine %s: %v", engineKey, err)
		}
	}

	return result, nil
}
