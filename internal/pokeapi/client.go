package pokeapi

import (
	"net/http"
	"time"

	"github.com/chaeanthony/pokedexcli/internal/pokecache"
)

type Client struct {
  cache       *pokecache.Cache
	httpClient  http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
    cache: pokecache.NewCache(pokecache.DefaultCacheInterval),
	}
}

