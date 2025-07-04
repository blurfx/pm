package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ExecuteTranslatedCommand(pm PackageManager, translated *TranslatedCommand) error {
	// Build the final command
	finalArgs := append(translated.Command, translated.Flags...)
	finalArgs = append(finalArgs, translated.Args...)

	fmt.Println("finalArgs", finalArgs)
	// Execute the main command
	err := executeCommand(string(pm), finalArgs...)
	if err != nil {
		return err
	}

	fmt.Println("translated.Args", translated.Args)
	// Handle @types packages if needed
	typesToInstall := []string{}

	for _, pkg := range translated.Args {
		// Skip if already a @types package
		if strings.HasPrefix(pkg, "@types/") {
			continue
		}

		// Check if the package already has types
		if IsTypedPackage(pkg) {
			continue
		}

		// Check if @types package exists
		typesPackage := "@types/" + pkg
		exists, _ := CheckPackageExists(typesPackage)
		if exists {
			typesToInstall = append(typesToInstall, typesPackage)
		}
	}

	// Install @types packages as dev dependencies
	if len(typesToInstall) > 0 {
		devCommand := []string{}
		devFlag := []string{}

		switch pm {
		case PackageManagerNpm:
			devCommand = []string{"install"}
			devFlag = []string{"--save-dev"}
		case PackageManagerYarn:
			devCommand = []string{"add"}
			devFlag = []string{"--dev"}
		case PackageManagerPnpm:
			devCommand = []string{"add"}
			devFlag = []string{"--save-dev"}
		case PackageManagerBun:
			devCommand = []string{"add"}
			devFlag = []string{"--dev"}
		}

		typesArgs := append(devCommand, devFlag...)
		typesArgs = append(typesArgs, typesToInstall...)

		err = executeCommand(string(pm), typesArgs...)
		if err != nil {
			// Don't fail if @types installation fails
			fmt.Fprintf(os.Stderr, "Warning: Failed to install @types packages: %v\n", err)
		}
	}

	return nil
}

func executeCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
