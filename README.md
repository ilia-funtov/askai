# Askai
Command line tool to ask AI for help.
## Description
Askai is a tool with the command line interface to ask questions to AI provided by OpenAI and Cohere via REST API. It can read the user's prompt from its command line parameter or/and from stdin. The answer is printed to stdout. The question can be directed to one AI provider (OpenAI or Cohere) or both providers simultaneously. It is capable of processing of text inputs that are larger than maximum number of tokens supported by the given model. This is done by splitting the prompt into multiple text segments and making summary for each of them.

## Build
Prerequisites:
- go 1.20 or later
- make

Just run make command in the root directory of the project.
```
ilia:~/Projects/askai$ make
go build -o bin/askai
```
The binary can be found in bin directory.

## Usage
Getting help.
```
ilia:~/Projects/askai/bin$ ./askai --help
Usage of ./askai:
  -b    Batch mode, do not ask for prompt if stdin is empty
  -e string
        AI engine to use (default "cohere")
  -ea
        Use all supported AI engines
  -nostdin
        Skip reading prompt from stdin
  -p string
        Prompt to AI
  -pe
        Print engine name in output
  -pp
        Print prompt in output
```

Asking a question.
```
ilia:~/Projects/askai/bin$ ./askai "What does Cohere mean?"
Cohere is an adjective that means sticking together firmly as parts.
```

If the prompt is empty, it will be requested from stdin.
```
ilia:~/Projects/askai/bin$ ./askai
Enter prompt to AI:
How are you?
I'm doing well, thank you. How about you?
```

You can apply some prompt to the text redirected to stdin.
```
ilia:~/Projects/askai/bin$ man ls | ./askai "Summary"
ls is a command line tool for displaying information about files and directories in the Linux operating system. It can display information such as file size, modification time, and file type. It can also be used to sort files by size, time, or name.
```

## Configuration
The program's configuraion is stored in user's home directory: ~/.askai/config/askai.json  
Example of configuration:
```json
{
    "apikeys": {
        "cohere": "",
        "openai": ""
    },
    "engine": "cohere",
    "summarizeprompt": "Summarize:",
    "providermodel": {
        "cohere": "command-xlarge-nightly",
        "openai": "gpt-3.5-turbo"
    },
    "printaiengine": "#%s#",
    "loglevel": "trace",
    "logdir": "~/.askai/log",
    "logformat": ""
}
```

- section "apikeys" contains API keys for Cohere and OpenAI. You can fill this information in configuration file or it will be asked on the first run.
- parameter "engine" is used to specify the default engine to use (openai or cohere).
- parameter "summarizeprompt" is used to specify the prompt to summarize the text input.
- section "providermodel" is used to specify the default provider model to use for each AI provider.
- parameter "printaiengine" is used to specify print template to print AI engine name in output.
- parameter "loglevel" is used to specify the default log level. It can be trace, debug, info, warn, error, fatal.
- parameter "logdir" is used to specify the default log directory.
- parameter "logformat" is used to specify the default log format.

## License
The project is distributed under the terms of the MIT license.

## Third party
- [cohere-go](https://github.com/cohere-ai/cohere-go)
- [go-isatty](https://github.com/mattn/go-isatty)
- [tiktoken-go](https://github.com/pkoukk/tiktoken-go)
- [go-gpt3](https://github.com/sashabaranov/go-gpt3)
- [logrus](https://github.com/sirupsen/logrus)
- [testify](https://github.com/stretchr/testify)