package main

import (
	"strings"

	"github.com/spf13/pflag"
)

type flagAlias map[PackageManager][]string

type flags struct {
	Dev      flagAlias
	Peer     flagAlias
	Optional flagAlias
	Global   flagAlias
	Exact    flagAlias
}

var Flags = flags{
	Dev: flagAlias{
		PackageManagerNpm:  []string{"--save-dev"},
		PackageManagerYarn: []string{"--dev"},
		PackageManagerPnpm: []string{"--save-dev"},
		PackageManagerBun:  []string{"--dev"},
	},
	Peer: flagAlias{
		PackageManagerNpm:  []string{"--save-peer"},
		PackageManagerYarn: []string{"--peer"},
		PackageManagerPnpm: []string{"--save-peer"},
		PackageManagerBun:  []string{"--peer"},
	},
	Optional: flagAlias{
		PackageManagerNpm:  []string{"--save-optional"},
		PackageManagerYarn: []string{"--optional"},
		PackageManagerPnpm: []string{"--save-optional"},
		PackageManagerBun:  []string{"--optional"},
	},
	Global: flagAlias{
		PackageManagerNpm:  []string{"--global"},
		PackageManagerYarn: []string{"--global"},
		PackageManagerPnpm: []string{"--global"},
		PackageManagerBun:  []string{"--global"},
	},
	Exact: flagAlias{
		PackageManagerNpm:  []string{"--save-exact"},
		PackageManagerYarn: []string{"--exact"},
		PackageManagerPnpm: []string{"--save-exact"},
		PackageManagerBun:  []string{"--exact"},
	},
}

func isFlag(arg string) (bool, bool) {
	if strings.HasPrefix(arg, "--") {
		return true, false
	}
	if strings.HasPrefix(arg, "-") {
		return true, true
	}
	return false, false
}

func filterFlags(args []string, flags *pflag.FlagSet) []string {
	filteredArgs := []string{}

	for i := 0; i < len(args); {
		arg := args[i]

		if ok, isShorthand := isFlag(arg); ok {
			var flag *pflag.Flag
			if isShorthand {
				flag = flags.ShorthandLookup(strings.TrimLeft(arg, "-"))
			} else {
				flag = flags.Lookup(strings.TrimLeft(arg, "-"))
			}
			if flag != nil {
				// If the flag has a value (e.g. --flag=value or -f value), skip it
				if flag.Value.Type() != "bool" && i+1 < len(args) {
					i++
				}
			} else {
				filteredArgs = append(filteredArgs, arg)
			}
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
		i++
	}

	return filteredArgs
}
