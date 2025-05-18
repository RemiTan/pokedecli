package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	locationsResp := RespShallowLocations{}
	url := baseURL + "/location-area"
	if pageURL != nil {
		url = *pageURL
	}

	cachedVal, found := c.cache.Get(url)

	if found {
		if err := json.Unmarshal(cachedVal, &locationsResp); err != nil {
			return RespShallowLocations{}, nil
		}

		return locationsResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	err = json.Unmarshal(body, &locationsResp)
	if err != nil {
		return RespShallowLocations{}, err
	}

	c.cache.Add(url, body)

	return locationsResp, nil
}
