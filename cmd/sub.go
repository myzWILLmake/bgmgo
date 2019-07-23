package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Add a new subscription.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		sub(args)
	},
}

var name string
var progress float64

func init() {
	rootCmd.AddCommand(subCmd)

	subCmd.Flags().StringVarP(&name, "name", "n", "", "name for this subscription")
	subCmd.Flags().Float64VarP(&progress, "progress", "p", 0, "progress for this subscription")
}

func sub(args []string) {
	fmt.Println(args)
	if len(args) != 1 {
		if len(args) == 0 {
			fmt.Println("Please enter a pattern for the subscription.")
		} else if len(args) > 1 {
			fmt.Println("Only pattern can be entered. Please use flags to define name or progress.")
		}
		return
	}

	pattern := args[0]

	globalData.SubMaxNo++
	no := globalData.SubMaxNo

	subItem := SubItem{no, name, 0, pattern, progress}
	globalData.Sublist[no] = subItem

	err := writeData()
	if err != nil {
		return
	}
	fmt.Println("Subscribe succeed!")
}
