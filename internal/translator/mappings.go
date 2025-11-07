package translator

import "pm/internal/detector"

// CommandAlias maps package managers to their command equivalents
type CommandAlias map[detector.PackageManager][]string

type commands struct {
	Add       CommandAlias
	CI        CommandAlias
	Install   CommandAlias
	Uninstall CommandAlias
	Run       CommandAlias
}

// Commands contains pre-defined command mappings between package managers
var Commands = commands{
	Add: CommandAlias{
		detector.NPM:       []string{"install"},
		detector.Yarn:      []string{"add"},
		detector.YarnBerry: []string{"add"},
		detector.Pnpm:      []string{"add"},
		detector.Bun:       []string{"add"},
	},
	CI: CommandAlias{
		detector.NPM:       []string{"ci"},
		detector.Yarn:      []string{"install", "--frozen-lockfile"},
		detector.YarnBerry: []string{"install", "--immutable"},
		detector.Pnpm:      []string{"install", "--frozen-lockfile"},
		detector.Bun:       []string{"install", "--frozen-lockfile"},
	},
	Install: CommandAlias{
		detector.NPM:       []string{"install"},
		detector.Yarn:      []string{"install"},
		detector.YarnBerry: []string{"install"},
		detector.Pnpm:      []string{"install"},
		detector.Bun:       []string{"install"},
	},
	Uninstall: CommandAlias{
		detector.NPM:       []string{"uninstall"},
		detector.Yarn:      []string{"remove"},
		detector.YarnBerry: []string{"remove"},
		detector.Pnpm:      []string{"remove"},
		detector.Bun:       []string{"remove"},
	},
	Run: CommandAlias{
		detector.NPM:       []string{"run"},
		detector.Yarn:      []string{"run"},
		detector.YarnBerry: []string{"run"},
		detector.Pnpm:      []string{"run"},
		detector.Bun:       []string{"run"},
	},
}
