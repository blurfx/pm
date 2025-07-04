package main

import (
	"fmt"
	"log"
	"os"

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
			if len(args) == 0 {
				script, err := showScriptPrompt()
				if err != nil {
					if err.Error() != "cancelled" {
						fmt.Fprintln(os.Stderr, err)
					}
					return
				}
				Exec(Commands.Run, script.Name)
				return
			}

			translator := NewCommandTranslator(packageManager)
			translated := translator.Translate(args)

			if err := ExecuteTranslatedCommand(packageManager, translated); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
		Args: cobra.ArbitraryArgs,
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
