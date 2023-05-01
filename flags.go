package main

import (
	"github.com/spf13/pflag"
	"strings"
)

type flagAlias struct {
	NPM  []string
	Yarn []string
	Pnpm []string
}

type flags struct {
	Dev      flagAlias
	Peer     flagAlias
	Optional flagAlias
	Global   flagAlias
	Exact    flagAlias
}

var Flags = flags{
	Dev: flagAlias{
		NPM:  []string{"--save-dev"},
		Yarn: []string{"--dev"},
		Pnpm: []string{"--save-dev"},
	},
	Peer: flagAlias{
		NPM:  []string{"--save-peer"},
		Yarn: []string{"--peer"},
		Pnpm: []string{"--save-peer"},
	},
	Optional: flagAlias{
		NPM:  []string{"--save-optional"},
		Yarn: []string{"--optional"},
		Pnpm: []string{"--save-optional"},
	},
	Global: flagAlias{
		NPM:  []string{"--Global"},
		Yarn: []string{"--Global"},
		Pnpm: []string{"--Global"},
	},
	Exact: flagAlias{
		NPM:  []string{"--save-exact"},
		Yarn: []string{"--exact"},
		Pnpm: []string{"--save-exact"},
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
