package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Add a new subscription with a pattern and a name.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sub(args)
	},
}

var progress float64

func init() {
	rootCmd.AddCommand(subCmd)

	subCmd.Flags().Float64VarP(&progress, "progress", "p", 0, "set progress for this subscription")
}

func sub(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: bgmgo sub [Pattern] [Name]")
		return
	}

	pattern := args[0]
	name := args[1]

	globalData.SubMaxNo++
	no := globalData.SubMaxNo

	subItem := SubItem{no, name, 0, pattern, progress}
	globalData.Sublist[no] = &subItem

	err := writeData()
	if err != nil {
		fmt.Println("Subscription failed:", err)
		return
	}
	fmt.Println("Subscription succeed! Sub-number is", no)
}
