package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:                "pm",
		Short:              "A universal package manager wrapper",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			packageManager, err := DetectPackageManager()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			if IsBuiltInCommand(packageManager, args...) {
				if err := PassThrough(args...); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
				return
			}
			allArgs := append([]string{"run"}, args...)
			if err := PassThrough(allArgs...); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
		Args: cobra.ArbitraryArgs,
	}

	addCmd := &cobra.Command{
		Use:     "add <package name>",
		Aliases: []string{"i", "install"},
		Short:   "Add dependency if package name is given.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalf("No package name given")
			}

			cmd.Flags().Parse(args)

			dev := GetBoolFlag(cmd, "save-dev", "dev")
			peer := GetBoolFlag(cmd, "save-peer", "peer")
			optional := GetBoolFlag(cmd, "save-optional", "optional")
			global := GetBoolFlag(cmd, "global")
			exact := GetBoolFlag(cmd, "exact")
			frozenLockfile := GetBoolFlag(cmd, "frozen-lockfile")

			if frozenLockfile {
				Exec(Commands.CI)
				return
			}

			var flags []flagAlias
			if dev {
				flags = append(flags, Flags.Dev)
			}
			if peer {
				flags = append(flags, Flags.Peer)
			}
			if optional {
				flags = append(flags, Flags.Optional)
			}
			if global {
				flags = append(flags, Flags.Global)
			}
			if exact {
				flags = append(flags, Flags.Exact)
			}
			args = filterFlags(args, cmd.Flags())

			ExecWithFlag(Commands.Add, flags, args...)

			nonTypedPackages := []string{}
			for _, arg := range args {
				if !strings.HasPrefix(arg, "@types/") {
					ok, _ := CheckPackageExists("@types/" + arg)
					if !ok {
						continue
					}
					nonTypedPackages = append(nonTypedPackages, "@types/"+arg)
				}
			}
			if !dev && !peer && !optional && !global {
				for _, packageName := range nonTypedPackages {
					ExecWithFlag(Commands.Add, []flagAlias{Flags.Dev}, packageName)
				}
			}
		},
		DisableFlagParsing: true,
	}

	addCmd.Flags().BoolP("save-dev", "D", false, "install package as dev dependency")
	addCmd.Flags().Bool("dev", false, "install package as dev dependency")
	addCmd.Flags().MarkHidden("dev")
	addCmd.Flags().BoolP("save-peer", "P", false, "install package as peer dependency")
	addCmd.Flags().Bool("peer", false, "install package as peer dependency")
	addCmd.Flags().MarkHidden("peer")
	addCmd.Flags().BoolP("save-optional", "O", false, "install package as optional dependency")
	addCmd.Flags().Bool("optional", false, "install package as optional dependency")
	addCmd.Flags().MarkHidden("optional")
	addCmd.Flags().BoolP("global", "g", false, "install package globally")
	addCmd.Flags().BoolP("exact", "E", false, "install exact version")
	addCmd.Flags().Bool("frozen-lockfile", false, "don't generate a lockfile and fail if an update is needed")
	addCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}

	installCmd := &cobra.Command{
		Use:     "install [package name]",
		Aliases: []string{"i"},
		Short:   "Add dependency if package name is given. Otherwise, install the package",
		Run: func(cmd *cobra.Command, args []string) {
			frozenLockfile, _ := cmd.Flags().GetBool("frozen-lockfile")

			if frozenLockfile {
				Exec(Commands.CI)
				return
			}

			if len(args) > 0 {
				addCmd.Run(cmd, args)
			} else {
				Exec(Commands.Install)
			}
		},
		DisableFlagParsing: true,
	}

	installCmd.Flags().Bool("frozen-lockfile", false, "don't generate a lockfile and fail if an update is needed")
	installCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}

	uninstallCmd := &cobra.Command{
		Use:     "uninstall <package name>",
		Aliases: []string{"rm", "remove", "un"},
		Short:   "Uninstall a package",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalf("No package name given")
			}

			Exec(Commands.Uninstall, args...)
		},
		DisableFlagParsing: true,
	}
	uninstallCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}

	ciCmd := &cobra.Command{
		Use:   "ci",
		Short: "CI command",
		Run: func(cmd *cobra.Command, args []string) {
			Exec(Commands.CI)
		},
	}
	ciCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run command",
		Run: func(cmd *cobra.Command, args []string) {
			Exec(Commands.Run, args...)
		},
		DisableFlagParsing: true,
	}
	runCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}

	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(ciCmd)
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
