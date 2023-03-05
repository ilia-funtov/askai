package main

const programName = "askai"

const defaultConfigDir = "config"
const defaultLogDir = "log"
const defaultApiKeysConfigFileName = "apikeys.conf"
const defaultLogFileName = programName + ".log"
const defaultPrintAIEngineTemplate = "#%s#\n"
const defaultEngine = "openai:gpt-3.5-turbo"

var defaultProviderModel = map[string]string{
	"openai": "gpt-3.5-turbo",
	"cohere": "command-xlarge-nightly",
}
