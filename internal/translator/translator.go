package translator

import (
	"pm/internal/detector"
	"strings"
)

// Translator translates universal package manager commands to package-manager-specific commands
type Translator struct {
	packageManager detector.PackageManager
}

// New creates a new command translator for the given package manager
func New(pm detector.PackageManager) *Translator {
	return &Translator{
		packageManager: pm,
	}
}

// Command represents a translated command with its flags and arguments
type Command struct {
	Command []string
	Flags   []string
	Args    []string
}

// Translate translates a universal command to a package-manager-specific command
func (t *Translator) Translate(packageManager detector.PackageManager, args []string) *Command {
	if len(args) == 0 {
		return &Command{}
	}

	baseCommand := args[0]
	remainingArgs := args[1:]

	switch baseCommand {
	case "i", "install":
		return t.translateInstall(remainingArgs)
	case "add":
		return t.translateAdd(remainingArgs)
	case "rm", "remove", "uninstall", "un":
		return t.translateUninstall(remainingArgs)
	case "ci":
		return t.translateCI(remainingArgs)
	default:
		if IsBuiltIn(packageManager, baseCommand) {
			return &Command{
				Command: []string{baseCommand},
				Args:    remainingArgs,
			}
		}
		return &Command{
			Command: []string{"run"},
			Args:    args,
		}
	}
}

func (t *Translator) translateInstall(args []string) *Command {
	parsed := t.parseArgs(args)

	if len(parsed.packages) > 0 {
		return t.translateAdd(args)
	}

	var command []string
	switch t.packageManager {
	case detector.NPM:
		command = []string{"install"}
	case detector.Yarn, detector.YarnBerry:
		command = []string{"install"}
	case detector.Pnpm:
		command = []string{"install"}
	case detector.Bun:
		command = []string{"install"}
	}

	flags := t.translateInstallFlags(parsed.flags)

	return &Command{
		Command: command,
		Flags:   flags,
		Args:    parsed.packages,
	}
}

func (t *Translator) translateAdd(args []string) *Command {
	parsed := t.parseArgs(args)

	var command []string
	switch t.packageManager {
	case detector.NPM:
		command = []string{"install"}
	case detector.Yarn, detector.YarnBerry:
		command = []string{"add"}
	case detector.Pnpm:
		command = []string{"add"}
	case detector.Bun:
		command = []string{"add"}
	}

	flags := t.translateAddFlags(parsed.flags)

	return &Command{
		Command: command,
		Flags:   flags,
		Args:    parsed.packages,
	}
}

func (t *Translator) translateUninstall(args []string) *Command {
	parsed := t.parseArgs(args)

	var command []string
	switch t.packageManager {
	case detector.NPM:
		command = []string{"uninstall"}
	case detector.Yarn, detector.YarnBerry:
		command = []string{"remove"}
	case detector.Pnpm:
		command = []string{"remove"}
	case detector.Bun:
		command = []string{"remove"}
	}

	flags := t.translateUninstallFlags(parsed.flags)

	return &Command{
		Command: command,
		Flags:   flags,
		Args:    parsed.packages,
	}
}

func (t *Translator) translateCI(args []string) *Command {
	var command []string
	switch t.packageManager {
	case detector.NPM:
		command = []string{"ci"}
	case detector.Yarn:
		command = []string{"install", "--frozen-lockfile"}
	case detector.YarnBerry:
		command = []string{"install", "--immutable"}
	case detector.Pnpm:
		command = []string{"install", "--frozen-lockfile"}
	case detector.Bun:
		command = []string{"install", "--frozen-lockfile"}
	}

	return &Command{
		Command: command,
		Args:    args,
	}
}

type parsedArgs struct {
	packages []string
	flags    map[string]string
}

