package main

import (
	"flag"
	"strings"

	"golang.org/x/exp/maps"
)

type ProgramOptions struct {
	cmdPrompt     string
	aiEngineList  string
	engines       []string
	allEngines    bool
	batchMode     bool
	printAIEngine bool
	printPrompt   bool
	noStdin       bool
}

func (po *ProgramOptions) add() {
	flag.StringVar(&po.cmdPrompt, "p", "", "Prompt to AI")
	flag.BoolVar(&po.batchMode, "b", false, "Batch mode, do not ask for prompt if stdin is empty")
	flag.StringVar(&po.aiEngineList, "e", defaultEngine, "AI engine to use")
	flag.BoolVar(&po.allEngines, "ea", false, "Use all supported AI engines")
	flag.BoolVar(&po.printAIEngine, "pe", false, "Print engine name in output")
	flag.BoolVar(&po.printPrompt, "pp", false, "Print prompt in output")
	flag.BoolVar(&po.noStdin, "nostdin", false, "Skip reading prompt from stdin")
}

func (po *ProgramOptions) parse() {
	flag.Parse()

	po.aiEngineList = strings.ToLower(po.aiEngineList)

	if po.allEngines {
		po.engines = maps.Keys(engineMap)
	} else {
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
