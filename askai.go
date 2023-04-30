package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func run() error {
	programConfig, err := initProgramConfig()
	if err != nil {
		return fmt.Errorf("failed to init program configuration: %w", err)
	}

	var progOptions ProgramOptions
	progOptions.add(programConfig.Engine)
	progOptions.parse()

	log.Debugf("Program options: %v", progOptions)

	err = initAPIKeysConfig(progOptions, programConfig)
	if err != nil {
		return fmt.Errorf("failed to init API keys configuration: %w", err)
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

	responseMap, err := askAI(progOptions.engines, message, *programConfig)
	if err != nil {
		return fmt.Errorf("failed to ask AI: %w", err)
	}

	printResponses(responseMap, progOptions, *programConfig)

	return nil
}

func printResponses(responseMap map[string][]string, progOptions ProgramOptions, progConfig ProgramConfig) {
	for engineKey, responses := range responseMap {
		log.Infof("Engine: %s", engineKey)
		log.Infof("Number of responses: %d", len(responses))
		log.Tracef("Responses: %v", responses)

		if progOptions.printAIEngine {
			fmt.Println(fmt.Sprintf(progConfig.PrintAIEngineTemplate, engineKey))
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
