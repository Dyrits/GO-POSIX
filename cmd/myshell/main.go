package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type CommandFunction func(argument string)

func main() {
	reader := bufio.NewReader(os.Stdin)

	commands := make(map[string]CommandFunction)

	commands["exit"] = func(arguments string) {
		if integer, err := strconv.Atoi(arguments); err == nil {
			os.Exit(integer)
		} else {
			_, _ = fmt.Fprint(os.Stdout, "Invalid argument for exit\n")
		}
	}

	commands["echo"] = func(arguments string) {
		_, _ = fmt.Fprint(os.Stdout, arguments+"\n")
	}

	commands["pwd"] = func(arguments string) {
		path, err := os.Getwd()
		if err == nil {
			_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v\n", path))
		}
	}

	commands["cd"] = func(arguments string) {
		first := arguments[0]
		if first == '~' {
			arguments = os.Getenv("HOME") + arguments[1:]
		}
		err := os.Chdir(arguments)
		if err != nil {
			_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v: No such file or directory\n", arguments))
		}
	}

	commands["type"] = func(arguments string) {
		if _, exists := commands[arguments]; exists {
			_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v is a shell builtin\n", arguments))
		} else {
			path, err := exec.LookPath(arguments)
			if err == nil {
				_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v is %v\n", arguments, path))
			} else {
				_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v: not found\n", arguments))
			}
		}
	}

	for {
		input := read(reader)
		fields := strings.Fields(input)
		command := fields[0]
		arguments := fields[1:]
		if function, exists := commands[command]; exists {
			// Run a registered command.
			function(strings.Join(arguments, " "))
		} else {
			// Run a system command (usually in the PATH).
			run := exec.Command(command, arguments...)
			run.Stdout = os.Stdout
			run.Stderr = os.Stderr
			err := run.Run()
			if err != nil {
				_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v: command not found\n", command))
			}
		}
	}
}

func read(reader *bufio.Reader) string {
	_, _ = fmt.Fprint(os.Stdout, "$ ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
