package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chaeanthony/pokedexcli/internal/pokeapi"
)

type PokedexConfig struct {
  Client            pokeapi.Client
  nextLocationsUrl  *string
  prevLocationsUrl  *string
  Pokemon           map[string]pokeapi.Pokemon
}

func RunPokedex(cfg *PokedexConfig) {
  scanner := bufio.NewScanner(os.Stdin)
  cmds := GetCommands()  
  
  for {
    fmt.Print("Pokedex > ")
    scanner.Scan()
    line := scanner.Text()

    inputs := cleanInput(line)
    if len(inputs) == 0 {
      continue
    }
    
    cmdName := inputs[0]
    var args []string 
    if len(inputs) >= 1 {
      args = inputs[1:]
    }

    cmd, ok := cmds[cmdName]
    if !ok {
      fmt.Printf("Command '%s' not found\n", cmdName)
      continue
    }
    if err := cmd.callback(cfg, args...); err != nil {
      fmt.Printf("\nerror executing %s, got: %v\n", cmdName, err)
      continue
    }

  }
}


func cleanInput(line string) []string {
  words := strings.Fields(strings.TrimSpace(strings.ToLower(line)))
  return words
}
