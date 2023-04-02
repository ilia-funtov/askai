package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func readOrAskAPIKeys(path string, progOptions *ProgramOptions) (map[string]string, error) {
	config := viper.New()

	dir, filename := filepath.Split(path)
	config.SetConfigName(filename)
	config.AddConfigPath(dir)

	apiKeys := make(map[string]string)

	err := config.ReadInConfig()

	if err != nil {
		errtype := reflect.TypeOf(err)
		if !progOptions.batchMode && errtype == reflect.TypeOf(viper.ConfigFileNotFoundError{}) {
			apiKeys, err = askAndStoreAPIKeys(progOptions.engines, path, config)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("failed to open config file with API keys %s: %w", path, err)
		}
	} else {
		allKeys := config.AllKeys()
		for _, key := range allKeys {
			apiKeys[key] = config.GetString(key)
		}

		missedKeys := make([]string, 0, len(progOptions.engines))
		for _, key := range progOptions.engines {
			if _, exists := apiKeys[key]; !exists {
				missedKeys = append(missedKeys, key)
			}
		}

		newAPIKeys, err := askAndStoreAPIKeys(missedKeys, path, config)
		if err == nil {
			for key, value := range newAPIKeys {
				apiKeys[key] = value
			}
		}
	}

	return apiKeys, nil
}

func askAndStoreAPIKeys(engines []string, path string, config *viper.Viper) (map[string]string, error) {
	apiKeys, err := askAPIKeys(engines)
	if err != nil {
		return nil, err
	}

	err = storeAPIKeys(apiKeys, config, path)
	if err != nil {
		log.Warnf("failed to write API keys to %s: %v", path, err)
	}

	return apiKeys, nil
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

func storeAPIKeys(apiKeys map[string]string, config *viper.Viper, keysConfigPath string) error {
	for provider, apiKey := range apiKeys {
		config.Set(provider, apiKey)
	}

	ext := filepath.Ext(keysConfigPath)
	if ext == "" {
		if !strings.HasSuffix(keysConfigPath, ".") {
			keysConfigPath += "."
		}
		keysConfigPath += defaultAPIKeysConfigExtension
	}

	err := config.WriteConfigAs(keysConfigPath)
	if err != nil {
		return fmt.Errorf("failed to write API keys to %s: %w", keysConfigPath, err)
	}

	return nil
}
