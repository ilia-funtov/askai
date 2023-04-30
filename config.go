package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"

	log "github.com/sirupsen/logrus"
)

type ProgramConfig struct {
	APIKeys               map[string]string `json:"apikeys"`
	Engine                string            `json:"engine"`
	SummarizePrompt       string            `json:"summarizeprompt"`
	ProviderModel         map[string]string `json:"providermodel"`
	PrintAIEngineTemplate string            `json:"printaiengine"`
	LogLevel              string            `json:"loglevel"`
	LogDir                string            `json:"logdir"`
	LogFormatter          string            `json:"logformat"`
	configFilePath        string            // don't serialize this
}

func initProgramConfig() (*ProgramConfig, error) {
	userProgramDir, err := getProgramUserDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(
		userProgramDir,
		defaultConfigDir)

	var config ProgramConfig

	config.configFilePath = filepath.Join(
		configDir,
		programName+"."+defaultConfigFileExtension)

	config.Engine = defaultEngine
	config.SummarizePrompt = defaultSummarizePrompt
	config.ProviderModel = defaultProviderModel
	config.PrintAIEngineTemplate = defaultPrintAIEngineTemplate

	data, err := os.ReadFile(config.configFilePath)
	if err == nil {
		err = json.Unmarshal(data, &config)

		if err != nil {
			log.Warningf("failed to deserialize config file: %v", err)
		}
	} else {
		log.Warningf("failed to read config file: %v", err)
	}

	if config.LogDir == "" {
		config.LogDir = filepath.Join(
			userProgramDir,
			defaultLogDir)
	}

	initLoggingToFile(config)

	return &config, nil
}

func initAPIKeysConfig(progOptions ProgramOptions, config *ProgramConfig) error {
	newAPIKeys, err := processMissedAPIKeys(config.APIKeys, progOptions.engines)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(config.APIKeys, newAPIKeys) {
		config.APIKeys = newAPIKeys

		data, err := json.MarshalIndent(config, "", " ")
		if err == nil {
			err = os.WriteFile(config.configFilePath, data, 0600)
		}

		if err != nil {
			log.Warningf("failed to write to config file: %v", err)
		}
	}

	return nil
}
