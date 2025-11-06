package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"pm/internal/detector"
	"pm/internal/executor"
	"pm/internal/translator"
	"pm/internal/ui"
	"pm/internal/version"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:                "pm",
		Short:              "A universal package manager wrapper",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 && (args[0] == "-v" || args[0] == "--version") {
				fmt.Println(version.GetVersion())
				return
			}

			pm, err := detector.Detect()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			if len(args) == 0 {
				script, err := ui.ShowScriptPrompt()
				if err != nil {
					if err.Error() != "cancelled" {
						fmt.Fprintln(os.Stderr, err)
					}
					return
				}
				executor.Run(pm, translator.Commands.Run, script.Name)
				return
			}

			tr := translator.New(pm)
			translated := tr.Translate(pm, args)
			if err := executor.Execute(pm, translated); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		},
		Args: cobra.ArbitraryArgs,
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
