package main

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
