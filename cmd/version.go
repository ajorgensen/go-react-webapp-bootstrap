package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Tarkov Maps v0.1")
	},
}
