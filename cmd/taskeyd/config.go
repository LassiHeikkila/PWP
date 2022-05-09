package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	URL          string `json:"serviceURL"`
	AccessToken  string `json:"accessToken"`
	Organization string `json:"organization"`
}

func loadConfig(path string, c *Config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(c)
	if err != nil {
		return err
	}

	return nil
}
