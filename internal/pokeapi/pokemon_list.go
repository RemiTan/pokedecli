package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListPokemons(areaName string) (RespShallowPokemons, error) {
	pokemonsResp := RespShallowPokemons{}

	url := baseURL + "/location-area/" + areaName

	cachedVal, found := c.cache.Get(url)

	if found {
		if err := json.Unmarshal(cachedVal, &pokemonsResp); err != nil {
			return RespShallowPokemons{}, nil
		}

		return pokemonsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowPokemons{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return RespShallowPokemons{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RespShallowPokemons{}, err
	}

	err = json.Unmarshal(body, &pokemonsResp)
	if err != nil {
		return RespShallowPokemons{}, err
	}

	c.cache.Add(url, body)

	return pokemonsResp, nil
}
