package main

import (
	"fmt"
	"strings"
)

type CommandTranslator struct {
	packageManager PackageManager
}

func NewCommandTranslator(pm PackageManager) *CommandTranslator {
	return &CommandTranslator{
		packageManager: pm,
	}
}

type TranslatedCommand struct {
	Command []string
	Flags   []string
	Args    []string
}

func (ct *CommandTranslator) Translate(args []string) *TranslatedCommand {
	if len(args) == 0 {
		return &TranslatedCommand{}
	}

	baseCommand := args[0]
	remainingArgs := args[1:]

	switch baseCommand {
	case "i", "install":
		return ct.translateInstall(remainingArgs)
	case "add":
		return ct.translateAdd(remainingArgs)
	case "rm", "remove", "uninstall", "un":
		return ct.translateUninstall(remainingArgs)
	case "ci":
		return ct.translateCI(remainingArgs)
	case "run":
		return ct.translateRun(remainingArgs)
	default:
		return &TranslatedCommand{
			Command: []string{baseCommand},
			Args:    remainingArgs,
		}
	}
}

func (ct *CommandTranslator) translateInstall(args []string) *TranslatedCommand {
	parsed := ct.parseArgs(args)

	if len(parsed.packages) > 0 {
		return ct.translateAdd(args)
	}

	command := []string{}
	switch ct.packageManager {
	case PackageManagerNpm:
		command = []string{"install"}
	case PackageManagerYarn:
		command = []string{"install"}
	case PackageManagerPnpm:
		command = []string{"install"}
	case PackageManagerBun:
		command = []string{"install"}
	}

	flags := ct.translateInstallFlags(parsed.flags)

	return &TranslatedCommand{
		Command: command,
		Flags:   flags,
		Args:    parsed.packages,
	}
}

func (ct *CommandTranslator) translateAdd(args []string) *TranslatedCommand {
	parsed := ct.parseArgs(args)

	command := []string{}
	switch ct.packageManager {
	case PackageManagerNpm:
		command = []string{"install"}
	case PackageManagerYarn:
		command = []string{"add"}
	case PackageManagerPnpm:
		command = []string{"add"}
	case PackageManagerBun:
		command = []string{"add"}
	}

	flags := ct.translateAddFlags(parsed.flags)
	fmt.Println("flags", flags)

	// Check if this is a dev/peer/optional/global install
	// _, hasDev := parsed.flags["save-dev"]
	// _, hasD := parsed.flags["D"]
	// _, hasPeer := parsed.flags["save-peer"]
	// _, hasP := parsed.flags["P"]
	// _, hasOptional := parsed.flags["save-optional"]
	// _, hasO := parsed.flags["O"]
	// _, hasGlobal := parsed.flags["global"]
	// _, hasG := parsed.flags["g"]

	// isSpecialInstall := hasDev || hasD || hasPeer || hasP || hasOptional || hasO || hasGlobal || hasG

	// Store info for @types handling
	result := &TranslatedCommand{
		Command: command,
		Flags:   flags,
		Args:    parsed.packages,
	}

	// Mark if we need to handle @types (this will be handled by the executor)
	// if !isSpecialInstall {
	// 	result.ShouldHandleTypes = true
	// }

	return result
}

func (ct *CommandTranslator) translateUninstall(args []string) *TranslatedCommand {
	parsed := ct.parseArgs(args)

	command := []string{}
	switch ct.packageManager {
	case PackageManagerNpm:
		command = []string{"uninstall"}
	case PackageManagerYarn:
		command = []string{"remove"}
	case PackageManagerPnpm:
		command = []string{"remove"}
	case PackageManagerBun:
		command = []string{"remove"}
	}

	flags := ct.translateUninstallFlags(parsed.flags)

	return &TranslatedCommand{
		Command: command,
		Flags:   flags,
		Args:    parsed.packages,
	}
}

func (ct *CommandTranslator) translateCI(args []string) *TranslatedCommand {
	command := []string{}
	switch ct.packageManager {
	case PackageManagerNpm:
		command = []string{"ci"}
	case PackageManagerYarn:
		command = []string{"install", "--frozen-lockfile"}
	case PackageManagerPnpm:
		command = []string{"install", "--frozen-lockfile"}
	case PackageManagerBun:
		command = []string{"install", "--frozen-lockfile"}
	}

	return &TranslatedCommand{
		Command: command,
		Args:    args,
	}
}

func (ct *CommandTranslator) translateRun(args []string) *TranslatedCommand {
	command := []string{"run"}

	return &TranslatedCommand{
		Command: command,
		Args:    args,
	}
}

