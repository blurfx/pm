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

	cmd := exec.Command(string(packageManager), args...)

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

	if err != nil {
		return err
	}

	flagArgs := make([]string, len(flags))
	for i, flag := range flags {
		switch packageManager {
		case PackageManagerNpm:
			flagArgs[i] = flag[PackageManagerNpm][0]
		case PackageManagerYarn:
			flagArgs[i] = flag[PackageManagerYarn][0]
		case PackageManagerPnpm:
			flagArgs[i] = flag[PackageManagerPnpm][0]
		case PackageManagerBun:
			flagArgs[i] = flag[PackageManagerBun][0]
		default:
			return fmt.Errorf("unknown package manager: %s", packageManager)
		}
	}

	cmd := exec.Command(string(packageManager), append(command[packageManager], append(flagArgs, args...)...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
