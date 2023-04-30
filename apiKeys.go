package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func processMissedAPIKeys(apiKeys map[string]string, engines []string) (map[string]string, error) {
	missedKeys := make([]string, 0, len(engines))
	for _, key := range engines {
		if _, exists := apiKeys[key]; !exists {
			missedKeys = append(missedKeys, key)
		}
	}

	if len(missedKeys) == 0 {
		return apiKeys, nil
	}

	newAPIKeys, err := askAPIKeys(missedKeys)
	if err != nil {
		return nil, err
	}

	for key, value := range apiKeys {
		newAPIKeys[key] = value
	}

	return newAPIKeys, nil
}

func askAPIKeys(engines []string) (map[string]string, error) {
	reader := bufio.NewReader(os.Stdin)
	if reader == nil {
		return nil, fmt.Errorf("bufio.NewReader failed")
	}

	apiKeys := make(map[string]string)

	for _, engine := range engines {
		aiProvider, _, err := splitEngineName(engine)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Enter API key for %s:\n", aiProvider)
		apiKey, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read API key from stdin: %w", err)
		}

		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" {
			continue
		}

		apiKeys[aiProvider] = apiKey
	}

	if len(apiKeys) == 0 {
		return nil, fmt.Errorf("no API keys provided")
	}

	return apiKeys, nil
}
