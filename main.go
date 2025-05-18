package main

import (
	"strings"
	"time"

	"github.com/RemiTan/pokedexcli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Second)
	cfg := &config{
		pokeapiClient: pokeClient,
		Pokemons:      make(map[string]pokeapi.Pokemon),
	}

	startRepl(cfg)
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	return strings.Fields(lowerText)
}
