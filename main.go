package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "add",
				Aliases:     []string{"i", "install"},
				Description: "Add dependency if package name is given.",
				Usage:       "pm add <package name> ",
				ArgsUsage:   "<package name>",
				Action: func(ctx *cli.Context) error {
					if ctx.Bool("frozen-lockfile") {
						return Exec(Commands.CI)
					}
					if !ctx.Args().Present() {
						log.Fatalf("pack")
						return nil
					}
					if ctx.NumFlags() != 0 {
						log.Fatalf("fuck")
					}

					return Exec(Commands.Add, ctx.Args().Slice()...)
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "dev",
						Aliases: []string{"save-dev", "d", "D"},
						Usage:   "install package as dev dependency",
					},
					&cli.BoolFlag{
						Name:  "frozen-lockfile",
						Usage: "don't generate a lockfile and fail if an update is needed",
					},
				},
			},
			{
				Name:        "install",
				Aliases:     []string{"i"},
				Description: "Add dependency if package name is given. Otherwise, install the package",
				Action: func(ctx *cli.Context) error {
					if ctx.Bool("frozen-lockfile") {
						return Exec(Commands.CI)
					}
					if ctx.Args().Present() {
						return Exec(Commands.Add, ctx.Args().Slice()...)
					}
					if ctx.NumFlags() != 0 {
						log.Fatalf("")
					}

					return Exec(Commands.Install, ctx.Args().Slice()...)
				},
			},
			{
				Name: "ci",
				Action: func(ctx *cli.Context) error {
					return Exec(Commands.CI)
				},
			},
			{
				Name: "run",
				Action: func(ctx *cli.Context) error {
					return Exec(Commands.Run, ctx.Args().Slice()...)
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			return PassThrough(ctx.Args().Slice()...)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
