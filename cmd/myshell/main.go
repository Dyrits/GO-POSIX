package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type CommandFunction func(argument string, writer io.Writer)

func main() {
	reader := bufio.NewReader(os.Stdin)

	commands := make(map[string]CommandFunction)

	commands["exit"] = func(arguments string, writer io.Writer) {
		if integer, err := strconv.Atoi(arguments); err == nil {
			os.Exit(integer)
		} else {
			_, _ = fmt.Fprint(writer, "Invalid argument for exit\n")
		}
	}

	commands["echo"] = func(arguments string, writer io.Writer) {
		_, _ = fmt.Fprint(writer, arguments+"\n")
	}

	commands["pwd"] = func(arguments string, writer io.Writer) {
		path, err := os.Getwd()
		if err == nil {
			_, _ = fmt.Fprint(writer, fmt.Sprintf("%v\n", path))
		}
	}

	commands["cd"] = func(arguments string, writer io.Writer) {
		first := arguments[0]
		if first == '~' {
			arguments = os.Getenv("HOME") + arguments[1:]
		}
		err := os.Chdir(arguments)
		if err != nil {
			_, _ = fmt.Fprint(writer, fmt.Sprintf("%v: No such file or directory\n", arguments))
		}
	}

	commands["type"] = func(arguments string, writer io.Writer) {
		if _, exists := commands[arguments]; exists {
			_, _ = fmt.Fprint(writer, fmt.Sprintf("%v is a shell builtin\n", arguments))
		} else {
			path, err := exec.LookPath(arguments)
			if err == nil {
				_, _ = fmt.Fprint(writer, fmt.Sprintf("%v is %v\n", arguments, path))
			} else {
				_, _ = fmt.Fprint(writer, fmt.Sprintf("%v: not found\n", arguments))
			}
		}
	}

	for {
		input := read(reader)
		fields := split(input)
		if len(fields) == 0 {
			continue
		}
		command := fields[0]
		arguments := fields[1:]

		var output string
		for index, argument := range arguments {
			if argument == ">" || argument == "1>" {
				if index+1 < len(arguments) {
					output = arguments[index+1]
					arguments = arguments[:index]
					break
				}
			}
		}

		var writer io.Writer = os.Stdout
		if output != "" {
			file, err := os.Create(output)
			if err != nil {
				_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("Error creating file: %v\n", err))
				continue
			}
			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("Error closing file: %v\n", err))
				}
			}(file)
			writer = file
		}

		if function, exists := commands[command]; exists {
			arguments := strings.Join(arguments, " ")
			// Run a registered command.
			function(arguments, writer)
		} else {
			// Run a system command (usually in the PATH).

			// Remove outer quotes if they exist.
			if (strings.HasPrefix(command, "'") && strings.HasSuffix(command, "'")) ||
				(strings.HasPrefix(command, "\"") && strings.HasSuffix(command, "\"")) {
				command = command[1 : len(command)-1]
			}

			command = strings.ReplaceAll(command, "'", `\'`)

			// Run the command
			run := exec.Command(command, arguments...)
			run.Stdout = writer
			run.Stderr = os.Stderr
			_ = run.Run()
		}
	}
}

func split(input string) []string {
	var result []string
	var current []rune
	single := false
	double := false
	escape := false
	path := false

	if strings.HasPrefix(input, "cat ") || strings.HasPrefix(input, "ls ") {
		path = true
	}

	escapable := func(character rune) bool {
		return character == '"' || character == '\'' || character == '\\' || character == ' ' || character == 'n'
	}

	for _, character := range input {
		if escape {
			// If the character is an escapable character, and not in a single-quoted string, or used as a path, append it to the current argument.
			if escapable(character) && !single && !path {
				current = append(current, character)
			} else {
				// If not, append the escape character and the current character.
				current = append(current, '\\', character)
			}
			escape = false
		} else if character == '\\' {
			escape = true
		} else if character == '\'' && !double {
			single = !single
		} else if character == '"' && !single {
			double = !double
		} else if character == ' ' && !single && !double {
			if len(current) > 0 {
				result = append(result, string(current))
				current = nil
			}
		} else {
			current = append(current, character)
		}
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}

func read(reader *bufio.Reader) string {
	_, _ = fmt.Fprint(os.Stdout, "$ ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
