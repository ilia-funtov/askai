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

func readApiKeys(path string, po *ProgramOptions) (map[string]string, error) {
	v := viper.New()

	dir, filename := filepath.Split(path)
	v.SetConfigName(filename)
	v.AddConfigPath(dir)

	apiKeys := make(map[string]string)

	err := v.ReadInConfig()

	if err != nil {
		errtype := reflect.TypeOf(err)
		if !po.batchMode && errtype == reflect.TypeOf(viper.ConfigFileNotFoundError{}) {
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
					return nil, fmt.Errorf("failed to read API key from stdin: %w", err)
				}

				apiKey = strings.TrimSpace(apiKey)
				if apiKey == "" {
					continue
				}

				apiKeys[aiProvider] = apiKey
			}

			if len(apiKeys) != 0 {
				for provider, apiKey := range apiKeys {
					v.Set(provider, apiKey)
				}

				ext := filepath.Ext(path)
				if ext == "" {
					if !strings.HasSuffix(path, ".") {
						path += "."
					}
					path += defaultAPIKeysConfigExtension
				}

				err = v.WriteConfigAs(path)
				if err != nil {
					log.Warnf("failed to write API keys to %s: %v", path, err)
				}
			}
		} else {
			return nil, fmt.Errorf("failed to open config file with API keys %s: %w", path, err)
		}
	} else {
		allKeys := v.AllKeys()
		for _, key := range allKeys {
			apiKeys[key] = v.GetString(key)
		}
	}

	return apiKeys, nil
}
