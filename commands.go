package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args ...string) error
}

func commandExit(c *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range getCommands() {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMap(cfg *config, args ...string) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}
	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Printf("%s\n", loc.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("You are on the first page")
	}

	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}
	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	for _, loc := range locationsResp.Results {
		fmt.Printf("%s\n", loc.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("You need to provide a city name to explore")
	}

	areaName := args[0]

	fmt.Printf("Exploring %s...\n", areaName)
	pokemonsResp, err := cfg.pokeapiClient.ListPokemons(areaName)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemonE := range pokemonsResp.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemonE.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("You need to provide a Pokemon name to capture")
	}

	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	pokemon, err := cfg.pokeapiClient.DetailPokemon(pokemonName)
	if err != nil {
		return err
	}

	baseExperience := pokemon.BaseExperience

	if rand.Intn(baseExperience) < 50 {
		fmt.Printf("%s was caught!\n", pokemonName)
		cfg.Pokemons[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return errors.New("You need to provide a Pokemon name to inspect")
	}
	pokemonName := args[0]

	pokemon, ok := cfg.Pokemons[pokemonName]
	if !ok {
		return errors.New("You have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("   -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, type_ := range pokemon.Types {
		fmt.Printf("   -%s\n", type_.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	if len(cfg.Pokemons) == 0 {
		return errors.New("Your pokedex is empty.")
	}

	fmt.Println("Your pokedex:")
	for key := range cfg.Pokemons {
		fmt.Printf(" - %s\n", key)
	}

	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
		"map": {
			name:        "map",
			description: "Get the next page of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Get the previous page of locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Get the list of Pokemon located here",
			callback:    commandExplore,
		},

		"catch": {
			name:        "catch",
			description: "Try to catch the mentionned pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect the pokemon you caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the pokemon you have caught so far",
			callback:    commandPokedex,
		},
	}
}
