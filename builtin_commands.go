package main

// Internal CLI commands

// v11
var NPM_COMMANDS = [][]string{
	{"access"},
	{"adduser"},
	{"audit"},
	{"bugs"},
	{"cache"},
	{"ci"},
	{"completion"},
	{"config"},
	{"dedupe"},
	{"deprecate"},
	{"diff"},
	{"dist-tag"},
	{"docs"},
	{"doctor"},
	{"edit"},
	{"exec"},
	{"explain"},
	{"explore"},
	{"find-dupes"},
	{"fund"},
	{"help"},
	{"help-search"},
	{"init"},
	{"install"},
	{"install-ci-test"},
	{"install-test"},
	{"link"},
	{"login"},
	{"logout"},
	{"ls"},
	{"org"},
	{"outdated"},
	{"owner"},
	{"pack"},
	{"ping"},
	{"pkg"},
	{"prefix"},
	{"profile"},
	{"prune"},
	{"publish"},
	{"query"},
	{"rebuild"},
	{"repo"},
	{"restart"},
	{"root"},
	{"run"},
	{"sbom"},
	{"search"},
	{"shrinkwrap"},
	{"star"},
	{"stars"},
	{"start"},
	{"stop"},
	{"team"},
	{"test"},
	{"token"},
	{"undeprecate"},
	{"uninstall"},
	{"unpublish"},
	{"unstar"},
	{"update"},
	{"version"},
	{"view"},
	{"whoami"},
}

var YARN_CLASSIC_COMMANDS = [][]string{
	{"add"},
	{"audit"},
	{"autoclean"},
	{"bin"},
	{"cache"},
	{"check"},
	{"config"},
	{"create"},
	{"dedupe"},
	{"generate-lock-entry"},
	{"global"},
	{"help"},
	{"import"},
	{"info"},
	{"init"},
	{"install"},
	{"licenses"},
	{"link"},
	{"list"},
	{"lockfile"},
	{"login"},
	{"logout"},
	{"outdated"},
	{"owner"},
	{"pack"},
	{"policies"},
	{"prune"},
	{"publish"},
	{"remove"},
	{"run"},
	{"self-update"},
	{"tag"},
	{"team"},
	{"test"},
	{"unlink"},
	{"upgrade"},
	{"upgrade-interactive"},
	{"version"},
	{"versions"},
	{"why"},
	{"workspace"},
	{"workspaces"},
}

var YARN2_COMMANDS = [][]string{
	{"add"},
	{"bin"},
	{"cache", "clean"},
	{"config", "get"},
	{"config", "set"},
	{"config", "unset"},
	{"dedupe"},
	{"dlx"},
	{"exec"},
	{"explain"},
	{"explain", "peer-requirements"},
	{"info"},
	{"init"},
	{"init"},
	{"install"},
	{"link"},
	{"node"},
	{"npm", "audit"},
	{"pack"},
	{"patch"},
	{"patch-commit"},
	{"rebuild"},
	{"remove"},
	{"run"},
	{"set", "resolution"},
	{"set", "version"},
	{"set", "version", "from", "sources"},
	{"stage"},
	{"unlink"},
	{"unplug"},
	{"up"},
	{"why"},
	{"constraints"},
	{"constraints", "query"},
	{"constraints", "source"},
	{"npm", "info"},
	{"npm", "login"},
	{"npm", "logout"},
	{"npm", "publish"},
	{"npm", "tag", "add"},
	{"npm", "tag", "list"},
	{"npm", "tag", "remove"},
	{"npm", "whoami"},
	{"plugin", "check"},
	{"plugin", "import", "from", "sources"},
	{"plugin", "list"},
	{"plugin", "remove"},
	{"plugin", "runtime"},
	{"search"},
	{"upgrade-interactive"},
	{"version", "apply"},
	{"version", "check"},
	{"workspace"},
	{"workspaces", "focus"},
	{"workspaces", "foreach"},
	{"workspaces", "list"},
}

var PNPM_COMMANDS = [][]string{
	// Manage dependencies
	{"add"},
	{"install"},
	{"update"},
	{"remove"},
	{"link"},
	{"unlink"},
	{"import"},
	{"rebuild"},
	{"prune"},
	{"fetch"},
	{"install-test"},
	{"dedupe"},

	// Patch dependencies
	{"patch"},
	{"patch-commit"},
	{"patch-remove"},

	// Review dependencies
	{"audit"},
	{"list"},
	{"outdated"},
	{"why"},
	{"licenses"},

	// Run scripts
	{"run"},
	{"test"},
	{"exec"},
	{"dlx"},
	{"create"},
	{"start"},
	{"approve-builds"},
	{"ignored-builds"},

	// Manage environments
	{"env"},

	// Inspect the store
	{"cat-file"},
	{"cat-index"},
	{"find-hash"},

	// Manage cache
	{"cache", "list"},
	{"cache", "list-registries"},
	{"cache", "view"},
	{"cache", "delete"},

	// Miscellaneous
	{"self-update"},
	{"publish"},
	{"pack"},
	{"-r"}, // pnpm -r, --recursive
	{"--recursive"},
	{"recursive"},
	{"server"},
	{"store"},
	{"root"},
	{"bin"},
	{"setup"},
	{"init"},
	{"deploy"},
	{"doctor"},
	{"config"},
}

var BUN_COMMANDS = [][]string{
	{"run"},
	{"test"},
	{"x"},
	{"repl"},
	{"exec"},
	{"install"},
	{"i"},
	{"add"},
	{"a"},
	{"remove"},
	{"rm"},
	{"audit"},
	{"outdated"},
	{"link"},
	{"unlink"},
	{"publish"},
	{"patch"},
	{"pm"},
	{"info"},
	{"build"},
	{"init"},
	{"create"},
	{"c"},
	{"upgrade"},
}

var IsBuiltInCommand = func(packageManager PackageManager, arg ...string) bool {
	var commands [][]string
	switch packageManager {
	case PackageManagerNpm:
		commands = NPM_COMMANDS
	case PackageManagerYarn:
		commands = YARN_CLASSIC_COMMANDS
	case PackageManagerPnpm:
		commands = PNPM_COMMANDS
	case PackageManagerBun:
		commands = BUN_COMMANDS
	default:
		return false
	}
	for _, cmd := range commands {
		if len(arg) < len(cmd) {
			continue
		}
		match := true
		for i := range cmd {
			if arg[i] != cmd[i] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
