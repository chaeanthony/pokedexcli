package repl

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"os"

)

type cliCommand struct {
	name        string
	description string
  callback    func(*PokedexConfig, ...string) error
}

func GetCommands() map[string]cliCommand {
  return map[string]cliCommand{
    "catch" : {
      name: "catch",
      description: "Attempt to catch a pokemon",
      callback: commandCatch,
    },
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
    "explore" : {
      name: "explore",
      description: "Get pokemon in location area",
      callback: commandExplore,
    },
    "inspect" : {
      name: "inspect",
      description: "Inspect pokemon details",
      callback: commandInspect,
    },
    "map": {
      name: "map",
      description: "Get next page of locations",
      callback: commandMap,
    },
    "mapb": {
      name: "mapb",
      description: "Gets previous page of locations",
      callback: commandMapb,
    },
    "pokedex": {
      name: "pokemon",
      description: "List my pokemon",
      callback: commandPokedex,
    },
  }
}

//TODO add naming to prevent duplicate pokemon in pokedex
func commandCatch(cfg *PokedexConfig, args ...string) error {
  if len(args) != 1 {
    return errors.New("Must provide a single pokemon name")
  }
  input := args[0]
  pokemon, err := cfg.Client.GetPokemon(input)
  if err != nil {
    return err
  }

  fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
  random := rand.IntN(pokemon.BaseExperience)
  if random >= 40 {
    fmt.Printf("%s was caught!\n", pokemon.Name)
    fmt.Printf("Base experience: %d\n", pokemon.BaseExperience)
    cfg.Pokemon[pokemon.Name] = pokemon
  } else {
    fmt.Printf("%s escaped!\n", pokemon.Name)
  }

  return nil
}

func commandHelp(cfg *PokedexConfig, args ...string) error {
  fmt.Println("Welcome to the Pokedex!")
  fmt.Println("Commands:")
  fmt.Println()
  for _, cmd := range GetCommands() {
    fmt.Printf(" %s: %s\n", cmd.name, cmd.description)
  }
  return nil
}

func commandExit(cfg *PokedexConfig, args ...string) error {
  fmt.Println("Exiting Pokedex")
  os.Exit(0)
  return nil
}

func commandExplore(cfg *PokedexConfig, args ...string) error {
  if len(args) != 1 {
    return errors.New("You must provide a single location name")
  }
  fmt.Println("Exploring area ...")
  areaName := args[0]
  
  // get pokemon data from location area name
  locationPokemon, err := cfg.Client.GetLocationPokemon(areaName)
  if err != nil {
    return fmt.Errorf("error finding pokemon. got: %v\n", err)
  } else if len(locationPokemon.PokemonEncounters) == 0 {
    return fmt.Errorf("no pokemon found in location: %s\n", areaName)
  }
  
  // parse data to list of pokemon names in location area
  names := []string{}
  for _, encounter := range locationPokemon.PokemonEncounters {
    pokemon := encounter.Pokemon
    names = append(names, pokemon.Name)
  }
  if len(names) == 0 {
    return fmt.Errorf("no pokemon found in location %s\n", areaName)
  }
  fmt.Println("Found pokemon: ")
  for _, name := range names {
    fmt.Printf(" -%s\n", name)
  }
  fmt.Println()
  
  return nil
}

//TODO: add prompui selection
func commandInspect(cfg *PokedexConfig, args ...string) error {
  if len(args) != 1 {
    return errors.New("please provide a single pokemon name to inspect")
  }
  
  name := args[0]
  pokemon, ok := cfg.Pokemon[name]
  if !ok {
    return errors.New("you have not caught that pokemon")
  }

  fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)
  fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
  fmt.Printf("Types:\n")
  for _, v := range pokemon.Types {
    fmt.Printf(" -%s\n", v.Type.Name)
  }

  return nil
}

func commandMap(cfg *PokedexConfig, args ...string) error {
  resp, err := cfg.Client.GetLocationData(cfg.nextLocationsUrl)
  if err != nil {
    fmt.Println("error getting locations data")
    return err
  }

  cfg.nextLocationsUrl = resp.Next
  cfg.prevLocationsUrl = resp.Previous

  fmt.Println() 
  fmt.Printf("Number of locations: %d\n", resp.Count)
  fmt.Println()

  for _, loc := range resp.Results {
    fmt.Printf(" -%s\n", loc.Name)
  }
  fmt.Println()

  return nil
}

func commandMapb(cfg *PokedexConfig, args ...string) error {
  if cfg.prevLocationsUrl == nil {
    return errors.New("On first page. No previous pages.")
  }

  resp, err := cfg.Client.GetLocationData(cfg.prevLocationsUrl)
  if err != nil {
    fmt.Println("error getting locations data")
    return err
  }

  cfg.nextLocationsUrl = resp.Next
  cfg.prevLocationsUrl = resp.Previous
  
  fmt.Println()
  fmt.Printf("Number of locations: %d\n", resp.Count)
  fmt.Println()

  for _, loc := range resp.Results {
    fmt.Printf(" -%s\n", loc.Name)
  }
  fmt.Println()

  return nil
  
}

func commandPokedex(cfg *PokedexConfig, args ...string) error {
  if len(cfg.Pokemon) == 0 {
    return errors.New("No Pokemon yet")
  }

  for k, _ := range cfg.Pokemon {
    fmt.Printf(" -%s\n", k)
  }
  fmt.Println()
  return nil
}


