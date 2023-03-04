package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
)

func getProgramUserDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}

	return filepath.Join(user.HomeDir, "."+programName), nil
}

func readPromptFromStdin(po *ProgramOptions) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	if reader == nil {
		return "", fmt.Errorf("bufio.NewReader failed")
	}

	var stdinPrompt string

	isTerminal := isatty.IsTerminal(os.Stdin.Fd())

	if isTerminal {
		if po.cmdPrompt == "" && !po.batchMode {
			fmt.Println("Enter prompt to AI:")

			var err error
			stdinPrompt, err = reader.ReadString('\n')
			if err != nil {
				return "", fmt.Errorf("failed to read prompt from stdin: %v", err)
			}
		}
	} else {
		scanner := bufio.NewScanner(reader)
		if scanner == nil {
			return "", fmt.Errorf("bufio.NewScanner failed")
		}

		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return "", fmt.Errorf("failed to read prompt from stdin: %v", err)
			}

			text := scanner.Text()
			if stdinPrompt != "" && text != "" {
				stdinPrompt += "\n"
			}
			stdinPrompt += text
		}
	}

	return strings.TrimSpace(stdinPrompt), nil
}

func makePrompt(po *ProgramOptions, stdinPrompt string) string {
	if po.cmdPrompt != "" && stdinPrompt != "" {
		return po.cmdPrompt + "\n" + stdinPrompt
	} else if po.cmdPrompt != "" {
		return po.cmdPrompt
	} else if stdinPrompt != "" {
		return stdinPrompt
	}

	return ""
}

func initLoggingToFile(logFilePath string) *os.File {
	dirPath := filepath.Dir(logFilePath)

	if fileInfo, err := os.Stat(dirPath); os.IsNotExist(err) || !fileInfo.IsDir() {
		err := os.MkdirAll(dirPath, 0770)
		if err != nil {
			log.Errorf("failed to create log directory: %v", err)
			return nil
		}
	}

	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Warningf("failed to create log file: %v\n", err)
	}

	if logFile != nil {
		runtime.SetFinalizer(
			logFile,
			func(logFile *os.File) {
				log.SetOutput(os.Stderr)
			})

		log.SetOutput(logFile)
		log.SetLevel(log.InfoLevel)
	}

	return logFile
}

func splitEngineName(engineName string) (string, string, error) {
	parts := strings.Split(engineName, ":")
	if len(parts) == 0 {
		return "", "", fmt.Errorf("failed to split engine name: %s", engineName)
	}

	aiProvider := strings.TrimSpace(parts[0])
	aiModel := ""
	if len(parts) > 1 {
		aiModel = strings.TrimSpace(parts[1])
	}

	return aiProvider, aiModel, nil
}
