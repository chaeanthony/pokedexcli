package pokeapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type LocationResp struct {
	Count    int    `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationPokemon struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetLocationData(pageUrl *string) (LocationResp, error) {
  url := BaseURL + "/location-area"
  if pageUrl != nil { 
    url = *pageUrl
  }

  // make new request if data not found in cache
  data, ok := c.cache.Get(url)
  if !ok {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      return LocationResp{}, err
    }
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
      return LocationResp{}, err
    }
    defer resp.Body.Close()

    data, err = io.ReadAll(resp.Body)
    if err != nil {
      return LocationResp{}, err
    }

    c.cache.Add(url, data)
  } 

  // decode/unmarshal data []byte to LocationResp
  var locationData LocationResp
  if err := json.Unmarshal(data, &locationData); err != nil {
    return LocationResp{}, err 
  }

  return locationData, nil
}

func (c *Client) GetLocationPokemon(areaName string) (LocationPokemon, error){
  url := BaseURL + "/location-area/" + areaName

  data, ok := c.cache.Get(url) 
  if !ok {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
      return LocationPokemon{}, errors.New("error creating request to get data")
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
      return LocationPokemon{}, errors.New("error getting data")
    }
    defer resp.Body.Close()

    data, err = io.ReadAll(resp.Body)
    if err != nil {
      return LocationPokemon{}, errors.New("could not read data")
    }

    c.cache.Add(url, data)
  }

  var locationPokemon LocationPokemon
  if err := json.Unmarshal(data, &locationPokemon); err != nil {
    return LocationPokemon{}, err
  }

  return locationPokemon, nil
}
