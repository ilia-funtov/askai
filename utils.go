package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getProgramUserDir() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	return filepath.Join(user.HomeDir, "."+programName), nil
}

func readPromptFromStdin(progOptions *ProgramOptions) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	if reader == nil {
		return "", fmt.Errorf("bufio.NewReader failed")
	}

	var stdinPrompt string
	var err error

	isTerminal := isatty.IsTerminal(os.Stdin.Fd())

	if isTerminal {
		if progOptions.cmdPrompt == "" && !progOptions.batchMode {
			fmt.Println("Enter prompt to AI:")

			stdinPrompt, err = reader.ReadString('\n')
			if err != nil {
				return "", fmt.Errorf("failed to read prompt from stdin: %w", err)
			}
		}
	} else {
		stdinPrompt, err = readStreamedPrompt(reader)
		if err != nil {
			return "", err
		}
	}

	return strings.TrimSpace(stdinPrompt), nil
}

func readStreamedPrompt(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
	if scanner == nil {
		return "", fmt.Errorf("bufio.NewScanner failed")
	}

	var stdinPrompt string

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read prompt from stdin: %w", err)
		}

		data := scanner.Bytes()
		if !utf8.Valid(data) {
			return "", fmt.Errorf("input from stdin is not valid utf-8")
		}

		text := string(data)
		if stdinPrompt != "" && text != "" {
			stdinPrompt += "\n"
		}
		stdinPrompt += text
	}

	return stdinPrompt, nil
}

func initLoggingToFile(programName string, configDir string, logFileDir string) *os.File {
	logLevel := log.InfoLevel

	viper := viper.New()

	viper.SetConfigName(programName)
	viper.AddConfigPath(configDir)

	err := viper.ReadInConfig()

	if err == nil {
		dir := viper.GetString("logdir")
		if dir != "" {
			dir, err = filepath.Abs(dir)
			if err == nil {
				logFileDir = dir
			}
		}

		levelStr := viper.GetString("level")
		level, err := log.ParseLevel(levelStr)
		if err == nil {
			logLevel = level
		}

		if viper.GetString("formatter") == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		}
	} else {
		log.Warningf("failed to read log config file: %v", err)
	}

	logFilePath := filepath.Join(logFileDir, defaultLogFileName)

	return initLoggingToFileConfigless(logFilePath, logLevel)
}

func initLoggingToFileConfigless(logFilePath string, level log.Level) *os.File {
	dirPath := filepath.Dir(logFilePath)

	if fileInfo, err := os.Stat(dirPath); os.IsNotExist(err) || !fileInfo.IsDir() {
		const dirPermissionMask = 0770
		err := os.MkdirAll(dirPath, dirPermissionMask)
		if err != nil {
			log.Warningf("failed to create log directory: %v", err)

			return nil
		}
	}

	const logfilePermissionMask = 0640
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, logfilePermissionMask)
	if err != nil {
		log.Warningf("failed to create log file: %v", err)
	}

	if logFile != nil {
		log.SetOutput(logFile)
		log.SetLevel(level)
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

func makeFullPrompt(prompt string, context string) string {
	if prompt != "" && context != "" {
		return prompt + "\n" + context
	} else if prompt != "" {
		return prompt
	} else if context != "" {
		return context
	}

	return ""
}
