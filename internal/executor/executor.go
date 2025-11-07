package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"pm/internal/detector"
	"pm/internal/project"
	"pm/internal/registry"
	"pm/internal/translator"
)

// Execute runs a translated command with automatic @types package handling for TypeScript projects
func Execute(pm detector.PackageManager, cmd *translator.Command) error {
	// Build the final command
	finalArgs := append(cmd.Command, cmd.Flags...)
	finalArgs = append(finalArgs, cmd.Args...)

	// Execute the main command
	err := run(pmBinary(pm), finalArgs...)
	if err != nil {
		return err
	}

	// Handle @types packages if needed
	if !project.IsTypeScript() {
		return nil
	}

	typesToInstall := []string{}

	for _, pkg := range cmd.Args {
		// Skip if already a @types package
		if strings.HasPrefix(pkg, "@types/") {
			continue
		}

		// Check if the package already has types
		if registry.IsTyped(pkg) {
			continue
		}

		// Check if @types package exists
		typesPackage := "@types/" + pkg
		exists, _ := registry.PackageExists(typesPackage)
		if exists {
			typesToInstall = append(typesToInstall, typesPackage)
		}
	}

	// Install @types packages as dev dependencies
	if len(typesToInstall) > 0 {
		devCommand := []string{}
		devFlag := []string{}

		switch pm {
		case detector.NPM:
			devCommand = []string{"install"}
			devFlag = []string{"--save-dev"}
		case detector.Yarn, detector.YarnBerry:
			devCommand = []string{"add"}
			devFlag = []string{"--dev"}
		case detector.Pnpm:
			devCommand = []string{"add"}
			devFlag = []string{"--save-dev"}
		case detector.Bun:
			devCommand = []string{"add"}
			devFlag = []string{"--dev"}
		}

		typesArgs := append(devCommand, devFlag...)
		typesArgs = append(typesArgs, typesToInstall...)

		err = run(pmBinary(pm), typesArgs...)
		if err != nil {
			// Don't fail if @types installation fails
			fmt.Fprintf(os.Stderr, "Warning: Failed to install @types packages: %v\n", err)
		}
	}

	return nil
}

// Run executes a command with the given package manager and arguments
func Run(pm detector.PackageManager, command translator.CommandAlias, args ...string) error {
	return RunWithFlags(pm, command, []FlagAlias{}, args...)
}

// RunWithFlags executes a command with flags using the detected package manager
func RunWithFlags(pm detector.PackageManager, command translator.CommandAlias, flags []FlagAlias, args ...string) error {
	flagArgs := make([]string, len(flags))
	for i, flag := range flags {
		switch pm {
		case detector.NPM:
			flagArgs[i] = flag[detector.NPM][0]
		case detector.Yarn:
			flagArgs[i] = flag[detector.Yarn][0]
		case detector.YarnBerry:
			if yarnValues, ok := flag[detector.YarnBerry]; ok {
				flagArgs[i] = yarnValues[0]
			} else if classicValues, ok := flag[detector.Yarn]; ok {
				flagArgs[i] = classicValues[0]
			} else {
				return fmt.Errorf("flag not available for yarn berry")
			}
		case detector.Pnpm:
			flagArgs[i] = flag[detector.Pnpm][0]
		case detector.Bun:
			flagArgs[i] = flag[detector.Bun][0]
		default:
			return fmt.Errorf("unknown package manager: %s", pm)
		}
	}

	cmd := exec.Command(pmBinary(pm), append(command[pm], append(flagArgs, args...)...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// FlagAlias maps package managers to their flag equivalents
type FlagAlias map[detector.PackageManager][]string

func run(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pmBinary(pm detector.PackageManager) string {
	switch pm {
	case detector.YarnBerry:
		return "yarn"
	default:
		return string(pm)
	}
}
