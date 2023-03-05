package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
		err := os.MkdirAll(dirPath, 0770)
		if err != nil {
			log.Warningf("failed to create log directory: %v", err)
			return nil
		}
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
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
