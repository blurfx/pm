package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackageManagerType(t *testing.T) {
	tests := []struct {
		pm   PackageManager
		want string
	}{
		{NPM, "npm"},
		{Yarn, "yarn"},
		{YarnBerry, "yarn-berry"},
		{Pnpm, "pnpm"},
		{Bun, "bun"},
	}

	for _, tt := range tests {
		t.Run(string(tt.pm), func(t *testing.T) {
			if string(tt.pm) != tt.want {
				t.Errorf("PackageManager = %s, want %s", tt.pm, tt.want)
			}
		})
	}
}

func TestFindPackageJSON(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir := t.TempDir()

	// Create nested directories
	nestedDir := filepath.Join(tmpDir, "a", "b", "c")
	err := os.MkdirAll(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Create package.json in the root
	packageJSONPath := filepath.Join(tmpDir, "package.json")
	err = os.WriteFile(packageJSONPath, []byte(`{"name":"test"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Change to nested directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	err = os.Chdir(nestedDir)
	if err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Test finding package.json
	found, err := FindPackageJSON()
	if err != nil {
		t.Errorf("FindPackageJSON() error = %v", err)
	}
	// Use filepath.EvalSymlinks to resolve /private/var to /var on macOS
	foundResolved, _ := filepath.EvalSymlinks(found)
	expectedResolved, _ := filepath.EvalSymlinks(packageJSONPath)
	if foundResolved != expectedResolved {
		t.Errorf("FindPackageJSON() = %s, want %s", foundResolved, expectedResolved)
	}
}

func TestFindPackageJSONNotFound(t *testing.T) {
	// Create a temporary directory without package.json
	tmpDir := t.TempDir()

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	_, err = FindPackageJSON()
	if err == nil {
		t.Error("FindPackageJSON() expected error, got nil")
	}
}

func TestFindProjectRoot(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	nestedDir := filepath.Join(tmpDir, "src", "components")
	err := os.MkdirAll(nestedDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested dir: %v", err)
	}

	// Create package.json in the root
	err = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(`{"name":"test"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	// Change to nested directory
	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	err = os.Chdir(nestedDir)
	if err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	// Test finding project root
	root, err := FindProjectRoot()
	if err != nil {
		t.Errorf("FindProjectRoot() error = %v", err)
	}
	// Use filepath.EvalSymlinks to resolve /private/var to /var on macOS
	rootResolved, _ := filepath.EvalSymlinks(root)
	tmpDirResolved, _ := filepath.EvalSymlinks(tmpDir)
	if rootResolved != tmpDirResolved {
		t.Errorf("FindProjectRoot() = %s, want %s", rootResolved, tmpDirResolved)
	}
}

func TestDetectPackageManagerByLockfile(t *testing.T) {
	tests := []struct {
		name        string
		lockfile    string
		wantPM      PackageManager
		expectError bool
	}{
		{
			name:     "npm - package-lock.json",
			lockfile: "package-lock.json",
			wantPM:   NPM,
		},
		{
			name:     "yarn - yarn.lock",
			lockfile: "yarn.lock",
			wantPM:   Yarn,
		},
		{
			name:     "pnpm - pnpm-lock.yaml",
			lockfile: "pnpm-lock.yaml",
			wantPM:   Pnpm,
		},
		{
			name:     "bun - bun.lock",
			lockfile: "bun.lock",
			wantPM:   Bun,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create package.json
			packageJSON := filepath.Join(tmpDir, "package.json")
			err := os.WriteFile(packageJSON, []byte(`{"name":"test"}`), 0644)
			if err != nil {
				t.Fatalf("Failed to create package.json: %v", err)
			}

			// Create lockfile
			lockfilePath := filepath.Join(tmpDir, tt.lockfile)
			err = os.WriteFile(lockfilePath, []byte(``), 0644)
			if err != nil {
				t.Fatalf("Failed to create lockfile: %v", err)
			}

			// Change to temp directory
			originalWd, _ := os.Getwd()
			defer os.Chdir(originalWd)

			err = os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("Failed to change dir: %v", err)
			}

			// Test detection
			pm, err := Detect()
			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if pm != tt.wantPM {
				t.Errorf("Detect() = %s, want %s", pm, tt.wantPM)
			}
		})
	}
}

func TestDetectYarnBerryByLockfileIndicators(t *testing.T) {
	tmpDir := t.TempDir()

	packageJSON := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(packageJSON, []byte(`{"name":"test"}`), 0644); err != nil {
		t.Fatalf("Failed to create package.json: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "yarn.lock"), []byte(``), 0644); err != nil {
		t.Fatalf("Failed to create yarn.lock: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, ".yarnrc.yml"), []byte("yarnPath: .yarn/releases/yarn-berry.cjs"), 0644); err != nil {
		t.Fatalf("Failed to create .yarnrc.yml: %v", err)
	}

	originalWd, _ := os.Getwd()
	defer os.Chdir(originalWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change dir: %v", err)
	}

	pm, err := Detect()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if pm != YarnBerry {
		t.Fatalf("Detect() = %s, want YarnBerry", pm)
	}
}

func TestDetectPackageManagerByPackageJSON(t *testing.T) {
	tests := []struct {
		name        string
		packageJSON string
		wantPM      PackageManager
	}{
		{
			name: "npm specified in packageManager field",
			packageJSON: `{
				"name": "test",
				"packageManager": "npm@9.0.0"
			}`,
			wantPM: NPM,
		},
		{
			name: "yarn classic specified in packageManager field",
			packageJSON: `{
				"name": "test",
				"packageManager": "yarn@1.22.19"
			}`,
			wantPM: Yarn,
		},
		{
			name: "yarn berry specified in packageManager field",
			packageJSON: `{
				"name": "test",
				"packageManager": "yarn@3.0.0"
			}`,
			wantPM: YarnBerry,
		},
		{
			name: "pnpm specified in packageManager field",
			packageJSON: `{
				"name": "test",
				"packageManager": "pnpm@8.0.0"
			}`,
			wantPM: Pnpm,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create package.json
			packageJSONPath := filepath.Join(tmpDir, "package.json")
			err := os.WriteFile(packageJSONPath, []byte(tt.packageJSON), 0644)
			if err != nil {
				t.Fatalf("Failed to create package.json: %v", err)
			}

			// Change to temp directory
			originalWd, _ := os.Getwd()
			defer os.Chdir(originalWd)

			err = os.Chdir(tmpDir)
			if err != nil {
				t.Fatalf("Failed to change dir: %v", err)
			}

			// Test detection
			pm, err := Detect()
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if pm != tt.wantPM {
				t.Errorf("Detect() = %s, want %s", pm, tt.wantPM)
			}
		})
	}
}