func (t *Translator) parseArgs(args []string) *parsedArgs {
	parsed := &parsedArgs{
		packages: []string{},
		flags:    make(map[string]string),
	}

	i := 0
	for i < len(args) {
		arg := args[i]

		if strings.HasPrefix(arg, "--") {
			flagName := strings.TrimPrefix(arg, "--")

			if strings.Contains(flagName, "=") {
				parts := strings.SplitN(flagName, "=", 2)
				parsed.flags[parts[0]] = parts[1]
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && t.isFlagWithValue(flagName) {
				parsed.flags[flagName] = args[i+1]
				i++
			} else {
				parsed.flags[flagName] = ""
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 && !strings.HasPrefix(arg, "--") {
			flagName := strings.TrimPrefix(arg, "-")

			// Handle multiple short flags like -gD
			for j, char := range flagName {
				shortFlag := string(char)
				fullFlag := t.expandShortFlag(shortFlag)

				// Only the last flag in a group can have a value
				if j == len(flagName)-1 && t.isFlagWithValue(fullFlag) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					parsed.flags[fullFlag] = args[i+1]
					i++
				} else {
					parsed.flags[fullFlag] = ""
				}
			}
		} else {
			parsed.packages = append(parsed.packages, arg)
		}

		i++
	}

	return parsed
}

func (t *Translator) isFlagWithValue(flag string) bool {
	flagsWithValues := []string{
		"omit", "registry", "tag", "workspace", "workspaces",
		"production", "only", "also", "save-bundle", "save-exact",
		"loglevel", "logs-max", "logs-dir", "script-shell",
		"cache-folder", "cache-dir", "prefix", "userconfig",
	}

	for _, f := range flagsWithValues {
		if f == flag {
			return true
		}
	}

	return false
}

func (t *Translator) expandShortFlag(short string) string {
	shortFlagMap := map[string]string{
		"D": "save-dev",
		"P": "save-peer",
		"O": "save-optional",
		"E": "save-exact",
		"g": "global",
		"S": "save",
		"B": "save-bundle",
		"f": "force",
		"s": "silent",
		"d": "loglevel",
	}

	if full, ok := shortFlagMap[short]; ok {
		return full
	}

	return short
}

func (t *Translator) translateInstallFlags(flags map[string]string) []string {
	var translated []string

	for flag, value := range flags {
		switch flag {
		case "frozen-lockfile":
			translated = append(translated, t.translateFrozenLockfileFlag()...)
		case "omit":
			translated = append(translated, t.translateOmitFlag(value)...)
		case "global", "g":
			translated = append(translated, t.translateGlobalFlag()...)
		case "production":
			translated = append(translated, t.translateProductionFlag()...)
		default:
			if value != "" {
				translated = append(translated, "--"+flag, value)
			} else {
				translated = append(translated, "--"+flag)
			}
		}
	}

	return translated
}

func (t *Translator) translateAddFlags(flags map[string]string) []string {
	var translated []string

	for flag, value := range flags {
		switch flag {
		case "save-dev", "D":
			translated = append(translated, t.translateDevFlag()...)
		case "save-peer", "P":
			translated = append(translated, t.translatePeerFlag()...)
		case "save-optional", "O":
			translated = append(translated, t.translateOptionalFlag()...)
		case "save-exact", "E":
			translated = append(translated, t.translateExactFlag()...)
		case "global", "g":
			translated = append(translated, t.translateGlobalFlag()...)
		case "omit":
			translated = append(translated, t.translateOmitFlag(value)...)
		default:
			if value != "" {
				translated = append(translated, "--"+flag, value)
			} else {
				translated = append(translated, "--"+flag)
			}
		}
	}

	return translated
}

func (t *Translator) translateUninstallFlags(flags map[string]string) []string {
	var translated []string

	for flag, value := range flags {
		switch flag {
		case "save-dev", "D":
			translated = append(translated, t.translateDevFlag()...)
		case "global", "g":
			translated = append(translated, t.translateGlobalFlag()...)
		default:
			if value != "" {
				translated = append(translated, "--"+flag, value)
			} else {
				translated = append(translated, "--"+flag)
			}
		}
	}

	return translated
}

func (t *Translator) translateFrozenLockfileFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{}
	case detector.YarnBerry:
		return []string{"--immutable"}
	case detector.Yarn, detector.Pnpm, detector.Bun:
		return []string{"--frozen-lockfile"}
	}
	return []string{}
}

func (t *Translator) translateDevFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--save-dev"}
	case detector.Yarn, detector.YarnBerry:
		return []string{"--dev"}
	case detector.Pnpm:
		return []string{"--save-dev"}
	case detector.Bun:
		return []string{"--dev"}
	}
	return []string{}
}

func (t *Translator) translatePeerFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--save-peer"}
	case detector.Yarn, detector.YarnBerry:
		return []string{"--peer"}
	case detector.Pnpm:
		return []string{"--save-peer"}
	case detector.Bun:
		return []string{"--peer"}
	}
	return []string{}
}

func (t *Translator) translateOptionalFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--save-optional"}
	case detector.Yarn, detector.YarnBerry:
		return []string{"--optional"}
	case detector.Pnpm:
		return []string{"--save-optional"}
	case detector.Bun:
		return []string{"--optional"}
	}
	return []string{}
}

func (t *Translator) translateExactFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--save-exact"}
	case detector.Yarn, detector.YarnBerry:
		return []string{"--exact"}
	case detector.Pnpm:
		return []string{"--save-exact"}
	case detector.Bun:
		return []string{"--exact"}
	}
	return []string{}
}

func (t *Translator) translateGlobalFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--global"}
	case detector.Yarn, detector.YarnBerry:
		return []string{"--global"}
	case detector.Pnpm:
		return []string{"--global"}
	case detector.Bun:
		return []string{"--global"}
	}
	return []string{}
}

func (t *Translator) translateProductionFlag() []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--production"}
	case detector.Yarn, detector.YarnBerry:
		return []string{"--production"}
	case detector.Pnpm:
		return []string{"--prod"}
	case detector.Bun:
		return []string{"--production"}
	}
	return []string{}
}

func (t *Translator) translateOmitFlag(value string) []string {
	switch t.packageManager {
	case detector.NPM:
		return []string{"--omit", value}
	case detector.Yarn, detector.YarnBerry:
		if value == "dev" {
			return []string{"--production"}
		}
		return []string{}
	case detector.Pnpm:
		if value == "dev" {
			return []string{"--prod"}
		}
		return []string{"--omit", value}
	case detector.Bun:
		if value == "dev" {
			return []string{"--production"}
		}
		return []string{}
	}
	return []string{}
}
