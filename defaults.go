package main

const programName = "askai"

const defaultConfigDir = "config"
const defaultLogDir = "log"
const defaultAPIKeysConfigFileName = "apikeys"
const defaultAPIKeysConfigExtension = "json"
const defaultLogFileName = programName + ".log"
const defaultPrintAIEngineTemplate = "#%s#\n"
const defaultEngine = "cohere:command-xlarge-nightly"
const defaultTLDRPrompt = "TL;DR"

var defaultProviderModel = map[string]string{
	"openai": "gpt-3.5-turbo",
	"cohere": "command-xlarge-nightly",
}
