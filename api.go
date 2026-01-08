package main

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"
	"time"

	"github.com/marco-04/godex/internal"
)

type Direction string
const (
	next Direction = "next"
	prev Direction = "prev"
)

type mapNavigation struct {
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func getLocationBatch(url string, c internal.Cache) (mapNavigation, error) {
	var data []byte

	cache, ok := c.Get(url)
	if ok {
		data = cache
	} else {
		res, err := http.Get(url)
		if err != nil {
			return mapNavigation{}, fmt.Errorf("network error: %w", err)
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return mapNavigation{}, fmt.Errorf("I/O error: %w", err)
		}

		c.Add(url, body)

		data = body
	}

	var nav mapNavigation
	if err := json.Unmarshal(data, &nav); err != nil {
		return mapNavigation{}, fmt.Errorf("json decoding error: %w", err)
	}

	return nav, nil
}

func LocationNavigator() func (Direction) ([]Location, error) {
	url := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	nextURL := ""
	prevURL := ""

	c := internal.NewCache(10 * time.Second)

	return func(direction Direction) ([]Location, error) {
		switch direction {
		case next:
			if nextURL != "" {
				url = nextURL
			}
		case prev:
			if prevURL != "" {
				url = prevURL
			}
		}

		nav, err := getLocationBatch(url, c)
		if err != nil {
			return nil, fmt.Errorf("could not get locations: %w", err)
		}

		if nav.Next != nil {
			nextURL = *nav.Next
		}
		if nav.Previous != nil {
			prevURL = *nav.Previous
		}

		return nav.Results, nil
	}
}
