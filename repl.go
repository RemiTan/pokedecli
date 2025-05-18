package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/RemiTan/pokedexcli/internal/pokeapi"
)

type config struct {
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	Pokemons         map[string]pokeapi.Pokemon
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	supportedCommands := getCommands()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		cliCommand, ok := supportedCommands[words[0]]

		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		args := []string{}
		if len(words) > 1 {
			args = words[1:]
		}

		if err := cliCommand.callback(cfg, args...); err != nil {
			fmt.Println(err)
		}
	}
}
