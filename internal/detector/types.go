package detector

type PackageManager string

const (
	NPM  PackageManager = "npm"
	Yarn PackageManager = "yarn"
	Pnpm PackageManager = "pnpm"
	Bun  PackageManager = "bun"
)
