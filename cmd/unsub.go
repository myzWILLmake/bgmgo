package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var unsubCmd = &cobra.Command{
	Use:   "unsub",
	Short: "",
	Long:  `Delete subscriptions by sub-no`,
	Run: func(cmd *cobra.Command, args []string) {
		unsub(args)
	},
}

func init() {
	rootCmd.AddCommand(unsubCmd)
}

func unsub(args []string) {
	if len(args) < 1 {
		fmt.Println("Please enter the sub-number you want to ubsubscribe.")
		return
	}

	no, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		fmt.Println("Please enter a sub-number")
		return
	}

	if _, ok := globalData.Sublist[int(no)]; !ok {
		fmt.Println("Please enter a sub-number")
		return
	}

	delete(globalData.Sublist, int(no))
	fmt.Println("Unsubscribe succeed!")
}
