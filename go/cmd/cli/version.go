package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev" // This will be set during build

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Personal Agent",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Personal Agent %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
