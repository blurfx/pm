package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PackageJSON struct {
	Scripts map[string]string `json:"scripts"`
}

type Script struct {
	Name    string
	Command string
}

type OrderedScripts struct {
	Scripts []Script
}

func (os *OrderedScripts) UnmarshalJSON(data []byte) error {
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

		os.Scripts = append(os.Scripts, Script{
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

func getScriptsOrdered() ([]Script, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	packageJSONPath := filepath.Join(cwd, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	var ordered OrderedScripts
	if err := json.Unmarshal(data, &ordered); err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	return ordered.Scripts, nil
}
