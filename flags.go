package main

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
