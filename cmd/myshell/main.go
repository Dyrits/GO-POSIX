package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	reader := bufio.NewReader(os.Stdin)

	commands := map[string]func(){
		"exit 0": func() {
			os.Exit(0)
		},
	}

	for {
		input := read(reader)
		if function, exists := commands[input]; exists {
			function()
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
