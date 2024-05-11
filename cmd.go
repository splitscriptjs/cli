package main

import (
	"fmt"
	"os"
)

type Handle func(args []string)
type Handler struct {
	arguments args
	function  Handle
}
type args map[string]string

var commands map[string]Handler = make(map[string]Handler)

func addCmd(name string, args args, function Handle) {
	commands[name] = Handler{
		arguments: args,
		function:  function,
	}
}
func parseCmd() {
	name := os.Args[1]
	handler := commands[name]
	if handler.arguments == nil || handler.function == nil {
		fmt.Println("Command not found")
		return
	}
	handler.function(os.Args[2:])
}
