package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Scripts         map[string]string `json:"scripts"`
	OrderedScripts  []Script
}

type Script struct {
	Name    string
	Command string
}

func (p *PackageJSON) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	scriptsData, exists := raw["scripts"]
	if !exists {
		return nil
	}

	dec := json.NewDecoder(bytes.NewReader(scriptsData))

	token, err := dec.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('{') {
		return fmt.Errorf("expected '{' but got %v", token)
	}

	for dec.More() {
		token, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return fmt.Errorf("expected string key but got %v", token)
		}

		// Read value
		var value string
		if err := dec.Decode(&value); err != nil {
			return err
		}

		p.OrderedScripts = append(p.OrderedScripts, Script{
			Name:    key,
			Command: value,
		})
	}

	token, err = dec.Token()
	if err != nil {
		return err
	}
	if token != json.Delim('}') {
		return fmt.Errorf("expected '}' but got %v", token)
	}

	return nil
}
