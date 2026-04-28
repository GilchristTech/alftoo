package main

import (
	"fmt"
	"strings"
)

type ColonCommand interface {
	// The ColonCommand interface used by Alftoo to register
	// behaviors for commands which start with a ':', such as
	// :font. By making these an interface, the underlying data
	// type can be anything, which lets commands which control a
	// single value to be stored as primative data types, such as a
	// string. This design allows each command to take up only the
	// amount of memory they require.

	Run(name, args string) error
}

var colon_commands map[string]ColonCommand

// Recreates an empty registery of colon commands
func ClearColonCommands() {
	colon_commands = make(map[string]ColonCommand)
}

func AddColonCommand(name string, cc ColonCommand) {
	var err_h string = "alftoo.AddColonCommand"

	if colon_commands == nil {
		ClearColonCommands()
	}

	if cc == nil {
		panic(err_h + ": command is nil")

	} else if !strings.HasPrefix(name, ":") {
		panic(fmt.Sprintf(
			`%s: command name "%s" does not start with a colon`,
			err_h, name,
		))

	} else if _, exists := colon_commands[name]; exists {
		panic(fmt.Sprintf(
			`%s: command with name "%s" already exists`,
			err_h, name,
		))

	} else {
		colon_commands[name] = cc
	}
}
