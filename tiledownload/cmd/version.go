package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "check version",
	Long:  `check version .`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("1.0.8")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

}
