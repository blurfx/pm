package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestTranslator(t *testing.T) {
	testCases := []struct {
		name           string
		packageManager PackageManager
		input          []string
		expected       string
	}{
		// Test basic install with flags
		{
			name:           "npm install with -g and -D flags",
			packageManager: PackageManagerNpm,
			input:          []string{"i", "-g", "axios", "-D", "--omit", "dev"},
			expected:       "install --global --save-dev --omit dev axios",
		},
		{
			name:           "yarn add with -g and -D flags",
			packageManager: PackageManagerYarn,
			input:          []string{"i", "-g", "axios", "-D", "--omit", "dev"},
			expected:       "add --global --dev --production axios",
		},
		{
			name:           "pnpm add with -g and -D flags",
			packageManager: PackageManagerPnpm,
			input:          []string{"i", "-g", "axios", "-D", "--omit", "dev"},
			expected:       "add --global --save-dev --prod axios",
		},

		// Test ci command
		{
			name:           "npm ci",
			packageManager: PackageManagerNpm,
			input:          []string{"ci"},
			expected:       "ci",
		},
		{
			name:           "yarn ci",
			packageManager: PackageManagerYarn,
			input:          []string{"ci"},
			expected:       "install --frozen-lockfile",
		},
		{
			name:           "pnpm ci",
			packageManager: PackageManagerPnpm,
			input:          []string{"ci"},
			expected:       "install --frozen-lockfile",
		},

		// Test install without packages
		{
			name:           "npm install",
			packageManager: PackageManagerNpm,
			input:          []string{"install"},
			expected:       "install",
		},

		// Test add command
		{
			name:           "npm add",
			packageManager: PackageManagerNpm,
			input:          []string{"add", "express", "react"},
			expected:       "install express react",
		},
		{
			name:           "yarn add",
			packageManager: PackageManagerYarn,
			input:          []string{"add", "express", "react"},
			expected:       "add express react",
		},

		// Test remove/uninstall
		{
			name:           "npm uninstall",
			packageManager: PackageManagerNpm,
			input:          []string{"rm", "express"},
			expected:       "uninstall express",
		},
		{
			name:           "yarn remove",
			packageManager: PackageManagerYarn,
			input:          []string{"rm", "express"},
			expected:       "remove express",
		},

		// Test run command
		{
			name:           "run script",
			packageManager: PackageManagerNpm,
			input:          []string{"run", "test"},
			expected:       "run test",
		},

		// test package.json scripts
		{
			name:           "run script with `run` command omitted",
			packageManager: PackageManagerNpm,
			input:          []string{"test"},
			expected:       "run test",
		},
		{
			name:           "run script with `run` command omitted",
			packageManager: PackageManagerNpm,
			input:          []string{"dev"},
			expected:       "run dev",
		},
		{
			name:           "run script with `run` command omitted",
			packageManager: PackageManagerNpm,
			input:          []string{"dev --host 0.0.0.0"},
			expected:       "run dev --host 0.0.0.0",
		},
	}

	for _, tc := range testCases {
		translator := NewCommandTranslator(tc.packageManager)
		result := translator.Translate(tc.packageManager, tc.input)

		// Build the actual command string
		actualParts := append(result.Command, result.Flags...)
		actualParts = append(actualParts, result.Args...)
		actual := strings.Join(actualParts, " ")

		if actual != tc.expected {
			fmt.Printf("FAIL: %s\n", tc.name)
			fmt.Printf("  Input: %v\n", tc.input)
			fmt.Printf("  Expected: %s\n", tc.expected)
			fmt.Printf("  Actual: %s\n", actual)
		} else {
			fmt.Printf("PASS: %s\n", tc.name)
		}
	}
}
