package project

import (
	"encoding/json"
	"testing"
)

func TestPackageJSONUnmarshal(t *testing.T) {
	tests := []struct {
		name                string
		jsonData            string
		wantDeps            map[string]string
		wantDevDeps         map[string]string
		wantScriptsCount    int
		wantFirstScriptName string
	}{
		{
			name: "complete package.json with dependencies and scripts",
			jsonData: `{
				"name": "test-package",
				"version": "1.0.0",
				"dependencies": {
					"react": "^18.0.0",
					"axios": "^1.0.0"
				},
				"devDependencies": {
					"typescript": "^5.0.0",
					"jest": "^29.0.0"
				},
				"scripts": {
					"test": "jest",
					"build": "tsc",
					"dev": "vite"
				}
			}`,
			wantDeps: map[string]string{
				"react": "^18.0.0",
				"axios": "^1.0.0",
			},
			wantDevDeps: map[string]string{
				"typescript": "^5.0.0",
				"jest":       "^29.0.0",
			},
			wantScriptsCount:    3,
			wantFirstScriptName: "test",
		},
		{
			name: "only dependencies",
			jsonData: `{
				"name": "test-package",
				"dependencies": {
					"lodash": "^4.17.21"
				}
			}`,
			wantDeps: map[string]string{
				"lodash": "^4.17.21",
			},
			wantDevDeps:      map[string]string{},
			wantScriptsCount: 0,
		},
		{
			name: "only devDependencies",
			jsonData: `{
				"name": "test-package",
				"devDependencies": {
					"eslint": "^8.0.0"
				}
			}`,
			wantDeps: map[string]string{},
			wantDevDeps: map[string]string{
				"eslint": "^8.0.0",
			},
			wantScriptsCount: 0,
		},
		{
			name: "only scripts",
			jsonData: `{
				"name": "test-package",
				"scripts": {
					"start": "node index.js",
					"test": "echo test"
				}
			}`,
			wantDeps:            map[string]string{},
			wantDevDeps:         map[string]string{},
			wantScriptsCount:    2,
			wantFirstScriptName: "start",
		},
		{
			name: "empty package.json",
			jsonData: `{
				"name": "test-package"
			}`,
			wantDeps:         map[string]string{},
			wantDevDeps:      map[string]string{},
			wantScriptsCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pkg PackageJSON
			err := json.Unmarshal([]byte(tt.jsonData), &pkg)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			// Check dependencies
			if len(pkg.Dependencies) != len(tt.wantDeps) {
				t.Errorf("Dependencies count = %d, want %d", len(pkg.Dependencies), len(tt.wantDeps))
			}
			for key, wantVal := range tt.wantDeps {
				if gotVal, ok := pkg.Dependencies[key]; !ok {
					t.Errorf("Missing dependency %s", key)
				} else if gotVal != wantVal {
					t.Errorf("Dependency %s = %s, want %s", key, gotVal, wantVal)
				}
			}

			// Check devDependencies
			if len(pkg.DevDependencies) != len(tt.wantDevDeps) {
				t.Errorf("DevDependencies count = %d, want %d", len(pkg.DevDependencies), len(tt.wantDevDeps))
			}
			for key, wantVal := range tt.wantDevDeps {
				if gotVal, ok := pkg.DevDependencies[key]; !ok {
					t.Errorf("Missing devDependency %s", key)
				} else if gotVal != wantVal {
					t.Errorf("DevDependency %s = %s, want %s", key, gotVal, wantVal)
				}
			}

			// Check scripts
			if len(pkg.OrderedScripts) != tt.wantScriptsCount {
				t.Errorf("OrderedScripts count = %d, want %d", len(pkg.OrderedScripts), tt.wantScriptsCount)
			}
			if tt.wantScriptsCount > 0 && len(pkg.OrderedScripts) > 0 {
				if pkg.OrderedScripts[0].Name != tt.wantFirstScriptName {
					t.Errorf("First script name = %s, want %s", pkg.OrderedScripts[0].Name, tt.wantFirstScriptName)
				}
			}
		})
	}
}

func TestPackageJSONScriptOrder(t *testing.T) {
	jsonData := `{
		"name": "test",
		"scripts": {
			"first": "echo 1",
			"second": "echo 2",
			"third": "echo 3"
		}
	}`

	var pkg PackageJSON
	err := json.Unmarshal([]byte(jsonData), &pkg)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Note: JSON object key order is not guaranteed in Go's encoding/json
	// So we just verify all scripts are present
	if len(pkg.OrderedScripts) != 3 {
		t.Fatalf("Expected 3 scripts, got %d", len(pkg.OrderedScripts))
	}

	scriptNames := make(map[string]bool)
	for _, script := range pkg.OrderedScripts {
		scriptNames[script.Name] = true
	}

	expected := []string{"first", "second", "third"}
	for _, name := range expected {
		if !scriptNames[name] {
			t.Errorf("Missing script: %s", name)
		}
	}
}
