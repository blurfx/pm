package main

import (
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
	packageManager, err := DetectPackageManager()
	packageCommands := map[string][]string{
		"npm":  command.NPM,
		"yarn": command.Yarn,
		"pnpm": command.Pnpm,
	}
	if err != nil {
		return err
	}

	cmd := exec.Command(packageManager, append(packageCommands[packageManager], args...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
