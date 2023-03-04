package main

import (
	"flag"
	"strings"
)

type ProgramOptions struct {
	cmdPrompt     string
	aiEngineList  string
	engines       []string
	batchMode     bool
	printAIEngine bool
	printPrompt   bool
}

func (po *ProgramOptions) add() {
	flag.StringVar(&po.cmdPrompt, "prompt", "", "Prompt to AI")
	flag.BoolVar(&po.batchMode, "batch", false, "Batch mode, do not ask for prompt if stdin is empty")
	flag.StringVar(&po.aiEngineList, "engine", defaultEngine, "AI engine to use")
	flag.BoolVar(&po.printAIEngine, "printengine", false, "Print engine name in output")
	flag.BoolVar(&po.printPrompt, "printprompt", false, "Print prompt in output")
}

func (po *ProgramOptions) parse() {
	flag.Parse()

	po.aiEngineList = strings.ToLower(po.aiEngineList)

	{
		engines := strings.Split(po.aiEngineList, ",")
		engineMap := make(map[string]bool)
		for _, engine := range engines {
			engine = strings.TrimSpace(engine)
			if engine != "" {
				engineMap[engine] = true
			}
		}

		po.engines = make([]string, 0, len(engineMap))
		for k := range engineMap {
			po.engines = append(po.engines, k)
		}
	}

	if po.cmdPrompt == "" && flag.NArg() >= 1 {
		// try to take first argument of command as prompt
		po.cmdPrompt = flag.Arg(0)
	}

	po.cmdPrompt = strings.TrimSpace(po.cmdPrompt)
}