package main

const programName = "askai"

const defaultConfigDir = "config"
const defaultLogDir = "log"

const defaultConfigFileExtension = "json"
const defaultLogFileName = programName + ".log"
const defaultPrintAIEngineTemplate = "#%s#"
const defaultEngine = "cohere"
const defaultSummarizePrompt = "Summarize:"

var defaultProviderModel = map[string]string{
	"openai": "gpt-3.5-turbo",
	"cohere": "command-xlarge-nightly",
}
