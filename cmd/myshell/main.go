package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type CommandFunction func(argument string)

func main() {
	reader := bufio.NewReader(os.Stdin)

	commands := map[string]CommandFunction{
		"exit": func(argument string) {
			if integer, err := strconv.Atoi(argument); err == nil {
				os.Exit(integer)
			} else {
				fmt.Fprint(os.Stdout, "Invalid argument for exit\n")
			}
		},
		"echo": func(argument string) {
			fmt.Fprint(os.Stdout, argument+"\n")
		},
	}

	for {
		input := read(reader)
		// The shell only works for commands with 4 characters.
		command := input[0:4]
		argument := input[5:]
		if function, exists := commands[command]; exists {
			function(argument)
		} else {
			fmt.Fprint(os.Stdout, fmt.Sprintf("%v: command not found\n", input))
		}
	}
}

func read(reader *bufio.Reader) string {
	fmt.Fprint(os.Stdout, "$ ")
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
