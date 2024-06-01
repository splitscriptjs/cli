package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
)

type Handle func(args []string)
type Handler struct {
	arguments   Args
	function    Handle
	description string
}
type Args map[string]string

var commands map[string]Handler = make(map[string]Handler)

func Add(name string, description string, args Args, function Handle) {
	commands[name] = Handler{
		arguments:   args,
		function:    function,
		description: description,
	}
}
func Parse() {
	name := "help"
	args := []string{}
	if len(os.Args) > 1 {
		name = os.Args[1]
		args = os.Args[2:]
	}

	handler := commands[name]
	if handler.arguments == nil || handler.function == nil {
		fmt.Println("Command not found")
		return
	}
	handler.function(args)
}
func List() []table.Row {
	rows := make([]table.Row, len(commands))

	i := 0
	for cmd, val := range commands {
		rows[i] = []string{cmd, val.description}
		i++
	}
	return rows
}
