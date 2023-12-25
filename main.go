package main

import (
	"fmt"
	"os"
    "log"
    "net/http"
    "io"
    "encoding/json"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type pokedexLocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
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

    var prevLocation string
    var nextLocation string

    commandMap := func() error {
        if nextLocation == "" {
            nextLocation = "https://pokeapi.co/api/v2/location-area/?limit=20"
        }
        if nextLocation == "null" {
            fmt.Println("You have reached the end of the map!")
            return nil
        }
        res, err := http.Get(nextLocation)
        if err != nil {
            log.Fatal(err)
        }
        body, err := io.ReadAll(res.Body)
        res.Body.Close()
        if res.StatusCode > 299 {
            log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
        }
        if err != nil {
            log.Fatal(err)
        }
        var currentLocationArea pokedexLocationArea
        json.Unmarshal(body, &currentLocationArea)
        nextLocation = currentLocationArea.Next
        prevLocation = currentLocationArea.Previous
        for _, location := range currentLocationArea.Results {
            fmt.Printf("%s\n", location.Name)
        }
        return nil
    }

    commandMapB := func() error {
        if prevLocation == "" || prevLocation == "null" {
            fmt.Println("You are at the beginning of the map!")
            return nil
        }
        res, err := http.Get(prevLocation)
        if err != nil {
            log.Fatal(err)
        }
        body, err := io.ReadAll(res.Body)
        res.Body.Close()
        if res.StatusCode > 299 {
            log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
        }
        if err != nil {
            log.Fatal(err)
        }
        var currentLocationArea pokedexLocationArea
        json.Unmarshal(body, &currentLocationArea)
        nextLocation = currentLocationArea.Next
        prevLocation = currentLocationArea.Previous
        for _, location := range currentLocationArea.Results {
            fmt.Printf("%s\n", location.Name)
        }
        return nil
    }

	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"q": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
        "map": {
            name:        "map",
            description: "Displays names of the next 20 location areas",
            callback:    commandMap,
        },
        "mapb": {
            name:        "mapb",
            description: "Displays names of the 20 previous location areas",
            callback:    commandMapB,
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
