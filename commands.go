package main

// key is package manager
type CommandAlias map[PackageManager][]string

type commands struct {
	Add       CommandAlias
	CI        CommandAlias
	Install   CommandAlias
	Uninstall CommandAlias
	Run       CommandAlias
}

var Commands = commands{
	Add: CommandAlias{
		PackageManagerNpm:  []string{"install"},
		PackageManagerYarn: []string{"add"},
		PackageManagerPnpm: []string{"add"},
		PackageManagerBun:  []string{"add"},
	},
	CI: CommandAlias{
		PackageManagerNpm:  []string{"ci"},
		PackageManagerYarn: []string{"install", "--frozen-lockfile"},
		PackageManagerPnpm: []string{"install", "--frozen-lockfile"},
		PackageManagerBun:  []string{"install", "--frozen-lockfile"},
	},
	Install: CommandAlias{
		PackageManagerNpm:  []string{"install"},
		PackageManagerYarn: []string{"install"},
		PackageManagerPnpm: []string{"install"},
		PackageManagerBun:  []string{"install"},
	},
	Uninstall: CommandAlias{
		PackageManagerNpm:  []string{"uninstall"},
		PackageManagerYarn: []string{"remove"},
		PackageManagerPnpm: []string{"remove"},
		PackageManagerBun:  []string{"remove"},
	},
	Run: CommandAlias{
		PackageManagerNpm:  []string{"run"},
		PackageManagerYarn: []string{"run"},
		PackageManagerPnpm: []string{"run"},
		PackageManagerBun:  []string{"run"},
	},
}
