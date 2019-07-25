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
	Short: "Print the version number of Bgmgo",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bgmgo Anime Subscription Tool v0.1")
	},
}
