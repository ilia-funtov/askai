package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type ProgramConfig struct {
	apiKeys map[string]string
}

func initProgramConfig(progOptions ProgramOptions) (*ProgramConfig, error) {
	userProgramDir, err := getProgramUserDir()
	if err != nil {
		return nil, err
	}

	// log file path if it is not set in config file
	altLogFileDir := filepath.Join(
		userProgramDir,
		defaultLogDir)

	configDir := filepath.Join(
		userProgramDir,
		defaultConfigDir)

	initLoggingToFile(
		programName,
		configDir,
		altLogFileDir)

	apiKeysConfigFilePath := filepath.Join(
		userProgramDir,
		defaultConfigDir,
		defaultAPIKeysConfigFileName)

	apiKeys, err := readOrAskAPIKeys(apiKeysConfigFilePath, &progOptions)
	if err != nil {
		return nil, err
	}

	return &ProgramConfig{apiKeys: apiKeys}, nil
}

func run() error {
	var progOptions ProgramOptions
	progOptions.add()
	progOptions.parse()

	log.Debugf("Program options: %v", progOptions)

	programConfig, err := initProgramConfig(progOptions)
	if err != nil {
		return err
	}

	var stdinPrompt string
	if !progOptions.noStdin {
		stdinPrompt, err = readPromptFromStdin(&progOptions)
		if err != nil {
			return err
		}
	}

	message := UserMessage{Prompt: progOptions.cmdPrompt, Context: stdinPrompt}
	prompt := message.GetFullPrompt()

	log.Infof("Prompt: %s", prompt)

	if prompt == "" {
		return fmt.Errorf("prompt to AI is empty")
	}

	if progOptions.printPrompt {
		fmt.Printf("Prompt: %s", prompt)
	}

	responseMap, err := askAI(progOptions.engines, message, programConfig.apiKeys)
	if err != nil {
		return fmt.Errorf("failed to ask AI: %w", err)
	}

	printResponses(responseMap, progOptions)

	return nil
}

func printResponses(responseMap map[string][]string, progOptions ProgramOptions) {
	for engineKey, responses := range responseMap {
		log.Infof("Engine: %s", engineKey)
		log.Infof("Number of responses: %d", len(responses))
		log.Tracef("Responses: %v", responses)

		if progOptions.printAIEngine {
			fmt.Printf(defaultPrintAIEngineTemplate, engineKey)
		}
		for _, response := range responses {
			fmt.Println(strings.TrimSpace(response))
		}
	}
}

func main() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)

	if err := run(); err != nil {
		log.Errorln(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
