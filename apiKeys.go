package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func readApiKeys(path string, po *ProgramOptions) (map[string]string, error) {
	apiConfigFile, err := os.Open(path)
	if apiConfigFile != nil {
		defer apiConfigFile.Close()
	}

	apiKeys := make(map[string]string)

	if err != nil {
		if !po.batchMode && errors.Is(err, os.ErrNotExist) {
			// Create config file for api keys. Iterate over all providers and ask for each one.

			reader := bufio.NewReader(os.Stdin)
			if reader == nil {
				return nil, fmt.Errorf("bufio.NewReader failed")
			}

			for _, engine := range po.engines {
				aiProvider, _, err := splitEngineName(engine)
				if err != nil {
					return nil, err
				}

				fmt.Printf("Enter API key for %s:\n", aiProvider)
				apiKey, err := reader.ReadString('\n')
				if err != nil {
					return nil, fmt.Errorf("failed to read API key from stdin: %v", err)
				}

				apiKey = strings.TrimSpace(apiKey)
				if apiKey == "" {
					continue
				}

				apiKeys[aiProvider] = apiKey
			}

			if len(apiKeys) != 0 {
				jsonContent, err := json.Marshal(apiKeys)
				if err != nil {
					return nil, fmt.Errorf("failed to serialize API keys to JSON: %v", err)
				}

				apiConfigFile, err = os.Create(path)
				if err != nil {
					return nil, fmt.Errorf("failed to create API keys config file: %v", err)
				}
				defer apiConfigFile.Close()

				_, err = apiConfigFile.Write(jsonContent)
				if err != nil {
					return nil, fmt.Errorf("failed to write API keys to config file: %v", err)
				}
			}
		} else {
			return nil, fmt.Errorf("failed to open config file with API keys %s: %v", path, err)
		}
	} else {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("failed to stat API keys config file: %v", err)
		}

		if fileInfo.Size() == 0 {
			return nil, fmt.Errorf("API keys config file is empty")
		}

		buffer := make([]byte, fileInfo.Size())
		_, err = apiConfigFile.Read(buffer)
		if err != nil {
			return nil, fmt.Errorf("failed to read API keys from config file: %v", err)
		}

		err = json.Unmarshal(buffer, &apiKeys)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize API keys from JSON: %v", err)
		}
	}

	return apiKeys, nil
}
