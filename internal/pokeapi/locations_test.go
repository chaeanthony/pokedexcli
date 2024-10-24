package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetLocation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := []byte(`{ 
        "count":1036,"next":"https://pokeapi.co/api/v2/location?offset=40&limit=20",
        "previous":"https://pokeapi.co/api/v2/location?offset=0&limit=20",
        "results":[{"name":"ruin-maniac-cave","url":"https://pokeapi.co/api/v2/location/22/"},
                    {"name":"trophy-garden","url":"https://pokeapi.co/api/v2/location/23/"}]
    }`)

		w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)
	}))
	defer server.Close()

  pokeapiClient := NewClient(5*time.Second)
	locationResp, err := pokeapiClient.GetLocationData(&server.URL) 
  if err != nil {
    t.Fatalf("Expected locations data, got %v", err)
  }

  if len(locationResp.Results) < 2 {
    t.Fatalf("Expected at least 2 locations, got %v", locationResp.Results)
  }

	expected := []string{"1036", "ruin-maniac-cave", "trophy-garden"}
  actual := []string{fmt.Sprintf("%d", locationResp.Count), locationResp.Results[0].Name, locationResp.Results[1].Name}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}


