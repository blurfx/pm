package main

import "github.com/spf13/cobra"

func GetBoolFlag(cmd *cobra.Command, names ...string) bool {
	for _, name := range names {
		val, err := cmd.Flags().GetBool(name)
		if err != nil || !val {
			continue
		}
		return val
	}
	return false
}
