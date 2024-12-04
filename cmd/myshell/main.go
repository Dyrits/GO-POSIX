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

	commands["exit"] = func(argument string) {
		if integer, err := strconv.Atoi(argument); err == nil {
			os.Exit(integer)
		} else {
			_, _ = fmt.Fprint(os.Stdout, "Invalid argument for exit\n")
		}
	}

	commands["echo"] = func(argument string) {
		_, _ = fmt.Fprint(os.Stdout, argument+"\n")
	}

	commands["type"] = func(argument string) {
		if _, exists := commands[argument]; exists {
			_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v is a shell builtin\n", argument))
		} else {
			path, err := exec.LookPath(argument)
			if err == nil {
				_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v is %v\n", argument, path))
			} else {
				_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v: not found\n", argument))
			}
		}
	}

	for {
		input := read(reader)
		fields := strings.Fields(input)
		command := fields[0]
		argument := strings.Join(fields[1:], " ")
		if function, exists := commands[command]; exists {
			function(argument)
		} else {
			_, _ = fmt.Fprint(os.Stdout, fmt.Sprintf("%v: command not found\n", input))
		}
	}
}

func read(reader *bufio.Reader) string {
	_, _ = fmt.Fprint(os.Stdout, "$ ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
