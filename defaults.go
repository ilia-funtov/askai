package main

const programName = "askai"

const defaultConfigDir = "config"
const defaultLogDir = "log"
const defaultApiKeysConfigFileName = "apikeys"
const defaultApiKeysConfigExtension = "json"
const defaultLogFileName = programName + ".log"
const defaultPrintAIEngineTemplate = "#%s#\n"
const defaultEngine = "cohere:command-xlarge-nightly"

var defaultProviderModel = map[string]string{
	"openai": "gpt-3.5-turbo",
	"cohere": "command-xlarge-nightly",
}
