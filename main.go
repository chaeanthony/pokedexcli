package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaeanthony/pokedexcli/internal/pokeapi"
	"github.com/chaeanthony/pokedexcli/repl"
)

func main() {
  // Create a channel to receive signals
  sigChan := make(chan os.Signal, 1)

  // Notify the sigChan for SIGINT and SIGTERM
  signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

  // Start your CLI tool's main logic in a separate goroutine
  pokemon := make(map[string]pokeapi.Pokemon)
  cfg := repl.PokedexConfig{Client: pokeapi.NewClient(5*time.Second), Pokemon: pokemon} 
  go repl.RunPokedex(&cfg)

  // Wait for a termination signal
  sig := <-sigChan
  fmt.Printf("\nReceived signal: %v\n", sig)

  // Perform cleanup

  fmt.Println("Graceful shutdown complete")
  os.Exit(0)
}



