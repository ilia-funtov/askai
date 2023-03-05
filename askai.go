package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func run() error {
	var po ProgramOptions
	po.add()
	po.parse()

	userProgramDir, err := getProgramUserDir()
	if err != nil {
		return err
	}

	// log file path if it is not set in config file
	altLogFileDir := filepath.Join(
		userProgramDir,
		defaultLogDir)

	configDir := filepath.Join(
		userProgramDir,
		defaultConfigDir)

	logFile := initLoggingToFile(
		programName,
		configDir,
		altLogFileDir)

	if logFile != nil {
		defer logFile.Close()
	}

	log.Debugf("Program options: %v", po)

	apiKeysConfigFilePath := filepath.Join(
		userProgramDir,
		defaultConfigDir,
		defaultApiKeysConfigFileName)

	apiKeys, err := readApiKeys(apiKeysConfigFilePath, &po)
	if err != nil {
		return err
	}

	var stdinPrompt string
	if !po.noStdin {
		stdinPrompt, err = readPromptFromStdin(&po)
		if err != nil {
			return err
		}
	}

	prompt := makePrompt(&po, stdinPrompt)
	log.Infof("Prompt: %s", prompt)

	if prompt == "" {
		return fmt.Errorf("prompt to AI is empty")
	}

	if po.printPrompt {
		fmt.Printf("Prompt: %s", prompt)
	}

	responseMap, err := askAI(po.engines, prompt, apiKeys)
	if err != nil {
		return fmt.Errorf("failed to ask AI: %v", err)
	}

	for engineKey, responses := range responseMap {
		log.Infof("Engine: %s", engineKey)
		log.Infof("Number of responses: %d", len(responses))
		log.Tracef("Responses: %v", responses)

		if po.printAIEngine {
			fmt.Printf(defaultPrintAIEngineTemplate, engineKey)
		}
		for _, response := range responses {
			fmt.Println(strings.TrimSpace(response))
		}
	}

	return nil
}

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
}

func main() {
	if err := run(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
