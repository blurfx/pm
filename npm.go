package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func CheckPackageExists(packageName string) (bool, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	resp, err := http.Head(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
}

type PackageInfo struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions map[string]struct {
		Types   string `json:"types"`
		Typings string `json:"typings"`
	} `json:"versions"`
}

func IsTypedPackage(packageName string) bool {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", packageName)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var data PackageInfo
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return false
	}

	latestVersion := data.DistTags.Latest
	packageInfo := data.Versions[latestVersion]

	return packageInfo.Types != "" || packageInfo.Typings != ""
}

type LocalPackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func IsTypeScriptPackage() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	tsconfigPath := filepath.Join(cwd, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); err == nil {
		return true
	}

	packageJSONPath := filepath.Join(cwd, "package.json")
	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return false
	}

	var pkg LocalPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}

	if _, ok := pkg.Dependencies["typescript"]; ok {
		return true
	}
	if _, ok := pkg.DevDependencies["typescript"]; ok {
		return true
	}

	return false
}
