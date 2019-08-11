package cmd

import (
	"fmt"

	"github.com/spf13/viper"

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

var subProgress float64
var subSource string

func init() {
	rootCmd.AddCommand(subCmd)

	subCmd.Flags().Float64VarP(&subProgress, "progress", "p", 0, "set progress for this subscription")
	subCmd.Flags().StringVarP(&subSource, "source", "s", "",
		"set data source where to search magnet links,\n"+
			"available options:\n"+
			"	dmhy\n"+
			"	bangumi_moe\n"+
			"	nyaa")
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

	if subSource == "" {
		subSource = viper.GetString("data-source")
	}

	subItem := SubItem{no, name, 0, pattern, subProgress, subSource}
	globalData.Sublist[no] = &subItem

	err := writeData()
	if err != nil {
		fmt.Println("Subscription failed:", err)
		return
	}
	fmt.Println("Subscription succeed! Sub-number is", no)
}
