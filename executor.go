package main

import (
	"fmt"
	"os"
	"os/exec"
)

func PassThrough(args ...string) error {
	packageManager, err := DetectPackageManager()
	if err != nil {
		return err
	}

	cmd := exec.Command(packageManager, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Exec(command CommandAlias, args ...string) error {
	return ExecWithFlag(command, []flagAlias{}, args...)
}

func ExecWithFlag(command CommandAlias, flags []flagAlias, args ...string) error {
	packageManager, err := DetectPackageManager()
	packageCommands := map[string][]string{
		"npm":  command.NPM,
		"yarn": command.Yarn,
		"pnpm": command.Pnpm,
	}
	if err != nil {
		return err
	}

	flagArgs := make([]string, len(flags))
	for i, flag := range flags {
		if packageManager == "npm" {
			flagArgs[i] = flag.NPM[0]
		} else if packageManager == "yarn" {
			flagArgs[i] = flag.Yarn[0]
		} else if packageManager == "pnpm" {
			flagArgs[i] = flag.Pnpm[0]
		} else {
			return fmt.Errorf("unknown package manager: %s", packageManager)
		}
	}

	cmd := exec.Command(packageManager, append(packageCommands[packageManager], append(flagArgs, args...)...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
