package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var unsubCmd = &cobra.Command{
	Use:   "unsub",
	Short: "Remove subscriptions by sub-number",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		unsub(args)
	},
}

func init() {
	rootCmd.AddCommand(unsubCmd)
}

func unsub(args []string) {
	if len(args) < 1 {
		fmt.Println("Please enter the sub-number(s) you want to ubsubscribe.")
		fmt.Println("	Usage: bgmgo unsub [Sub-number ...]")
		return
	}

	for _, argv := range args {
		no, err := strconv.ParseInt(args[0], 10, 32)
		if err != nil {
			fmt.Println("Unsubcription failed:", argv, "is not a vaild sub-number.")
			return
		}

		if _, ok := globalData.Sublist[int(no)]; !ok {
			fmt.Println("Unsubcription failed:", argv, "is not existed.")
			return
		}

		delete(globalData.Sublist, int(no))
	}

	err := writeData()
	if err != nil {
		fmt.Println("Unsubscribe failed:", err)
		return
	}

	fmt.Println("Unsubscription succeed! Sub-number(s):", args)
}
