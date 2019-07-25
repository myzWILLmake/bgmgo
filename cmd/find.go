package cmd

import (
	"fmt"

	"github.com/myzWILLmake/bgmgo/parser"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find torrent descriptions from a torrent site by a provided pattern.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		find(args)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}

func find(args []string) {
	dataSource := viper.GetString("data-source")
	if parser.ParserCtor[dataSource] == nil {
		fmt.Println("Error: Not a valid data source", dataSource)
		return
	}

	web := parser.ParserCtor[dataSource]()
	err := web.Request(args)
	if err != nil {
		fmt.Println("Error:", err)
	}

	filterMap := map[string]int{
		"no":           0,
		"type":         1,
		"organization": 2,
		"size":         3,
		"title":        4,
	}

	infos := web.ShowFindResult(filterMap, 5)

	n := len(infos)
	if n == 0 {
		fmt.Println("Not found any result, please check your pattern")
		return
	}
	fmt.Printf("Found %d record(s):\n", n)
	for _, info := range infos {
		for _, s := range info {
			fmt.Print(s, "  ")
		}
		fmt.Println()
	}
}
