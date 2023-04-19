package main

type CommandAlias struct {
	NPM  []string
	Yarn []string
	Pnpm []string
}

type commands struct {
	Add     CommandAlias
	CI      CommandAlias
	Install CommandAlias
	Run     CommandAlias
}

var Commands = commands{
	Add: CommandAlias{
		NPM:  []string{"install"},
		Yarn: []string{"add"},
		Pnpm: []string{"add"},
	},
	CI: CommandAlias{
		NPM:  []string{"ci"},
		Yarn: []string{"install", "--frozen-lockfile"},
		Pnpm: []string{"install", "--frozen-lockfile"},
	},
	Install: CommandAlias{
		NPM:  []string{"install"},
		Yarn: []string{"install"},
		Pnpm: []string{"install"},
	},
	Run: CommandAlias{
		NPM:  []string{"run"},
		Yarn: []string{"run"},
		Pnpm: []string{"run"},
	},
}
