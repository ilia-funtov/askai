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

	logFilePath := filepath.Join(
		userProgramDir,
		defaultLogDir,
		defaultLogFileName)

	logFile := initLoggingToFile(logFilePath)
	if logFile != nil {
		defer logFile.Close()
	}

	log.Tracef("Program options: %v\n", po)

	apiKeysConfigFilePath := filepath.Join(
		userProgramDir,
		defaultConfigDir,
		defaultApiKeysConfigFileName)

	apiKeys, err := readApiKeys(apiKeysConfigFilePath, &po)
	if err != nil {
		return err
	}

	stdinPrompt, err := readPromptFromStdin(&po)
	if err != nil {
		return err
	}

	prompt := makePrompt(&po, stdinPrompt)
	log.Tracef("Prompt: %s\n", prompt)

	if prompt == "" {
		return fmt.Errorf("prompt to AI is empty")
	}

	if po.printPrompt {
		fmt.Printf("Prompt: %s\n", prompt)
	}

	responseMap, err := askAI(po.engines, prompt, apiKeys)
	if err != nil {
		return fmt.Errorf("failed to ask AI: %v", err)
	}

	for engineKey, responses := range responseMap {
		log.Tracef("Engine: %s\n", engineKey)
		log.Tracef("Responses: %v\n", responses)

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
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	log.SetLevel(log.WarnLevel)
}

func main() {
	if err := run(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}
