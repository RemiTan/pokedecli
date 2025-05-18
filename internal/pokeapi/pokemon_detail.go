package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) DetailPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	cachedVal, found := c.cache.Get(url)

	if found {
		pokemon := Pokemon{}
		if err := json.Unmarshal(cachedVal, &pokemon); err != nil {
			return Pokemon{}, nil
		}

		return pokemon, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	// fmt.Printf("body: %v\n", string(body))

	pokemon := Pokemon{}
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, body)

	return pokemon, nil
}
