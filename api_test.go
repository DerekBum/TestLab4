package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const baseURL = "https://swapi.dev/api"

type Response struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  json.RawMessage `json:"results"`
}

type Planet struct {
	Name       string `json:"name"`
	Population string `json:"population"`
}

func TestPlanetsEndpoint(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/planets", baseURL))
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(data.Results) == 0 {
		t.Fatalf("Expected non-empty response, but got empty results array")
	}

	var planets []Planet
	if err := json.Unmarshal(data.Results, &planets); err != nil {
		t.Fatalf("Error decoding planets: %v", err)
	}

	// test the first planet in the list
	planet := planets[0]
	if planet.Name != "Tatooine" {
		t.Errorf("Expected first planet to be Tatooine, but got %s", planet.Name)
	}

	if planet.Population != "200000" {
		t.Errorf("Expected Alderaan population to be 200000, but got %s", planet.Population)
	}
}

func TestPeopleEndpoint(t *testing.T) {
	// Choose a name to search for
	name := "Luke Skywalker"

	// Send a GET request to the /people endpoint
	resp, err := http.Get(fmt.Sprintf("%s/people", baseURL))
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(data.Results) == 0 {
		t.Fatalf("Expected non-empty response, but got empty results array")
	}

	var people []struct {
		Name   string `json:"name"`
		Gender string `json:"gender"`
	}
	if err := json.Unmarshal(data.Results, &people); err != nil {
		t.Fatalf("Error decoding people: %v", err)
	}

	// Search for the name in the list of people
	found := false
	for _, person := range people {
		if person.Name == name {
			found = true
			break
		}
	}

	luke := people[0]
	if luke.Name != "Luke Skywalker" {
		t.Errorf("Expected first person to be Luke Skywalker, but got %s", luke.Name)
	}

	if luke.Gender != "male" {
		t.Errorf("Expected Luke`s geder to be male, but got %s", luke.Gender)
	}

	if !found {
		t.Errorf("Expected to find a person with name %s, but none was found", name)
	}
}