type ParsedArgs struct {
	packages []string
	flags    map[string]string
}

func (ct *CommandTranslator) parseArgs(args []string) *ParsedArgs {
	parsed := &ParsedArgs{
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
			} else if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") && ct.isFlagWithValue(flagName) {
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
				fullFlag := ct.expandShortFlag(shortFlag)

				// Only the last flag in a group can have a value
				if j == len(flagName)-1 && ct.isFlagWithValue(fullFlag) && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
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

func (ct *CommandTranslator) isFlagWithValue(flag string) bool {
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

func (ct *CommandTranslator) expandShortFlag(short string) string {
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

func (ct *CommandTranslator) translateInstallFlags(flags map[string]string) []string {
	translated := []string{}

	for flag, value := range flags {
		switch flag {
		case "frozen-lockfile":
			if ct.packageManager != PackageManagerNpm {
				translated = append(translated, "--frozen-lockfile")
			}
		case "omit":
			translated = append(translated, ct.translateOmitFlag(value)...)
		case "global", "g":
			translated = append(translated, ct.translateGlobalFlag()...)
		case "production":
			translated = append(translated, ct.translateProductionFlag()...)
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

func (ct *CommandTranslator) translateAddFlags(flags map[string]string) []string {
	translated := []string{}

	for flag, value := range flags {
		switch flag {
		case "save-dev", "D":
			translated = append(translated, ct.translateDevFlag()...)
		case "save-peer", "P":
			translated = append(translated, ct.translatePeerFlag()...)
		case "save-optional", "O":
			translated = append(translated, ct.translateOptionalFlag()...)
		case "save-exact", "E":
			translated = append(translated, ct.translateExactFlag()...)
		case "global", "g":
			translated = append(translated, ct.translateGlobalFlag()...)
		case "omit":
			translated = append(translated, ct.translateOmitFlag(value)...)
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

func (ct *CommandTranslator) translateUninstallFlags(flags map[string]string) []string {
	translated := []string{}

	for flag, value := range flags {
		switch flag {
		case "save-dev", "D":
			translated = append(translated, ct.translateDevFlag()...)
		case "global", "g":
			translated = append(translated, ct.translateGlobalFlag()...)
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

func (ct *CommandTranslator) translateDevFlag() []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--save-dev"}
	case PackageManagerYarn:
		return []string{"--dev"}
	case PackageManagerPnpm:
		return []string{"--save-dev"}
	case PackageManagerBun:
		return []string{"--dev"}
	}
	return []string{}
}

func (ct *CommandTranslator) translatePeerFlag() []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--save-peer"}
	case PackageManagerYarn:
		return []string{"--peer"}
	case PackageManagerPnpm:
		return []string{"--save-peer"}
	case PackageManagerBun:
		return []string{"--peer"}
	}
	return []string{}
}

func (ct *CommandTranslator) translateOptionalFlag() []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--save-optional"}
	case PackageManagerYarn:
		return []string{"--optional"}
	case PackageManagerPnpm:
		return []string{"--save-optional"}
	case PackageManagerBun:
		return []string{"--optional"}
	}
	return []string{}
}

func (ct *CommandTranslator) translateExactFlag() []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--save-exact"}
	case PackageManagerYarn:
		return []string{"--exact"}
	case PackageManagerPnpm:
		return []string{"--save-exact"}
	case PackageManagerBun:
		return []string{"--exact"}
	}
	return []string{}
}

func (ct *CommandTranslator) translateGlobalFlag() []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--global"}
	case PackageManagerYarn:
		return []string{"--global"}
	case PackageManagerPnpm:
		return []string{"--global"}
	case PackageManagerBun:
		return []string{"--global"}
	}
	return []string{}
}

func (ct *CommandTranslator) translateProductionFlag() []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--production"}
	case PackageManagerYarn:
		return []string{"--production"}
	case PackageManagerPnpm:
		return []string{"--prod"}
	case PackageManagerBun:
		return []string{"--production"}
	}
	return []string{}
}

func (ct *CommandTranslator) translateOmitFlag(value string) []string {
	switch ct.packageManager {
	case PackageManagerNpm:
		return []string{"--omit", value}
	case PackageManagerYarn:
		if value == "dev" {
			return []string{"--production"}
		}
		return []string{}
	case PackageManagerPnpm:
		if value == "dev" {
			return []string{"--prod"}
		}
		return []string{"--omit", value}
	case PackageManagerBun:
		if value == "dev" {
			return []string{"--production"}
		}
		return []string{}
	}
	return []string{}
}
