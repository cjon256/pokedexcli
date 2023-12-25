package main

import (
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	var commands map[string]cliCommand

	commandHelp := func() error {
        fmt.Println("\nWelcome to the Pokedex!")
		fmt.Println("Usage:\n")
		for _, command := range commands {
			fmt.Printf("\t%s: %s\n", command.name, command.description)
		}
        fmt.Println()
		return nil
	}

	commandExit := func() error {
		os.Exit(0)
		return nil // This is required to make the compiler happy
	}

	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}

	prompt := "pokedex> "
	for {
		fmt.Print(prompt)
		var response string
		fmt.Scanln(&response)
		command, ok := commands[response]
		if !ok {
			fmt.Printf("Unknown command: %s\n", response)
			continue
		}
		action := command.callback
		err := action()
		if err != nil {
            // this isn't used, but the interface seems to imply we will need it
			fmt.Printf("Error: %s\n", err)
		}
	}
}
